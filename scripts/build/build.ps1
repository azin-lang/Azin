<#
.SYNOPSIS
    Builds the Azin compiler (azc).
.DESCRIPTION
    Compiles ./cmd/azc into the build directory. Works from any working directory.
.PARAMETER OutputDir
    Target directory for the built binary (default: "build").
.PARAMETER Release
    Strips symbol tables and debug information (-s -w) for smaller release binaries.
.EXAMPLE
    .\scripts\build\build.ps1
.EXAMPLE
    .\scripts\build\build.ps1 -Release
#>
[CmdletBinding()]
param (
    [string]$OutputDir = "build",
    [switch]$Release
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$RepoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
Set-Location $RepoRoot

if (-not (Get-Command "go" -ErrorAction SilentlyContinue)) {
    Write-Error "Go executable not found in PATH. Please install Go before building."
    exit 1
}

# Compatibility with PowerShell < 7
$IsWin = if (Get-Variable -Name IsWindows -ErrorAction Ignore) {
    $IsWindows
} else {
    $env:OS -eq "Windows_NT"
}

$BinaryName = if ($IsWin) { "azc.exe" } else { "azc" }

$FullOutputDir = Join-Path $RepoRoot $OutputDir
if (-not (Test-Path $FullOutputDir)) {
    New-Item -ItemType Directory -Path $FullOutputDir -Force | Out-Null
}
$OutputPath = Join-Path $FullOutputDir $BinaryName
$SourcePath = "./cmd/azc"

$GoArgs = @("build", "-trimpath")
if ($Release) {
    $GoArgs += "-ldflags", "-s -w"
}
$GoArgs += "-o", $OutputPath, $SourcePath

Write-Host "Building Azin compiler (azc)..." -ForegroundColor Cyan
if ($Release) { Write-Host "Mode: Release (symbols stripped)" -ForegroundColor Yellow }

$Timer = [System.Diagnostics.Stopwatch]::StartNew()
& go @GoArgs
$Timer.Stop()

Write-Host ("Successfully built $OutputPath in {0:N2}s" -f $Timer.Elapsed.TotalSeconds) -ForegroundColor Green