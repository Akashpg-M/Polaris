$ErrorActionPreference='Stop'
$deploymentDir=Split-Path -Parent $MyInvocation.MyCommand.Path
$backendDir=Resolve-Path (Join-Path $deploymentDir '..')
$rootDir=Resolve-Path (Join-Path $backendDir '..')
$composeFile=Join-Path $deploymentDir 'docker-compose.yml'
$evidenceFile=Join-Path $rootDir 'PHASE_4_2_ROUTING_OVERLOAD_RESULT.json'
$stamp=[DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds();$tenant="phase42_overload_$stamp";$device="OVERLOAD-$stamp"

function New-Token([string]$kind){$rng=[Security.Cryptography.RandomNumberGenerator]::Create();try{$a=New-Object byte[] 8;$b=New-Object byte[] 32;$rng.GetBytes($a);$rng.GetBytes($b)}finally{$rng.Dispose()};"pol_${kind}_$(([BitConverter]::ToString($a)-replace '-','').ToLowerInvariant()).$(([BitConverter]::ToString($b)-replace '-','').ToLowerInvariant())"}
function Invoke-API([string]$method,[string]$uri,$body=$null,$headers=$null){$p=@{Method=$method;Uri=$uri};if($headers){$p.Headers=$headers};if($null-ne$body){$p.ContentType='application/json';$p.Body=($body|ConvertTo-Json -Depth 10)};Invoke-RestMethod @p}

$saved=@{workers=$env:MOBILITY_ROUTING_WORKERS;queue=$env:MOBILITY_ROUTING_QUEUE_CAPACITY;tenant=$env:MOBILITY_MAX_CONCURRENT_ROUTES_PER_TENANT;admin=$env:DEV_PLATFORM_ADMIN_TOKEN}
try{
  $env:MOBILITY_ROUTING_WORKERS='1';$env:MOBILITY_ROUTING_QUEUE_CAPACITY='4';$env:MOBILITY_MAX_CONCURRENT_ROUTES_PER_TENANT='128';$env:DEV_PLATFORM_ADMIN_TOKEN=New-Token 'op'
  docker compose -f $composeFile up -d --build --wait --force-recreate
  if($LASTEXITCODE-ne0){throw 'Overload stack failed to become ready'}
  $adminHeaders=@{Authorization="Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)"}
  Invoke-API POST 'http://127.0.0.1:6081/api/v1/tenants' @{tenant_id=$tenant;display_name='Phase 4.2 overload'} $adminHeaders|Out-Null
  $headers=@{Authorization="Bearer $($env:DEV_PLATFORM_ADMIN_TOKEN)";'X-Tenant-ID'=$tenant}
  $project=(Invoke-API POST 'http://127.0.0.1:6081/api/v1/projects' @{name='Overload isolation';description='Routing must not block core'} $headers).data.project_id
  Invoke-API POST 'http://127.0.0.1:6081/api/v1/devices' @{device_id=$device;project_id=$project;device_type_id='connected_vehicle';display_name=$device} $headers|Out-Null
  Invoke-API PUT "http://127.0.0.1:6081/api/v1/devices/$device/capabilities/run_model" @{configuration=@{}} $headers|Out-Null
  Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$device/activate" @{} $headers|Out-Null
  $secret=(Invoke-API POST "http://127.0.0.1:6081/api/v1/devices/$device/credentials" @{} $headers).data.secret
  Push-Location $backendDir
  try{$result=go run ./cmd/routing-overload -admin-token $env:DEV_PLATFORM_ADMIN_TOKEN -tenant $tenant -project $project -device $device -device-token $secret;if($LASTEXITCODE-ne0){throw 'Routing overload harness failed'}}finally{Pop-Location}
  $result|Set-Content -LiteralPath $evidenceFile -Encoding UTF8
  Write-Host "PASS: bounded routing overload, unrelated telemetry/task delivery, and recovery; evidence: $evidenceFile"
}finally{
  $env:MOBILITY_ROUTING_WORKERS=$saved.workers;$env:MOBILITY_ROUTING_QUEUE_CAPACITY=$saved.queue;$env:MOBILITY_MAX_CONCURRENT_ROUTES_PER_TENANT=$saved.tenant;$env:DEV_PLATFORM_ADMIN_TOKEN=$saved.admin
  docker compose -f $composeFile up -d --wait --force-recreate engine|Out-Null
}
