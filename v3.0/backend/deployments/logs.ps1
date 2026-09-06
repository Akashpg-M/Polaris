param(
  [ValidateSet('all','frontend','gateway','engine','postgres','redis','redpanda','postgres-migrate','kafka-init')]
  [string]$Service = 'all',
  [int]$Tail = 200,
  [switch]$Follow
)
. (Join-Path $PSScriptRoot 'deployment-common.ps1')
Initialize-DeploymentEnvironment
Assert-DockerAvailable
$compose = Get-ComposeArguments
$arguments = @('logs', '--tail', "$Tail")
if ($Follow) { $arguments += '--follow' }
if ($Service -ne 'all') { $arguments += $Service }
& docker @compose @arguments
if ($LASTEXITCODE -ne 0) { throw 'Unable to read container logs.' }
