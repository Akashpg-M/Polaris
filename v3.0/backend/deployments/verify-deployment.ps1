param(
  [switch]$SkipCodeChecks,
  [switch]$SkipImageBuild,
  [switch]$EndToEnd
)
. (Join-Path $PSScriptRoot 'deployment-common.ps1')
Initialize-DeploymentEnvironment
Assert-DockerAvailable
$compose = Get-ComposeArguments
$backendDir = Resolve-Path (Join-Path $PSScriptRoot '..')
$frontendDir = Resolve-Path (Join-Path $PSScriptRoot '..\..\frontend')

& docker @compose config --quiet
if ($LASTEXITCODE -ne 0) { throw 'Resolved Compose configuration is invalid.' }
Write-Host 'PASS: Compose configuration resolves with all required values.'

if (-not $SkipCodeChecks) {
  Push-Location $backendDir
  try {
    go test ./...
    if ($LASTEXITCODE -ne 0) { throw 'Go tests failed.' }
    go vet ./...
    if ($LASTEXITCODE -ne 0) { throw 'Go vet failed.' }
  } finally { Pop-Location }
  Push-Location $frontendDir
  try {
    npm ci
    if ($LASTEXITCODE -ne 0) { throw 'Frontend dependency installation failed.' }
    npm audit --omit=dev --audit-level=high
    if ($LASTEXITCODE -ne 0) { throw 'Frontend production dependency audit failed.' }
    npm run build
    if ($LASTEXITCODE -ne 0) { throw 'Frontend production build failed.' }
  } finally { Pop-Location }
  Write-Host 'PASS: backend and frontend code checks.'
}

if (-not $SkipImageBuild) {
  & docker @compose build
  if ($LASTEXITCODE -ne 0) { throw 'Container image build failed.' }
  Write-Host 'PASS: all locally built container images.'
}

if ($EndToEnd) {
  & (Join-Path $PSScriptRoot 'smoke-test.ps1')
  if ($LASTEXITCODE -ne 0) { throw 'End-to-end Compose smoke test failed.' }
}

Write-Host 'PASS: Polaris deployment package verification completed.'
