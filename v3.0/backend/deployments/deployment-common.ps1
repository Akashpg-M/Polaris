$ErrorActionPreference = 'Stop'
$script:DeploymentDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$script:ComposeFile = Join-Path $script:DeploymentDir 'docker-compose.yml'
$script:EnvironmentFile = Join-Path $script:DeploymentDir '.env'
$script:EnvironmentExample = Join-Path $script:DeploymentDir '.env.example'

function New-DeploymentHex([int]$bytes) {
  $buffer = New-Object byte[] $bytes
  $generator = [Security.Cryptography.RandomNumberGenerator]::Create()
  try { $generator.GetBytes($buffer) } finally { $generator.Dispose() }
  return ([BitConverter]::ToString($buffer) -replace '-', '').ToLowerInvariant()
}

function Initialize-DeploymentEnvironment {
  if (-not (Test-Path -LiteralPath $script:EnvironmentFile)) {
    $content = Get-Content -Raw -LiteralPath $script:EnvironmentExample
    $content = $content.Replace('replace_with_a_url_safe_random_secret', (New-DeploymentHex 32))
    $operator = "pol_op_$(New-DeploymentHex 8).$(New-DeploymentHex 32)"
    $content = $content.Replace('pol_op_replace_with_a_secure_random_secret', $operator)
    Set-Content -LiteralPath $script:EnvironmentFile -Value $content -Encoding utf8
    Write-Host "Created $script:EnvironmentFile with secure random local credentials."
  }
  $raw = Get-Content -Raw -LiteralPath $script:EnvironmentFile
  if ($raw -match 'replace_with_' -or $raw -notmatch '(?m)^POSTGRES_PASSWORD=\S+' -or $raw -notmatch '(?m)^DEV_PLATFORM_ADMIN_TOKEN=pol_op_[^.]+\.[^\s]+') {
    throw "Deployment .env is incomplete. Copy .env.example to .env and set POSTGRES_PASSWORD and DEV_PLATFORM_ADMIN_TOKEN."
  }
}

function Get-DeploymentValue([string]$name, [string]$fallback) {
  if (-not (Test-Path -LiteralPath $script:EnvironmentFile)) { return $fallback }
  $line = Get-Content -LiteralPath $script:EnvironmentFile | Where-Object { $_ -match "^$([regex]::Escape($name))=" } | Select-Object -Last 1
  if (-not $line) { return $fallback }
  $value = ($line -split '=', 2)[1].Trim()
  if ($value) { return $value }
  return $fallback
}

function Assert-DockerAvailable {
  docker version | Out-Null
  if ($LASTEXITCODE -ne 0) { throw 'Docker Desktop/Engine is not available.' }
  docker compose version | Out-Null
  if ($LASTEXITCODE -ne 0) { throw 'Docker Compose v2 is not available.' }
}

function Get-ComposeArguments {
  return @('compose', '--env-file', $script:EnvironmentFile, '-f', $script:ComposeFile)
}
