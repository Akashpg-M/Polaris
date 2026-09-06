param([switch]$RemoveVolumes)
. (Join-Path $PSScriptRoot 'deployment-common.ps1')
Initialize-DeploymentEnvironment
Assert-DockerAvailable
$compose = Get-ComposeArguments
$down = @('down')
if ($RemoveVolumes) {
  Write-Warning 'Removing volumes permanently deletes Polaris PostgreSQL, Redis, and Redpanda data.'
  $down += '--volumes'
}
& docker @compose @down
if ($LASTEXITCODE -ne 0) { throw 'Compose shutdown failed.' }
Write-Host $(if ($RemoveVolumes) { 'Polaris stopped and persistent volumes removed.' } else { 'Polaris stopped; persistent data was retained.' })
