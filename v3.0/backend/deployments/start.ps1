param(
  [switch]$Pull,
  [switch]$NoBuild,
  [int]$WaitTimeoutSeconds = 300
)
. (Join-Path $PSScriptRoot 'deployment-common.ps1')
Initialize-DeploymentEnvironment
Assert-DockerAvailable
$compose = Get-ComposeArguments

& docker @compose config --quiet
if ($LASTEXITCODE -ne 0) { throw 'Resolved Compose configuration is invalid.' }
if ($Pull) {
  & docker @compose pull
  if ($LASTEXITCODE -ne 0) { throw 'Image pull failed.' }
}
$up = @('up', '-d', '--wait', '--wait-timeout', "$WaitTimeoutSeconds")
if (-not $NoBuild) { $up += '--build' }
& docker @compose @up
if ($LASTEXITCODE -ne 0) { throw 'Polaris did not become ready.' }

$frontendPort = Get-DeploymentValue 'FRONTEND_PORT' '5173'
$gatewayPort = Get-DeploymentValue 'GATEWAY_PORT' '6080'
$enginePort = Get-DeploymentValue 'ENGINE_PORT' '6081'
Write-Host 'PASS: Polaris containers are ready.'
Write-Host "Dashboard:        http://localhost:$frontendPort"
Write-Host "Gateway readiness: http://localhost:$gatewayPort/readyz"
Write-Host "Engine readiness:  http://localhost:$enginePort/readyz"
Write-Host "Credentials:       $script:EnvironmentFile (do not commit or share)"
