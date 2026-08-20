$ErrorActionPreference='Stop'
$deploymentDir=Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir=Resolve-Path (Join-Path $deploymentDir '..')
$composeFile=Join-Path $deploymentDir 'docker-compose.yml'
$engine='http://127.0.0.1:6081/api/v1'
function New-Token(){ $rng=[Security.Cryptography.RandomNumberGenerator]::Create();try{$a=New-Object byte[] 8;$b=New-Object byte[] 32;$rng.GetBytes($a);$rng.GetBytes($b)}finally{$rng.Dispose()};"pol_op_$(([BitConverter]::ToString($a)-replace '-','').ToLowerInvariant()).$(([BitConverter]::ToString($b)-replace '-','').ToLowerInvariant())" }
function API([string]$method,[string]$uri,$body=$null,$headers=$null){$p=@{Method=$method;Uri=$uri};if($headers){$p.Headers=$headers};if($null-ne$body){$p.ContentType='application/json';$p.Body=($body|ConvertTo-Json -Depth 10)};Invoke-RestMethod @p}
$env:DEV_PLATFORM_ADMIN_TOKEN=New-Token;$env:POLARIS_MODULE_MOBILITY_ENABLED='false'
try{
  docker compose -f $composeFile up -d --build --force-recreate --wait engine gateway
  if($LASTEXITCODE-ne0){throw 'Mobility-disabled engine failed readiness'}
  $ready=Invoke-RestMethod http://127.0.0.1:6081/readyz
  if($ready.core-ne'READY'-or $ready.modules.mobility){throw 'Mobility was not disabled cleanly'}
  $headers=@{Authorization="Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)";'X-Tenant-ID'='alpha_logistics'}
  try{API POST "$engine/tenants" @{tenant_id='alpha_logistics';display_name='Alpha Logistics'} @{Authorization="Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)"}|Out-Null}catch{if($_.Exception.Response.StatusCode.value__-ne409){throw}}
  $stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds();$device="P4-CAMERA-$stamp"
  $project=API POST "$engine/projects" @{name="Phase 4 generic proof $stamp"} $headers
  API POST "$engine/devices" @{device_id=$device;project_id=$project.data.project_id;device_type_id='static_camera';display_name=$device} $headers|Out-Null
  API PUT "$engine/devices/$device/capabilities/capture_image" @{configuration=@{}} $headers|Out-Null
  API POST "$engine/devices/$device/activate" @{} $headers|Out-Null
  $credential=API POST "$engine/devices/$device/credentials" @{} $headers
  $env:OPERATOR_TOKEN=$env:DEV_PLATFORM_ADMIN_TOKEN;$env:DEVICE_TOKEN=$credential.data.secret;$env:SMOKE_DEVICE_ID=$device;$env:SMOKE_NODE_TYPE='7'
  Push-Location $backendDir;try{go run ./cmd/smoke|Out-Host;if($LASTEXITCODE-ne0){throw 'Generic telemetry failed with Mobility disabled'}}finally{Pop-Location}
  $task=API POST "$engine/tasks" @{project_id=$project.data.project_id;task_type='CAPTURE_IMAGE';priority='NORMAL';requirements=@{required_capabilities=@('capture_image');project_id=$project.data.project_id};target=@{image_profile='overview'};expires_at=(Get-Date).ToUniversalTime().AddMinutes(2).ToString('o')} $headers
  if(-not $task.data.command.command_id-or $task.data.command.payload.image_profile-ne'overview'){throw 'Default planner failed with Mobility disabled'}
  Write-Host "PASS: Mobility disabled; telemetry, twin state, CAPTURE_IMAGE and durable generic command remain functional ($device)"
}finally{
  $env:POLARIS_MODULE_MOBILITY_ENABLED='true';docker compose -f $composeFile up -d --force-recreate --wait engine|Out-Null
}
