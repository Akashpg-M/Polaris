param(
  [int]$Devices = 1000,
  [int]$DurationSeconds = 45,
  [int]$Tasks = 120,
  [int]$RampPerSecond = 25,
  [int]$TelemetryIntervalSeconds = 5,
  [string]$EvidenceName = 'PHASE_4_2',
  [string]$AdminToken = '',
  [switch]$SkipCompose
)
$ErrorActionPreference = 'Stop'
if($Devices -lt 20){throw 'At least 20 devices are required for a heterogeneous soak'}
if($EvidenceName-notmatch'^[A-Za-z0-9_-]+$'){throw 'EvidenceName may contain only letters, digits, underscore, and dash'}
$deploymentDir=Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir=Resolve-Path (Join-Path $deploymentDir '..')
$rootDir=Resolve-Path (Join-Path $backendDir '..')
$composeFile=Join-Path $deploymentDir 'docker-compose.yml'
$stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
$tenant="phase42_soak_$stamp";$isolationTenant="phase42_isolation_$stamp";$project=[guid]::NewGuid().ToString()
$credentialFile=Join-Path $env:TEMP "polaris-phase42-devices-$stamp.json"
$sqlFile=Join-Path $env:TEMP "polaris-phase42-seed-$stamp.sql"
$stdoutFile=Join-Path $env:TEMP "polaris-phase42-output-$stamp.json"
$stderrFile=Join-Path $env:TEMP "polaris-phase42-error-$stamp.log"
$statsFile=Join-Path $rootDir "${EvidenceName}_CONTAINER_STATS.jsonl"
$resultFile=Join-Path $rootDir "${EvidenceName}_SOAK_RESULT.json"

function New-Token([string]$kind){
  $rng=[Security.Cryptography.RandomNumberGenerator]::Create()
  try{$a=New-Object byte[] 8;$b=New-Object byte[] 32;$rng.GetBytes($a);$rng.GetBytes($b)}finally{$rng.Dispose()}
  "pol_${kind}_$(([BitConverter]::ToString($a)-replace '-','').ToLowerInvariant()).$(([BitConverter]::ToString($b)-replace '-','').ToLowerInvariant())"
}
function Escape-SQL([string]$value){$value.Replace("'","''")}
function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null){
  $p=@{Method=$method;Uri=$uri};if($headers){$p.Headers=$headers};if($null-ne$body){$p.ContentType='application/json';$p.Body=($body|ConvertTo-Json -Depth 10)};Invoke-RestMethod @p
}

$env:DEV_PLATFORM_ADMIN_TOKEN=if($AdminToken){$AdminToken}else{New-Token 'op'}
try{
  # Recreate processes and Docker port forwarding for a reproducible proof;
  # named Kafka/PostgreSQL/Redis volumes remain durable.
  if(-not $SkipCompose){
    docker compose -f $composeFile up -d --build --wait --force-recreate
    if($LASTEXITCODE-ne0){throw 'Compose stack did not become ready'}
  }else{
    $engineReady=Invoke-RestMethod http://127.0.0.1:6081/readyz;$gatewayReady=Invoke-RestMethod http://127.0.0.1:6080/readyz
    if($engineReady.status-ne'ready'-or$gatewayReady.status-ne'ready'){throw 'SkipCompose requires an already-ready stack'}
  }
  $records=New-Object System.Collections.Generic.List[object]
  $sql=New-Object System.Text.StringBuilder
  [void]$sql.AppendLine('BEGIN;')
  [void]$sql.AppendLine("INSERT INTO tenants(tenant_id,display_name,status) VALUES('$tenant','Phase 4.2 soak','ACTIVE'),('$isolationTenant','Phase 4.2 isolation','ACTIVE') ON CONFLICT DO NOTHING;")
  [void]$sql.AppendLine("INSERT INTO projects(project_id,tenant_id,name,status) VALUES('$project','$tenant','Phase 4.2 mixed soak $stamp','ACTIVE') ON CONFLICT DO NOTHING;")
  for($i=0;$i-lt$Devices;$i++){
    $ratio=$i/[double]$Devices
    if($ratio-lt.40){$deviceType='connected_vehicle';$nodeType=3;$spatial=$true;$caps=@('navigate','receive_relocation_command','run_model')}
    elseif($ratio-lt.60){$deviceType='ground_robot';$nodeType=6;$spatial=$true;$caps=@('navigate','receive_relocation_command')}
    elseif($ratio-lt.85){$deviceType='static_camera';$nodeType=7;$spatial=$true;$caps=@('capture_image')}
    else{$deviceType='compute_node';$nodeType=0;$spatial=$false;$caps=@('run_model')}
    $id="SOAK-$stamp-$('{0:D4}'-f$i)";$token=New-Token 'dev';$prefix=($token.Split('.')[0] -split '_')[-1];$credential=[guid]::NewGuid().ToString()
    [void]$sql.AppendLine("INSERT INTO devices(tenant_id,device_id,project_id,device_type_id,display_name,lifecycle_status) VALUES('$tenant','$id','$project','$deviceType','$id','ACTIVE') ON CONFLICT DO NOTHING;")
    [void]$sql.AppendLine("INSERT INTO device_credentials(credential_id,tenant_id,device_id,token_prefix,token_hash,status) VALUES('$credential','$tenant','$id','$prefix',digest('$(Escape-SQL $token)','sha256'),'ACTIVE') ON CONFLICT(token_prefix) DO NOTHING;")
    foreach($cap in $caps){[void]$sql.AppendLine("INSERT INTO device_capabilities(tenant_id,device_id,capability_id,configuration,enabled) VALUES('$tenant','$id','$cap','{}',true) ON CONFLICT(tenant_id,device_id,capability_id) DO UPDATE SET enabled=true;")}
    $records.Add(@{tenant_id=$tenant;device_id=$id;token=$token;node_type=$nodeType;spatial=$spatial})
  }
  [void]$sql.AppendLine('COMMIT;')
  Set-Content -LiteralPath $sqlFile -Value $sql.ToString() -Encoding UTF8
  $records|ConvertTo-Json -Depth 4|Set-Content -LiteralPath $credentialFile -Encoding UTF8
  Get-Content -Raw -LiteralPath $sqlFile|docker compose -f $composeFile exec -T postgres psql -q -U polaris_user -d polaris_core -v ON_ERROR_STOP=1|Out-Null
  if($LASTEXITCODE-ne0){throw 'Soak registry seed failed'}
  # This harness bulk-loads in one transaction, unlike normal incremental
  # registration. Refresh statistics so the measured query plan represents
  # the populated fleet instead of PostgreSQL's pre-seed cardinality estimate.
  docker compose -f $composeFile exec -T postgres psql -q -U polaris_user -d polaris_core -c 'ANALYZE devices; ANALYZE device_capabilities; ANALYZE device_assignments;'|Out-Null
  if($LASTEXITCODE-ne0){throw 'Soak planner statistics refresh failed'}

  Set-Content -LiteralPath $statsFile -Value '' -Encoding UTF8
  $arguments=@('run','./cmd/system-soak','-devices',$credentialFile,'-admin-token',$env:DEV_PLATFORM_ADMIN_TOKEN,'-tenant',$tenant,'-project',$project,'-duration',"${DurationSeconds}s",'-tasks',"$Tasks",'-ramp-per-second',"$RampPerSecond",'-telemetry-interval',"${TelemetryIntervalSeconds}s")
  $process=Start-Process -FilePath 'go' -ArgumentList $arguments -WorkingDirectory $backendDir -RedirectStandardOutput $stdoutFile -RedirectStandardError $stderrFile -WindowStyle Hidden -PassThru
  while(-not $process.HasExited){
    try{$engineReady=Invoke-RestMethod http://127.0.0.1:6081/readyz}catch{$engineReady=@{status='unavailable';error=$_.Exception.Message}}
    try{$gatewayReady=Invoke-RestMethod http://127.0.0.1:6080/readyz}catch{$gatewayReady=@{status='unavailable';error=$_.Exception.Message}}
    $sample=@{captured_at=(Get-Date).ToUniversalTime().ToString('o');containers=(docker stats --no-stream --format '{{json .}}');engine=$engineReady;gateway=$gatewayReady}|ConvertTo-Json -Compress -Depth 12
    Add-Content -LiteralPath $statsFile -Value $sample
    Start-Sleep -Seconds 5
    $process.Refresh()
  }
  $process.WaitForExit()
  if(-not(Test-Path -LiteralPath $stdoutFile)-or(Get-Item -LiteralPath $stdoutFile).Length-eq0){
    $detail=if(Test-Path -LiteralPath $stderrFile){Get-Content -Raw $stderrFile}else{'no process output'}
    throw "System soak produced no evidence: $detail"
  }
  Copy-Item -LiteralPath $stdoutFile -Destination $resultFile -Force
  $soak=Get-Content -Raw -LiteralPath $resultFile|ConvertFrom-Json
  if([int]$soak.counters.connections_established-ne$Devices-or[int]$soak.counters.connection_errors-ne0-or[int]$soak.counters.identity_mutations-ne0-or[int]$soak.counters.duplicate_physical_executions-ne0){
    $detail=if(Test-Path -LiteralPath $stderrFile){Get-Content -Raw $stderrFile}else{''}
    throw "System soak connection/identity invariant failed: $detail; evidence: $resultFile"
  }
  if([int]$soak.error_totals.unexpected-ne0-or[int]$soak.error_totals.server_error-ne0-or[int]$soak.error_totals.transport_error-ne0){throw "Unexpected classified workload errors: $($soak.error_totals|ConvertTo-Json -Compress); evidence: $resultFile"}
  if([int64]$soak.counters.telemetry_sent-lt1-or[int]$soak.counters.task_requests_attempted-ne$Tasks-or[int]$soak.counters.physical_executions-ne$Tasks-or[int]$soak.counters.commands_delivered-lt$Tasks){throw "System soak workload did not exercise every required path; evidence: $resultFile"}

  foreach($group in @('polaris_engine_group','polaris_archive_group','polaris_traffic_group','polaris-command-dispatcher')){
    $deadline=(Get-Date).AddSeconds(90)
    do{$description=docker compose -f $composeFile exec -T redpanda rpk group describe $group -X brokers=localhost:29092;if(($description|Select-String 'TOTAL-LAG\s+0').Count-eq1){break};Start-Sleep -Milliseconds 500}while((Get-Date)-lt$deadline)
    if(($description|Select-String 'TOTAL-LAG\s+0').Count-ne1){throw "$group lag did not return to zero"}
  }
  $invariants=docker compose -f $composeFile exec -T postgres psql -U polaris_user -d polaris_core -tA -F ',' -c "SELECT (SELECT count(*) FROM telemetry_history WHERE tenant_id='$tenant'),(SELECT count(*) FROM (SELECT tenant_id,device_id,count(*) FROM device_assignments WHERE tenant_id='$tenant' AND status='ACTIVE' GROUP BY tenant_id,device_id HAVING count(*)>1)x),(SELECT count(*) FROM commands c JOIN tasks t USING(task_id) WHERE c.tenant_id='$tenant' AND (c.tenant_id<>t.tenant_id OR c.device_id<>t.assigned_device_id)),(SELECT count(*) FROM commands c JOIN tasks t USING(task_id) WHERE c.tenant_id='$tenant' AND t.requirements->>'planning_mode'='POLARIS_REQUIRED' AND (NOT(c.payload?'road_graph_version') OR NOT(c.payload?'routing_snapshot_version'))),(SELECT count(*) FROM tasks WHERE tenant_id='$tenant');"
  $parts=($invariants|Select-Object -Last 1).Trim().Split(',')
  if([int64]$parts[0]-lt1-or[int]$parts[1]-ne0-or[int]$parts[2]-ne0-or[int]$parts[3]-ne0-or[int]$parts[4]-lt$Tasks){throw "Database invariant failure: $invariants"}
  $headers=@{Authorization="Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)";'X-Tenant-ID'=$isolationTenant}
  $isolated=Invoke-API GET 'http://127.0.0.1:6081/api/v1/spatial/devices/nearby?lat=13.0067&lon=80.2206&radius_meters=5000&limit=20' $null $headers
  if($isolated.data.count-ne0){throw 'Cross-tenant spatial leakage detected'}
  $ready=Invoke-RestMethod http://127.0.0.1:6081/readyz
  if($ready.modules.mobility.details.routing_runtime.queue_depth-ge$ready.modules.mobility.details.routing_runtime.queue_capacity){throw 'Routing queue remained saturated'}
  $soak|Add-Member -NotePropertyName final_runtime -NotePropertyValue @{engine=$ready;gateway=(Invoke-RestMethod http://127.0.0.1:6080/readyz)} -Force
  $soak|ConvertTo-Json -Depth 20|Set-Content -LiteralPath $resultFile -Encoding UTF8
  Write-Host "PASS: $Devices-device heterogeneous full-system soak; evidence: $resultFile and $statsFile"
}finally{
  foreach($path in @($credentialFile,$sqlFile,$stdoutFile,$stderrFile)){if(Test-Path -LiteralPath $path){Remove-Item -LiteralPath $path -Force}}
}
