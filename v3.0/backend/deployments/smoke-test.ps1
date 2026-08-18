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
  Invoke-RestMethod -Method Post -Uri http://localhost:6081/api/v1/tenants -Headers @{ Authorization = "Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)" } -ContentType application/json -Body '{"tenant_id":"alpha_logistics","display_name":"Alpha Logistics"}' | Out-Null
} catch { if ($_.Exception.Response.StatusCode.value__ -ne 409) { throw } }
$projectResponse = Invoke-RestMethod -Method Post -Uri http://localhost:6081/api/v1/projects -Headers $operatorHeaders -ContentType application/json -Body (@{name="Chennai fleet demo $smokeID";description="Phase 2 reproducible identity proof"}|ConvertTo-Json)
$deviceBody = @{ device_id=$smokeID; project_id=$projectResponse.data.project_id; device_type_id="delivery_drone"; display_name="Authenticated smoke device" } | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri http://localhost:6081/api/v1/devices -Headers $operatorHeaders -ContentType application/json -Body $deviceBody | Out-Null
Invoke-RestMethod -Method Put -Uri "http://localhost:6081/api/v1/devices/$smokeID/capabilities/navigate" -Headers $operatorHeaders -ContentType application/json -Body '{"configuration":{}}' | Out-Null
Invoke-RestMethod -Method Put -Uri "http://localhost:6081/api/v1/devices/$smokeID/capabilities/receive_relocation_command" -Headers $operatorHeaders -ContentType application/json -Body '{"configuration":{}}' | Out-Null
Invoke-RestMethod -Method Post -Uri "http://localhost:6081/api/v1/devices/$smokeID/activate" -Headers $operatorHeaders | Out-Null
$credentialResponse = Invoke-RestMethod -Method Post -Uri "http://localhost:6081/api/v1/devices/$smokeID/credentials" -Headers $operatorHeaders -ContentType application/json -Body '{}'
$env:DEVICE_TOKEN = $credentialResponse.data.secret
$env:DEVICE_CREDENTIAL_ID = $credentialResponse.data.credential.credential_id
$env:OPERATOR_TOKEN = $env:DEV_PLATFORM_ADMIN_TOKEN
$env:SMOKE_DEVICE_ID = $smokeID
Push-Location $backendDir
try { $smokeResultJson = go run ./cmd/smoke } finally { Pop-Location }
if ($LASTEXITCODE -ne 0) { throw "Smoke client failed" }
if (-not $smokeResultJson) { throw "Smoke client returned no result" }
$smokeResult = $smokeResultJson | ConvertFrom-Json
$smokeID = $smokeResult.id

$count = 0
$archiveDeadline = (Get-Date).AddSeconds(15)
while ((Get-Date) -lt $archiveDeadline -and [int]$count -lt 1) {
  $count = docker compose -f (Join-Path $deploymentDir "docker-compose.yml") exec -T postgres `
    psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM telemetry_history WHERE device_id='$smokeID'"
  if ([int]$count -lt 1) { Start-Sleep -Milliseconds 250 }
}
if ([int]$count -lt 1) { throw "Telemetry was not archived in PostgreSQL" }

Write-Host "PASS: Simulator -> Gateway -> Kafka -> Engine -> Redis -> PostgreSQL -> Dashboard ($smokeID, $($smokeResult.end_to_end_latency_ms) ms)"
