$ErrorActionPreference = "Stop"
$deploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir = Resolve-Path (Join-Path $deploymentDir "..")

function New-RandomToken([string]$kind) {
  $rng=[Security.Cryptography.RandomNumberGenerator]::Create()
  try { $prefixBytes = New-Object byte[] 8; $rng.GetBytes($prefixBytes); $secretBytes = New-Object byte[] 32; $rng.GetBytes($secretBytes) } finally { $rng.Dispose() }
  $prefix = ([BitConverter]::ToString($prefixBytes) -replace '-','').ToLowerInvariant()
  $secret = ([BitConverter]::ToString($secretBytes) -replace '-','').ToLowerInvariant()
  return "pol_${kind}_${prefix}.${secret}"
}

$env:DEV_PLATFORM_ADMIN_TOKEN = New-RandomToken "op"
$smokeID = "SMOKE-$([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds())"

docker compose -f (Join-Path $deploymentDir "docker-compose.yml") up -d --build --wait
if ($LASTEXITCODE -ne 0) { throw "Docker Compose stack failed to start" }
$operatorHeaders = @{ Authorization = "Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)"; "X-Tenant-ID" = "alpha_logistics" }
try {
  Invoke-RestMethod -Method Post -Uri http://127.0.0.1:6081/api/v1/tenants -Headers @{ Authorization = "Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)" } -ContentType application/json -Body '{"tenant_id":"alpha_logistics","display_name":"Alpha Logistics"}' | Out-Null
} catch { if ($_.Exception.Response.StatusCode.value__ -ne 409) { throw } }
$projectResponse = Invoke-RestMethod -Method Post -Uri http://127.0.0.1:6081/api/v1/projects -Headers $operatorHeaders -ContentType application/json -Body (@{name="Chennai fleet demo $smokeID";description="Phase 2 reproducible identity proof"}|ConvertTo-Json)
$deviceBody = @{ device_id=$smokeID; project_id=$projectResponse.data.project_id; device_type_id="delivery_drone"; display_name="Authenticated smoke device" } | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri http://127.0.0.1:6081/api/v1/devices -Headers $operatorHeaders -ContentType application/json -Body $deviceBody | Out-Null
Invoke-RestMethod -Method Put -Uri "http://127.0.0.1:6081/api/v1/devices/$smokeID/capabilities/navigate" -Headers $operatorHeaders -ContentType application/json -Body '{"configuration":{}}' | Out-Null
Invoke-RestMethod -Method Put -Uri "http://127.0.0.1:6081/api/v1/devices/$smokeID/capabilities/receive_relocation_command" -Headers $operatorHeaders -ContentType application/json -Body '{"configuration":{}}' | Out-Null
Invoke-RestMethod -Method Post -Uri "http://127.0.0.1:6081/api/v1/devices/$smokeID/activate" -Headers $operatorHeaders | Out-Null
$credentialResponse = Invoke-RestMethod -Method Post -Uri "http://127.0.0.1:6081/api/v1/devices/$smokeID/credentials" -Headers $operatorHeaders -ContentType application/json -Body '{}'
$env:DEVICE_TOKEN = $credentialResponse.data.secret
$env:DEVICE_CREDENTIAL_ID = $credentialResponse.data.credential.credential_id
$env:OPERATOR_TOKEN = $env:DEV_PLATFORM_ADMIN_TOKEN
$env:SMOKE_DEVICE_ID = $smokeID
$env:SMOKE_TENANT_ID = 'alpha_logistics'
$env:SMOKE_NODE_TYPE = '5'
$env:SMOKE_LAT = '13.0067'
$env:SMOKE_LON = '80.2206'
$env:SMOKE_SEQUENCE = '1'
$env:SMOKE_WAIT_FOR_PROJECTION = 'true'
$env:SMOKE_WAIT_FOR_MATCH = 'true'
foreach ($name in @('SMOKE_BOOT_ID','SMOKE_BOOT_STARTED_AT')) {
  [Environment]::SetEnvironmentVariable($name, $null, 'Process')
}
Push-Location $backendDir
try { $smokeResultJson = go run ./cmd/smoke } finally { Pop-Location }
if ($LASTEXITCODE -ne 0) { throw "Smoke client failed" }
if (-not $smokeResultJson) { throw "Smoke client returned no result" }
$smokeResult = $smokeResultJson | ConvertFrom-Json
$smokeID = $smokeResult.id

$count = 0
$archiveDeadline = (Get-Date).AddSeconds(60)
while ((Get-Date) -lt $archiveDeadline -and [int]$count -lt 1) {
  $count = docker compose -f (Join-Path $deploymentDir "docker-compose.yml") exec -T postgres `
    psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM telemetry_history WHERE device_id='$smokeID'"
  if ([int]$count -lt 1) { Start-Sleep -Milliseconds 250 }
}
if ([int]$count -lt 1) {
  $archiveGroup = docker compose -f (Join-Path $deploymentDir "docker-compose.yml") exec -T redpanda `
    rpk group describe polaris_archive_group -X brokers=localhost:29092
  throw "Telemetry was not archived in PostgreSQL within 60 seconds. Archive consumer state:`n$archiveGroup"
}

Write-Host "PASS: Simulator -> Gateway -> Kafka -> Engine -> Redis -> PostgreSQL -> Dashboard ($smokeID, $($smokeResult.end_to_end_latency_ms) ms)"
