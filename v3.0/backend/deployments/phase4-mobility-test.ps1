param([switch]$FullRegression,[switch]$SkipLocalChecks)
$ErrorActionPreference = 'Stop'
$deploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir = Resolve-Path (Join-Path $deploymentDir '..')
$composeFile = Join-Path $deploymentDir 'docker-compose.yml'
$engine = 'http://127.0.0.1:6081/api/v1'

function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null) {
  $params=@{Method=$method;Uri=$uri}
  if($headers){$params.Headers=$headers}
  if($null-ne$body){$params.ContentType='application/json';$params.Body=($body|ConvertTo-Json -Depth 12)}
  Invoke-RestMethod @params
}

if(-not $SkipLocalChecks){Push-Location $backendDir
try { go test ./...; if($LASTEXITCODE-ne0){throw 'Go tests failed'}; go vet ./...; if($LASTEXITCODE-ne0){throw 'Go vet failed'} } finally { Pop-Location }}
& (Join-Path $deploymentDir 'smoke-test.ps1')
$headers=@{Authorization="Bearer $($env:OPERATOR_TOKEN)";'X-Tenant-ID'='alpha_logistics'}

$ready=Invoke-RestMethod http://127.0.0.1:6081/readyz
if($ready.core-ne'READY'-or $ready.modules.mobility.state-notin @('READY','DEGRADED')){throw 'Mobility readiness was not exposed'}
$twin=Invoke-API GET "$engine/devices/$($env:SMOKE_DEVICE_ID)/twin" $null $headers
if(-not $twin.data.components.'spatial/v1'-or-not $twin.data.components.'battery/v1'){throw 'Versioned twin components missing'}
$genericSmokeEnvironment=@{
  SMOKE_DEVICE_ID=$env:SMOKE_DEVICE_ID
  DEVICE_TOKEN=$env:DEVICE_TOKEN
  SMOKE_NODE_TYPE=$env:SMOKE_NODE_TYPE
  SMOKE_LAT=$env:SMOKE_LAT
  SMOKE_LON=$env:SMOKE_LON
}

$route=Invoke-API POST "$engine/routes" @{mobility_profile='ROAD_VEHICLE';origin=@{latitude=13.0067;longitude=80.2206};destination=@{latitude=13.02;longitude=80.23};policy='FASTEST'} $headers
if(-not $route.data.route_id-or -not $route.data.road_graph_version-or $route.data.snapshot_version-lt1-or $route.data.waypoints.Count-lt1){throw 'Versioned A* route was not returned'}

$stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds();$vehicle="P4-ROAD-$stamp"
$project=Invoke-API POST "$engine/projects" @{name="Phase 4 road proof $stamp";description='Mobility planner proof'} $headers
Invoke-API POST "$engine/devices" @{device_id=$vehicle;project_id=$project.data.project_id;device_type_id='connected_vehicle';display_name=$vehicle} $headers|Out-Null
Invoke-API PUT "$engine/devices/$vehicle/capabilities/navigate" @{configuration=@{}} $headers|Out-Null
Invoke-API POST "$engine/devices/$vehicle/activate" @{} $headers|Out-Null
$credential=Invoke-API POST "$engine/devices/$vehicle/credentials" @{} $headers
$env:SMOKE_DEVICE_ID=$vehicle;$env:DEVICE_TOKEN=$credential.data.secret;$env:SMOKE_NODE_TYPE='3';$env:SMOKE_LAT='13.04123';$env:SMOKE_LON='80.23876'
Push-Location $backendDir
try { go run ./cmd/smoke | Out-Host; if($LASTEXITCODE-ne0){throw 'Road telemetry proof failed'} } finally { Pop-Location }
$near=Invoke-API GET "$engine/spatial/devices/nearby?lat=$($env:SMOKE_LAT)&lon=$($env:SMOKE_LON)&radius_meters=25&limit=20" $null $headers
$nearbyDeviceIDs=@($near.data.devices | ForEach-Object { $_.state.device_id })
if($nearbyDeviceIDs -notcontains $vehicle){throw "Mobility nearby query did not return authenticated telemetry device '$vehicle'; returned [$($nearbyDeviceIDs -join ', ')]"}
$task=Invoke-API POST "$engine/tasks" @{task_type='NAVIGATE';priority='HIGH';requirements=@{required_capabilities=@('navigate');minimum_battery=20;max_distance_meters=10000;planning_mode='POLARIS_REQUIRED'};target=@{lat=13.02;lon=80.23;policy='FASTEST'};expires_at=(Get-Date).ToUniversalTime().AddMinutes(5).ToString('o')} $headers
if(-not $task.data.command.payload.route_id-or -not $task.data.command.payload.road_graph_version-or $task.data.command.payload.routing_snapshot_version-lt1){throw 'NAVIGATE did not persist graph and routing snapshot identity'}

foreach($name in $genericSmokeEnvironment.Keys){[Environment]::SetEnvironmentVariable($name,$genericSmokeEnvironment[$name],'Process')}
if($FullRegression){ & (Join-Path $deploymentDir 'reliability-test.ps1'); & (Join-Path $deploymentDir 'phase2-identity-test.ps1'); & (Join-Path $deploymentDir 'phase3-command-test.ps1') }
Write-Host "PASS: Phase 4 H3/R-tree spatial state, A* route snapshot, and durable NAVIGATE planning ($vehicle)"
