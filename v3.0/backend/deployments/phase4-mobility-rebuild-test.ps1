$ErrorActionPreference='Stop'
$deploymentDir=Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir=Resolve-Path (Join-Path $deploymentDir '..')
$rootDir=Resolve-Path (Join-Path $backendDir '..')
$composeFile=Join-Path $deploymentDir 'docker-compose.yml'
$evidenceFile=Join-Path $rootDir 'PHASE_4_2_MOBILITY_REBUILD_RESULT.json'
$engine='http://127.0.0.1:6081/api/v1'

function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null){$p=@{Method=$method;Uri=$uri};if($headers){$p.Headers=$headers};if($null-ne$body){$p.ContentType='application/json';$p.Body=($body|ConvertTo-Json -Depth 10)};Invoke-RestMethod @p}
function New-Tenant([string]$id,[string]$name){try{Invoke-API POST "$engine/tenants" @{tenant_id=$id;display_name=$name} @{Authorization="Bearer $($env:OPERATOR_TOKEN)"}|Out-Null}catch{if($_.Exception.Response.StatusCode.value__-ne409){throw}}}
function New-SpatialDevice([string]$tenant,[string]$project,[string]$id,[string]$type,[int]$nodeType,[double]$lat,[double]$lon){
  $headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'=$tenant}
  Invoke-API POST "$engine/devices" @{device_id=$id;project_id=$project;device_type_id=$type;display_name=$id} $headers|Out-Null
  Invoke-API PUT "$engine/devices/$id/capabilities/navigate" @{configuration=@{}} $headers|Out-Null
  Invoke-API POST "$engine/devices/$id/activate" @{} $headers|Out-Null
  $secret=(Invoke-API POST "$engine/devices/$id/credentials" @{} $headers).data.secret
  $env:SMOKE_TENANT_ID=$tenant;$env:SMOKE_DEVICE_ID=$id;$env:DEVICE_TOKEN=$secret;$env:SMOKE_NODE_TYPE="$nodeType";$env:SMOKE_LAT="$lat";$env:SMOKE_LON="$lon";$env:SMOKE_BOOT_ID="rebuild-$id";$env:SMOKE_SEQUENCE='1';$env:SMOKE_WAIT_FOR_PROJECTION='true';$env:SMOKE_WAIT_FOR_MATCH='false'
  Push-Location $backendDir
  try{go run ./cmd/smoke|Out-Host;if($LASTEXITCODE-ne0){throw "Telemetry failed for $id"}}finally{Pop-Location}
  return $secret
}
function Nearby([string]$tenant,[double]$lat,[double]$lon){$headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'=$tenant};(Invoke-API GET "$engine/spatial/devices/nearby?lat=$lat&lon=$lon&radius_meters=1000&limit=50" $null $headers).data.devices}

. (Join-Path $deploymentDir 'smoke-test.ps1')
$stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds();$tenant="phase42_rebuild_$stamp";$other="phase42_rebuild_other_$stamp"
New-Tenant $tenant 'Phase 4.2 rebuild';New-Tenant $other 'Phase 4.2 rebuild isolation'
$headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'=$tenant};$otherHeaders=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'=$other}
$project=(Invoke-API POST "$engine/projects" @{name='Rebuild devices';description='Derived projection restart proof'} $headers).data.project_id
$otherProject=(Invoke-API POST "$engine/projects" @{name='Isolation devices';description='Rebuild tenancy proof'} $otherHeaders).data.project_id
$road="REBUILD-ROAD-$stamp";$robot="REBUILD-ROBOT-$stamp";$static="REBUILD-STATIC-$stamp";$inactive="REBUILD-INACTIVE-$stamp";$foreign="REBUILD-FOREIGN-$stamp";$compute="REBUILD-COMPUTE-$stamp"
$roadSecret=New-SpatialDevice $tenant $project $road 'connected_vehicle' 3 13.2000 80.3000
New-SpatialDevice $tenant $project $robot 'ground_robot' 6 13.2010 80.3000|Out-Null
New-SpatialDevice $tenant $project $static 'static_camera' 7 13.2020 80.3000|Out-Null
New-SpatialDevice $tenant $project $inactive 'connected_vehicle' 3 13.2030 80.3000|Out-Null
New-SpatialDevice $other $otherProject $foreign 'connected_vehicle' 3 13.2005 80.3000|Out-Null
Invoke-API POST "$engine/devices" @{device_id=$compute;project_id=$project;device_type_id='compute_node';display_name=$compute} $headers|Out-Null
Invoke-API POST "$engine/devices/$compute/activate" @{} $headers|Out-Null
Invoke-API POST "$engine/devices/$inactive/suspend" @{} $headers|Out-Null

# Advance the road device, then inject a stale sequence without waiting for a
# projection event. Redis and Mobility must both retain sequence 2.
$env:SMOKE_TENANT_ID=$tenant;$env:SMOKE_DEVICE_ID=$road;$env:DEVICE_TOKEN=$roadSecret;$env:SMOKE_NODE_TYPE='3';$env:SMOKE_LAT='13.2002';$env:SMOKE_LON='80.3001';$env:SMOKE_BOOT_ID="rebuild-$road";$env:SMOKE_SEQUENCE='2';$env:SMOKE_WAIT_FOR_PROJECTION='true';$env:SMOKE_WAIT_FOR_MATCH='false'
Push-Location $backendDir;try{go run ./cmd/smoke|Out-Host;if($LASTEXITCODE-ne0){throw 'Newer road telemetry failed'}}finally{Pop-Location}
$env:SMOKE_LAT='13.2500';$env:SMOKE_LON='80.3500';$env:SMOKE_SEQUENCE='1';$env:SMOKE_WAIT_FOR_PROJECTION='false'
Push-Location $backendDir;try{go run ./cmd/smoke|Out-Host;if($LASTEXITCODE-ne0){throw 'Stale road telemetry injection failed'}}finally{Pop-Location}

$before=Nearby $tenant 13.2010 80.3000;$beforeRoad=$before|Where-Object{$_.state.device_id-eq$road}|Select-Object -First 1
if((-not $beforeRoad)-or[uint64]$beforeRoad.state.sequence_number-ne2){throw 'Pre-restart Mobility state did not retain the newer sequence'}
docker compose -f $composeFile restart engine|Out-Null
$sawNotReady=$false;$deadline=(Get-Date).AddMinutes(3);$ready=$null
do{
  try{$probe=Invoke-RestMethod http://127.0.0.1:6081/readyz;if($probe.modules.mobility.state-eq'READY'){$ready=$probe;break}else{$sawNotReady=$true}}catch{$sawNotReady=$true}
  Start-Sleep -Milliseconds 250
}while((Get-Date)-lt$deadline)
if((-not $sawNotReady)-or(-not $ready)){throw 'Readiness did not transition from unavailable/not-ready to Mobility READY during rebuild'}
$after=Nearby $tenant 13.2010 80.3000;$ids=@($after|ForEach-Object{$_.state.device_id});$afterRoad=$after|Where-Object{$_.state.device_id-eq$road}|Select-Object -First 1
foreach($id in @($road,$robot,$static)){if(($ids|Where-Object{$_-eq$id}).Count-ne1){throw "Active rebuilt device $id was missing or duplicated"}}
foreach($id in @($inactive,$foreign,$compute)){if($ids-contains$id){throw "Excluded device $id leaked into rebuilt spatial state"}}
if([uint64]$afterRoad.state.sequence_number-ne2-or[double]$afterRoad.state.reported_position.latitude-ne13.2002){throw 'Rebuild regressed the accepted road boot/sequence/position'}
$otherView=Nearby $other 13.2010 80.3000;$otherIDs=@($otherView|ForEach-Object{$_.state.device_id})
if($otherIDs-notcontains$foreign-or$otherIDs-contains$road){throw 'Rebuild tenant isolation failed'}
$result=@{measured_at=(Get-Date).ToUniversalTime().ToString('o');readiness_transition_observed=$sawNotReady;mobility_state=$ready.modules.mobility.state;recovered=@($road,$robot,$static);excluded=@($inactive,$compute,$foreign);road_sequence_before=[uint64]$beforeRoad.state.sequence_number;road_sequence_after=[uint64]$afterRoad.state.sequence_number;duplicates=0;cross_tenant=0}|ConvertTo-Json -Depth 10
$result|Set-Content -LiteralPath $evidenceFile -Encoding UTF8
Write-Host "PASS: Mobility restart/rebuild, freshness, exclusion, uniqueness, and tenancy; evidence: $evidenceFile"
