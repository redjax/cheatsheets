<#
    .SYNOPSIS
    Install script for chtsht (cheatsheets CLI) on Windows.

    .DESCRIPTION
    Downloads the latest release of chtsht from GitHub and installs it to
    %LOCALAPPDATA%\chtsht. Adds the install directory to the user's PATH
    if not already present.

    .PARAMETER Auto
    Install without prompts.

    .EXAMPLE
    # Run directly:
    .\.scripts\install.ps1

    # Run from the internet (one-liner):
    & ([scriptblock]::Create((irm https://raw.githubusercontent.com/redjax/cheatsheets/main/.scripts/install.ps1))) -Auto
#>

[CmdletBinding()]
Param(
    [switch] $Auto
)

$ErrorActionPreference = 'Stop'

$Repo = 'redjax/cheatsheets'
$BinName = 'chtsht'
$InstallPath = Join-Path $env:LOCALAPPDATA $BinName

## Create install directory if needed
if (-not (Test-Path $InstallPath)) {
    New-Item -ItemType Directory -Path $InstallPath -Force | Out-Null
}

## Check if already installed
$ExistingBin = Get-Command $BinName -ErrorAction SilentlyContinue
if ($ExistingBin) {
    if (-not $Auto) {
        $Confirm = Read-Host "$BinName is already installed at $($ExistingBin.Path). Download and install again? (y/N)"
        if ($Confirm -notin @('y', 'Y', 'yes', 'Yes', 'YES')) {
            Write-Host "Aborting."
            exit 0
        }
    } else {
        Write-Host "$BinName is already installed at $($ExistingBin.Path). Reinstalling..."
    }
}

## Fetch latest release from GitHub API
Write-Host "Fetching latest release info..."
try {
    $ReleaseApi = "https://api.github.com/repos/$Repo/releases/latest"
    $Release = Invoke-RestMethod -Uri $ReleaseApi -UseBasicParsing
} catch {
    Write-Error "Failed to fetch latest release info: $($_.Exception.Message)"
    throw
}

$ReleaseTag = $Release.tag_name
$Version = $ReleaseTag.TrimStart('v')
Write-Host "Latest version: $ReleaseTag"

## Detect CPU architecture
$ArchNorm = $null
try {
    $ArchCode = (Get-CimInstance Win32_Processor | Select-Object -First 1).Architecture
    switch ($ArchCode) {
        9  { $ArchNorm = 'amd64' }
        12 { $ArchNorm = 'arm64' }
        default {
            Write-Error "Unsupported architecture code: $ArchCode"
            throw "Unsupported architecture"
        }
    }
} catch {
    ## Fallback to environment variable
    $EnvArch = $env:PROCESSOR_ARCHITECTURE
    if ($EnvArch -match '^(AMD64|x86_64)$') {
        $ArchNorm = 'amd64'
    } elseif ($EnvArch -match '^ARM64$') {
        $ArchNorm = 'arm64'
    } else {
        Write-Error "Unsupported architecture: $EnvArch"
        throw "Unsupported architecture"
    }
}

if (-not $ArchNorm) {
    Write-Error "Failed to detect system architecture"
    throw "Failed to detect architecture"
}

## Build asset name: chtsht-windows-amd64-0.1.0.zip
$FileName = "$BinName-windows-$ArchNorm-$Version.zip"
$DownloadUrl = "https://github.com/$Repo/releases/download/$ReleaseTag/$FileName"

## Create temp directory
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("chtsht_install_" + [Guid]::NewGuid())
New-Item -ItemType Directory -Path $TempDir | Out-Null

$ZipPath = Join-Path $TempDir $FileName

## Download
Write-Host "Downloading $FileName..."
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $ZipPath -UseBasicParsing
} catch {
    Write-Error "Download failed: $($_.Exception.Message)"
    Write-Error "URL: $DownloadUrl"
    Remove-Item -Recurse -Force $TempDir -ErrorAction SilentlyContinue
    throw
}

## Extract
Write-Host "Extracting..."
try {
    Expand-Archive -Path $ZipPath -DestinationPath $TempDir -Force
} catch {
    Write-Error "Failed to extract archive: $($_.Exception.Message)"
    Remove-Item -Recurse -Force $TempDir -ErrorAction SilentlyContinue
    throw
}

## Find the binary
$ExtractedExe = Join-Path $TempDir "$BinName.exe"
if (-not (Test-Path $ExtractedExe)) {
    Write-Error "Binary '$BinName.exe' not found in archive."
    Get-ChildItem $TempDir | Format-Table Name
    Remove-Item -Recurse -Force $TempDir -ErrorAction SilentlyContinue
    throw "Binary not found"
}

## Copy to install path
$DestExePath = Join-Path $InstallPath "$BinName.exe"
try {
    Copy-Item -Path $ExtractedExe -Destination $DestExePath -Force
} catch {
    Write-Error "Failed to install ${BinName}.exe: $($_.Exception.Message)"
    Remove-Item -Recurse -Force $TempDir -ErrorAction SilentlyContinue
    throw
}

## Clean up temp directory
Remove-Item -Recurse -Force $TempDir -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "$BinName $ReleaseTag installed to $DestExePath" -ForegroundColor Green

## Add to PATH if not present
$UserPath = [Environment]::GetEnvironmentVariable('PATH', [EnvironmentVariableTarget]::User)
$PathEntries = $UserPath -split ';' | Where-Object { $_ -ne '' }

if ($InstallPath -notin $PathEntries) {
    try {
        [Environment]::SetEnvironmentVariable('PATH', "$UserPath;$InstallPath", 'User')
        Write-Host ""
        Write-Host "Added '$InstallPath' to user PATH." -ForegroundColor Cyan
        Write-Host "Close and reopen your terminal for PATH changes to take effect."
    } catch {
        Write-Warning @"

'$InstallPath' is not in your PATH. Add it manually by running:

  [Environment]::SetEnvironmentVariable('PATH', "`$env:PATH;$InstallPath", 'User')

Then close and reopen your terminal.
"@
    }
}

## Verify installation
Write-Host ""
try {
    & $DestExePath self version
} catch {
    Write-Host "Installed successfully, but could not run verification." -ForegroundColor Yellow
}
