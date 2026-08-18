$ErrorActionPreference = "Stop"
$deploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir = Resolve-Path (Join-Path $deploymentDir "..")
$rootDir = Resolve-Path (Join-Path $backendDir "..")

& (Join-Path $deploymentDir "smoke-test.ps1")

Push-Location $backendDir
try {
  go test ./...
  if ($LASTEXITCODE -ne 0) { throw "Unit reliability tests failed" }
  $env:REDIS_URL = "redis://localhost:6379/0"
  $env:POSTGRES_URL = "postgres://polaris_user:polaris_password@localhost:5432/polaris_core?sslmode=disable"
  go test -count=1 -tags=integration -v ./internal/application/stream
  if ($LASTEXITCODE -ne 0) { throw "Live dependency reliability tests failed" }
} finally { Pop-Location }

$gatewayHealth = Invoke-RestMethod http://localhost:6080/healthz
$gatewayReady = Invoke-RestMethod http://localhost:6080/readyz
$engineHealth = Invoke-RestMethod http://localhost:6081/healthz
$engineReady = Invoke-RestMethod http://localhost:6081/readyz
if ($gatewayHealth.status -ne "live" -or $gatewayReady.status -ne "ready" -or $engineHealth.status -ne "live" -or $engineReady.status -ne "ready") {
  throw "Health/readiness contract failed"
}

$partitions = docker compose -f (Join-Path $deploymentDir "docker-compose.yml") exec -T redpanda `
  rpk topic describe telemetry.ingress -p -X brokers=localhost:29092
if (($partitions | Select-String -Pattern '^\s*[0-2]\s+').Count -lt 3) { throw "telemetry.ingress does not have three partitions" }

$constraints = docker compose -f (Join-Path $deploymentDir "docker-compose.yml") exec -T postgres `
  psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM pg_indexes WHERE tablename='telemetry_history' AND indexname IN ('uq_telemetry_event_id','uq_telemetry_device_sequence')"
if ([int]$constraints -ne 2) { throw "PostgreSQL idempotency indexes are missing" }

Write-Host "PASS: Phase 1 at-least-once/idempotent reliability verification"
