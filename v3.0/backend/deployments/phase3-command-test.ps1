$ErrorActionPreference = "Stop"
$deploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir = Resolve-Path (Join-Path $deploymentDir "..")
$composeFile = Join-Path $deploymentDir "docker-compose.yml"
$engine = "http://127.0.0.1:6081/api/v1"

function New-RandomToken([string]$kind) {
  $rng=[Security.Cryptography.RandomNumberGenerator]::Create()
  try { $prefixBytes=New-Object byte[] 8; $rng.GetBytes($prefixBytes); $secretBytes=New-Object byte[] 32; $rng.GetBytes($secretBytes) } finally { $rng.Dispose() }
  $prefix=([BitConverter]::ToString($prefixBytes)-replace '-','').ToLowerInvariant()
  $secret=([BitConverter]::ToString($secretBytes)-replace '-','').ToLowerInvariant()
  return "pol_${kind}_${prefix}.${secret}"
}
function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null) {
  $params=@{Method=$method;Uri=$uri}
  if($headers){$params.Headers=$headers}
  if($null-ne$body){$params.ContentType='application/json';$params.Body=($body|ConvertTo-Json -Depth 10)}
  return Invoke-RestMethod @params
}
function Run-Check([string]$mode) {
  $env:ORCHESTRATION_CHECK_MODE=$mode
  Push-Location $backendDir
  try { $output=go run ./cmd/orchestrationcheck } finally { Pop-Location }
  if($LASTEXITCODE-ne0){throw "Phase 3 client failed: $mode"}
  $output | ForEach-Object { Write-Host $_ }
  return $output
}
function New-Device([string]$label,[string[]]$capabilities) {
  $stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
  $id="P3-$label-$stamp"
  $project=Invoke-API POST "$engine/projects" @{name="Phase 3 $label $stamp";description='Durable command proof'} $script:headers
  Invoke-API POST "$engine/devices" @{device_id=$id;project_id=$project.data.project_id;device_type_id='delivery_drone';display_name=$id} $script:headers|Out-Null
  foreach($capability in $capabilities){Invoke-API PUT "$engine/devices/$id/capabilities/$capability" @{configuration=@{}} $script:headers|Out-Null}
  Invoke-API POST "$engine/devices/$id/activate" @{} $script:headers|Out-Null
  $credential=Invoke-API POST "$engine/devices/$id/credentials" @{} $script:headers
  return [pscustomobject]@{ID=$id;Token=$credential.data.secret;CredentialID=$credential.data.credential.credential_id;ProjectID=$project.data.project_id}
}
function Select-Device($device) {
  $env:SMOKE_DEVICE_ID=$device.ID
  $env:DEVICE_TOKEN=$device.Token
  $env:DEVICE_CREDENTIAL_ID=$device.CredentialID
  $env:TASK_PROJECT_ID=$device.ProjectID
}
function Wait-CommandStatus([string]$commandID,[string]$expected,[int]$seconds=20) {
  $deadline=(Get-Date).AddSeconds($seconds)
  do{$value=Invoke-API GET "$engine/commands/$commandID" $null $script:headers;if($value.data.status-eq$expected){return $value.data};Start-Sleep -Milliseconds 200}while((Get-Date)-lt$deadline)
  throw "Command $commandID expected $expected, got $($value.data.status)"
}
function Wait-Online([string]$deviceID,[int]$seconds=15) {
  $deadline=(Get-Date).AddSeconds($seconds)
  do{$twin=Invoke-API GET "$engine/devices/$deviceID/twin" $null $script:headers;if($twin.data.connectivity.status-eq'ONLINE'){return};Start-Sleep -Milliseconds 200}while((Get-Date)-lt$deadline)
  throw "Device $deviceID did not become ONLINE before failure injection"
}

$env:COMMAND_ACK_TIMEOUT='1s'
$env:COMMAND_RECONCILE_INTERVAL='200ms'
$env:OUTBOX_POLL_INTERVAL='100ms'
$env:COMMAND_MAX_ATTEMPTS='3'
$env:CONNECTION_LEASE_TTL='6s'
$env:DEVICE_STALE_AFTER='30s'
$env:DEVICE_OFFLINE_AFTER='60s'
$env:OFFLINE_SCAN_INTERVAL='1s'

& (Join-Path $deploymentDir 'smoke-test.ps1')
$script:headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'='alpha_logistics'}
$smokeProject=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT project_id FROM devices WHERE tenant_id='alpha_logistics' AND device_id='$($env:SMOKE_DEVICE_ID)'"
$env:TASK_PROJECT_ID=$smokeProject.Trim()
Run-Check 'complete'|Out-Null

$duplicate=New-Device 'DUPLICATE' @('navigate','receive_relocation_command')
Select-Device $duplicate
Run-Check 'duplicate'|Out-Null

$offline=New-Device 'OFFLINE' @('navigate','receive_relocation_command')
Select-Device $offline
Run-Check 'offline'|Out-Null

$fenced=New-Device 'FENCE' @('navigate','receive_relocation_command')
Select-Device $fenced
$epochBefore=docker compose -f $composeFile exec -T redis redis-cli GET "polaris:connection-epoch:alpha_logistics:$($fenced.ID)"
Run-Check 'fencing'|Out-Null
$epochAfter=docker compose -f $composeFile exec -T redis redis-cli GET "polaris:connection-epoch:alpha_logistics:$($fenced.ID)"
if([int]$epochAfter-le[int]$epochBefore){throw 'Ownership fencing epoch did not advance'}

$wrongA=New-Device 'ACK-A' @('navigate','receive_relocation_command')
$wrongB=New-Device 'ACK-B' @('navigate','receive_relocation_command')
Select-Device $wrongA
$env:DEVICE_ID_B=$wrongB.ID
$env:DEVICE_TOKEN_B=$wrongB.Token
Run-Check 'wrong-ack'|Out-Null

$mismatch=New-Device 'NO-CAMERA' @('navigate','receive_relocation_command')
Select-Device $mismatch
Run-Check 'capability-mismatch'|Out-Null

$cancelTask=Invoke-API POST "$engine/tasks" @{project_id=$mismatch.ProjectID;task_type='CAPTURE_IMAGE';priority='NORMAL';requirements=@{required_capabilities=@('capture_image');project_id=$mismatch.ProjectID};target=@{lat=13.0067;lon=80.2206};expires_at=(Get-Date).ToUniversalTime().AddMinutes(1).ToString('o')} $script:headers
Invoke-API POST "$engine/tasks/$($cancelTask.data.task.task_id)/cancel" @{} $script:headers|Out-Null
$cancelled=Invoke-API GET "$engine/tasks/$($cancelTask.data.task.task_id)" $null $script:headers
if($cancelled.data.task.status-ne'CANCELLED'){throw 'Pending task cancellation failed'}

$crash=New-Device 'CRASH' @('navigate','receive_relocation_command')
Select-Device $crash
$receiveOutput=Run-Check 'receive-no-ack'
$commandLine=$receiveOutput|Where-Object{$_-like'COMMAND_ID=*'}|Select-Object -First 1
if(-not$commandLine){throw 'Gateway crash scenario did not expose command ID'}
$env:EXPECTED_COMMAND_ID=$commandLine.Substring('COMMAND_ID='.Length)
docker compose -f $composeFile restart gateway|Out-Null
docker compose -f $composeFile up -d --wait gateway|Out-Null
Run-Check 'resume'|Out-Null

# Commit command state while Kafka is unavailable, then prove outbox recovery.
$recovery=New-Device 'OUTBOX' @('navigate','receive_relocation_command')
Select-Device $recovery
$env:IDENTITY_CHECK_MODE='send'
$env:DEVICE_BOOT_ID="phase3-outbox-$([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds())"
Push-Location $backendDir;try{go run ./cmd/identitycheck|Out-Null}finally{Pop-Location}
Wait-Online $recovery.ID
docker compose -f $composeFile stop redpanda|Out-Null
$recoveryTask=Invoke-API POST "$engine/tasks" @{project_id=$recovery.ProjectID;task_type='RELOCATE';priority='HIGH';requirements=@{required_capabilities=@('receive_relocation_command');minimum_battery=30;project_id=$recovery.ProjectID};target=@{lat=13.0068;lon=80.2207};expires_at=(Get-Date).ToUniversalTime().AddMinutes(2).ToString('o')} $script:headers
$recoveryCommand=$recoveryTask.data.command.command_id
if(-not$recoveryCommand){throw 'Command was not committed while Kafka was unavailable'}
$pendingOutbox=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM outbox_events WHERE aggregate_id='$recoveryCommand' AND status<>'PUBLISHED'"
if([int]$pendingOutbox-lt1){throw 'Command outbox was not retained during Kafka outage'}
docker compose -f $composeFile start redpanda|Out-Null
$deadline=(Get-Date).AddSeconds(60)
do{$published=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM outbox_events WHERE aggregate_id='$recoveryCommand' AND event_type='command.created.v1' AND status='PUBLISHED'";if([int]$published-ge1){break};Start-Sleep -Milliseconds 500}while((Get-Date)-lt$deadline)
if([int]$published-lt1){throw 'Outbox did not recover after Kafka restart'}

# Replaying the durable command event must not create another command row.
docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -c "UPDATE outbox_events SET status='PENDING',published_at=NULL,next_attempt_at=NOW() WHERE aggregate_id='$recoveryCommand' AND event_type='command.created.v1'"|Out-Null
Start-Sleep -Seconds 2
$duplicates=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tAc "SELECT count(*) FROM commands WHERE command_id='$recoveryCommand'"
if([int]$duplicates-ne1){throw 'Kafka/outbox replay duplicated a command row'}

# Cross-tenant task reads are hidden and viewer mutation is forbidden.
try{Invoke-API POST "$engine/tenants" @{tenant_id='phase3_isolation';display_name='Phase 3 isolation'} @{Authorization="Bearer $($env:OPERATOR_TOKEN)"}|Out-Null}catch{if($_.Exception.Response.StatusCode.value__-ne409){throw}}
$tenantToken=New-RandomToken 'op';$tenantPrefix=($tenantToken.Split('.')[0].Split('_')[-1])
$viewerToken=New-RandomToken 'op';$viewerPrefix=($viewerToken.Split('.')[0].Split('_')[-1])
docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -v ON_ERROR_STOP=1 -c "INSERT INTO operator_api_keys(api_key_id,tenant_id,name,token_prefix,token_hash,role,status) VALUES(gen_random_uuid(),'phase3_isolation','phase3 tenant','$tenantPrefix',digest('$tenantToken','sha256'),'TENANT_ADMIN','ACTIVE'),(gen_random_uuid(),'alpha_logistics','phase3 viewer','$viewerPrefix',digest('$viewerToken','sha256'),'VIEWER','ACTIVE') ON CONFLICT(token_prefix) DO NOTHING"|Out-Null
try{Invoke-API GET "$engine/tasks/$($recoveryTask.data.task.task_id)" $null @{Authorization="Bearer $tenantToken"}|Out-Null;throw 'Cross-tenant task was exposed'}catch{if($_.Exception.Response.StatusCode.value__-ne404){throw}}
try{Invoke-API POST "$engine/tasks" @{task_type='RELOCATE';target=@{lat=13;lon=80};expires_at=(Get-Date).ToUniversalTime().AddMinutes(1).ToString('o')} @{Authorization="Bearer $viewerToken"}|Out-Null;throw 'Viewer created a task'}catch{if($_.Exception.Response.StatusCode.value__-ne403){throw}}

foreach($group in @('polaris_engine_group','polaris_archive_group','polaris_traffic_group','polaris-command-dispatcher')){
  $deadline=(Get-Date).AddSeconds(20);do{$description=docker compose -f $composeFile exec -T redpanda rpk group describe $group -X brokers=localhost:29092;if(($description|Select-String 'TOTAL-LAG\s+0').Count-eq1){break};Start-Sleep -Milliseconds 500}while((Get-Date)-lt$deadline)
  if(($description|Select-String 'TOTAL-LAG\s+0').Count-ne1){throw "$group has non-zero lag"}
}

Write-Host 'PASS: Phase 3 durable task assignment, fenced delivery, idempotent ACK/result, retry, expiry, recovery, RBAC and tenancy flow'
