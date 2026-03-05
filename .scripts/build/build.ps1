<#
    .SYNOPSIS
    Build script for chtsht CLI.

    .DESCRIPTION
    Accepts a common set of parameters to automate (re)building the chtsht Go app.
    Injects version, commit, and build date via ldflags.

    .PARAMETER BinName
    The name of your executable, i.e. ./`$BinName.

    .PARAMETER BuildOS
    The OS to build for. Full list available at https://github.com/golang/go/blob/master/src/internal/syslist/syslist.go

    .PARAMETER BuildArch
    The CPU architecture to build for.

    .PARAMETER BuildOutputDir
    The build artifact path, where build outputs will be saved.

    .PARAMETER BuildTarget
    The name of the file to build (the entrypoint for your app).

    .PARAMETER ModulePath
    The Go module path used for ldflags injection. Check go.mod for the correct value.

    .EXAMPLE
    .\build.ps1
    .\build.ps1 -BinName "chtsht" -BuildOS "windows" -BuildArch "amd64" -BuildOutputDir "dist/"
    .\build.ps1 -BuildOS "linux"
#>
Param(
    [Parameter(Mandatory = $false, HelpMessage = "The name of your executable.")]
    $BinName = "chtsht",
    [Parameter(Mandatory = $false, HelpMessage = "The OS to build for.")]
    $BuildOS = "windows",
    [Parameter(Mandatory = $false, HelpMessage = "The CPU architecture to build for.")]
    $BuildArch = "amd64",
    [Parameter(Mandatory = $false, HelpMessage = "The build artifact path.")]
    $BuildOutputDir = "dist/",
    [Parameter(Mandatory = $false, HelpMessage = "The entrypoint for your app.")]
    $BuildTarget = "./app/cmd/chtsht",
    [Parameter(Mandatory = $false, HelpMessage = "The Go module path for ldflags.")]
    $ModulePath = "github.com/redjax/cheatsheets/internal/version"
)

## Get Git metadata
try {
    $GitVersion = (git describe --tags --always).Trim()
} catch {
    $GitVersion = "dev"
}

try {
    $GitCommit = (git rev-parse --short HEAD).Trim()
} catch {
    $GitCommit = "none"
}

$BuildDate = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")

## Set up ldflags to inject version metadata
$LdFlags = "-s -w " +
    "-X `"$ModulePath.Version=$GitVersion`" " +
    "-X `"$ModulePath.Commit=$GitCommit`" " +
    "-X `"$ModulePath.Date=$BuildDate`""

Write-Debug "BinName: $BinName"
Write-Debug "BuildOS: $BuildOS"
Write-Debug "BuildArch: $BuildArch"
Write-Debug "BuildOutputDir: $BuildOutputDir"
Write-Debug "BuildTarget: $BuildTarget"
Write-Debug "GitVersion: $GitVersion"
Write-Debug "GitCommit: $GitCommit"
Write-Debug "BuildDate: $BuildDate"

if ( $null -eq $BinName ) {
    Write-Warning "No bin name provided, pass the name of your executable using the -BinName flag"
    exit(1)
}

if ( ($BuildOS -eq "windows") -and ( -not $BinName.EndsWith(".exe") ) ) {
    Write-Warning "Building for Windows but bin name does not end with '.exe'. Appending .exe to '$BinName'"
    $BinName += ".exe"
}

$env:GOOS = $BuildOS
$env:GOARCH = $BuildArch

## Ensure output directory exists
if ( -not (Test-Path $BuildOutputDir) ) {
    New-Item -ItemType Directory -Path $BuildOutputDir -Force | Out-Null
}

$BuildOutput = Join-Path -Path $BuildOutputDir -ChildPath $BinName
Write-Debug "Build output: $BuildOutput"

Write-Host "Building chtsht" -ForegroundColor Cyan
Write-Host "  Version: $GitVersion  Commit: $GitCommit  Date: $BuildDate"
Write-Host "  Target:  $BuildTarget -> $BuildOutput"
Write-Host "  OS/Arch: $BuildOS/$BuildArch"

Write-Information "-- [ Build start"
try {
    go build -trimpath -ldflags "$LdFlags" -o $BuildOutput $BuildTarget
    Write-Host "Build successful" -ForegroundColor Green
}
catch {
    Write-Error "Error building app. Details: $($_.Exception.Message)"
    exit(1)
}
finally {
    Write-Information "-- [ Build complete"
}

exit(0)
