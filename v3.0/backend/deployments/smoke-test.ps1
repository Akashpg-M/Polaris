$ErrorActionPreference = "Stop"
$deploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir = Resolve-Path (Join-Path $deploymentDir "..")

docker compose -f (Join-Path $deploymentDir "docker-compose.yml") up -d --build --wait
if ($LASTEXITCODE -ne 0) { throw "Docker Compose stack failed to start" }
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
