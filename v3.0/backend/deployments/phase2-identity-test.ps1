$ErrorActionPreference = "Stop"
$deploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir = Resolve-Path (Join-Path $deploymentDir "..")
$composeFile = Join-Path $deploymentDir "docker-compose.yml"

function New-RandomToken([string]$kind) {
  $rng=[Security.Cryptography.RandomNumberGenerator]::Create()
  try { $prefixBytes = New-Object byte[] 8; $rng.GetBytes($prefixBytes); $secretBytes = New-Object byte[] 32; $rng.GetBytes($secretBytes) } finally { $rng.Dispose() }
  $prefix=([BitConverter]::ToString($prefixBytes)-replace '-','').ToLowerInvariant()
  $secret=([BitConverter]::ToString($secretBytes)-replace '-','').ToLowerInvariant()
  return "pol_${kind}_${prefix}.${secret}"
}
function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null) {
  $params=@{Method=$method;Uri=$uri}
  if($headers){$params.Headers=$headers}
  if($null-ne $body){$params.ContentType='application/json';$params.Body=($body|ConvertTo-Json -Depth 8)}
  return Invoke-RestMethod @params
}
function Run-IdentityCheck([string]$mode) {
  $env:IDENTITY_CHECK_MODE=$mode
  Push-Location $backendDir
  try { go run ./cmd/identitycheck } finally { Pop-Location }
  if($LASTEXITCODE-ne 0){throw "Identity check failed: $mode"}
}
function Wait-Connectivity([string]$expected,[int]$timeoutSeconds) {
  $deadline=(Get-Date).AddSeconds($timeoutSeconds)
  do { $current=Invoke-API GET "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/twin" $null $headers; if($current.data.connectivity.status-eq$expected){return $current}; Start-Sleep -Milliseconds 250 } while((Get-Date)-lt$deadline)
  throw "Expected $expected, got $($current.data.connectivity.status)"
}

$env:DEVICE_STALE_AFTER="5s"
$env:DEVICE_OFFLINE_AFTER="8s"
$env:OFFLINE_SCAN_INTERVAL="1s"
& (Join-Path $deploymentDir "smoke-test.ps1")
$headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'='alpha_logistics'}

Run-IdentityCheck "basic"
$twin=Invoke-API GET "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/twin" $null $headers
if($twin.data.tenant_id-ne'alpha_logistics'-or $twin.data.reported_state.id-ne $env:SMOKE_DEVICE_ID){throw "Authenticated twin composition failed"}

$oldToken=$env:DEVICE_TOKEN
$oldID=$env:DEVICE_CREDENTIAL_ID
$rotated=Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/credentials/rotate" @{credential_id=$oldID} $headers
$env:DEVICE_TOKEN=$oldToken
Run-IdentityCheck "rejected"
$env:DEVICE_TOKEN=$rotated.data.secret
$env:DEVICE_CREDENTIAL_ID=$rotated.data.credential.credential_id
$env:DEVICE_BOOT_ID="phase2-rotation-$([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds())"
$env:TELEMETRY_SEQUENCE="1"
Run-IdentityCheck "send"

$env:TELEMETRY_SEQUENCE="2"
Run-IdentityCheck "revoke-session"
Run-IdentityCheck "rejected"

$offline=Wait-Connectivity "OFFLINE" 15
$recovery=Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/credentials" @{} $headers
$env:DEVICE_TOKEN=$recovery.data.secret
$env:DEVICE_CREDENTIAL_ID=$recovery.data.credential.credential_id
$env:DEVICE_BOOT_ID="phase2-recovery-$([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds())"
$env:TELEMETRY_SEQUENCE="1"
Run-IdentityCheck "send"
$online=Wait-Connectivity "ONLINE" 10
Run-IdentityCheck "ticket"

Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/suspend" @{} $headers | Out-Null
Run-IdentityCheck "rejected"
Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/activate" @{} $headers | Out-Null
Invoke-API PATCH http://127.0.0.1:6081/api/v1/tenants/alpha_logistics @{status='SUSPENDED'} @{Authorization="Bearer $($env:OPERATOR_TOKEN)"} | Out-Null
Run-IdentityCheck "rejected"
Invoke-API PATCH http://127.0.0.1:6081/api/v1/tenants/alpha_logistics @{status='ACTIVE'} @{Authorization="Bearer $($env:OPERATOR_TOKEN)"} | Out-Null

try { Invoke-API POST http://127.0.0.1:6081/api/v1/tenants @{tenant_id='isolation_tenant';display_name='Isolation Tenant'} @{Authorization="Bearer $($env:OPERATOR_TOKEN)"}|Out-Null } catch { if($_.Exception.Response.StatusCode.value__-ne409){throw} }
$tenantBToken=New-RandomToken "op"
$tenantBPrefix=($tenantBToken.Split('.')[0].Split('_')[-1])
docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -v ON_ERROR_STOP=1 -c "INSERT INTO operator_api_keys(api_key_id,tenant_id,name,token_prefix,token_hash,role,status) VALUES(gen_random_uuid(),'isolation_tenant','isolation test','$tenantBPrefix',digest('$tenantBToken','sha256'),'TENANT_ADMIN','ACTIVE') ON CONFLICT(token_prefix) DO NOTHING" | Out-Null
try { Invoke-API GET "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/twin" $null @{Authorization="Bearer $tenantBToken"}|Out-Null; throw "Cross-tenant twin was exposed" } catch { if($_.Exception.Response.StatusCode.value__-ne404){throw} }

Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$($env:SMOKE_DEVICE_ID)/decommission" @{} $headers | Out-Null
Run-IdentityCheck "rejected"

docker compose -f $composeFile run --rm postgres-migrate | Out-Null
docker compose -f $composeFile run --rm postgres-migrate | Out-Null
$plain=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM device_credentials WHERE encode(token_hash,'escape') LIKE '%pol_dev_%'"
if([int]$plain-ne0){throw "Plaintext credential found"}
$outbox=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM outbox_events WHERE tenant_id='alpha_logistics' AND status='PUBLISHED'"
if([int]$outbox-lt5){throw "Outbox relay did not publish lifecycle events"}
$audit=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM audit_events WHERE tenant_id='alpha_logistics'"
if([int]$audit-lt5){throw "Security mutations were not audited"}
foreach($group in @('polaris_engine_group','polaris_archive_group','polaris_traffic_group')) {
  $description=docker compose -f $composeFile exec -T redpanda rpk group describe $group -X brokers=localhost:29092
  if(($description|Select-String 'TOTAL-LAG\s+0').Count-ne1){throw "$group has non-zero lag"}
}

Write-Host "PASS: Phase 2 authenticated registry, credential lifecycle, tenant isolation, outbox, audit and digital-twin flow"
