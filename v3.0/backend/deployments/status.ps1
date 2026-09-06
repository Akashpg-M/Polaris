. (Join-Path $PSScriptRoot 'deployment-common.ps1')
Initialize-DeploymentEnvironment
Assert-DockerAvailable
$compose = Get-ComposeArguments
& docker @compose ps
if ($LASTEXITCODE -ne 0) { throw 'Unable to read Compose status.' }

$frontendPort = Get-DeploymentValue 'FRONTEND_PORT' '5173'
$gatewayPort = Get-DeploymentValue 'GATEWAY_PORT' '6080'
$enginePort = Get-DeploymentValue 'ENGINE_PORT' '6081'
foreach ($probe in @(
  @{ Name='frontend'; Uri="http://127.0.0.1:$frontendPort/healthz" },
  @{ Name='gateway'; Uri="http://127.0.0.1:$gatewayPort/readyz" },
  @{ Name='engine'; Uri="http://127.0.0.1:$enginePort/readyz" }
)) {
  try {
    $response = Invoke-WebRequest -UseBasicParsing -TimeoutSec 5 -Uri $probe.Uri
    Write-Host "$($probe.Name): HTTP $($response.StatusCode)"
  } catch {
    Write-Warning "$($probe.Name): unavailable ($($_.Exception.Message))"
  }
}
