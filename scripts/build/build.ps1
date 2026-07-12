# Usage: .\scripts\build\build.ps1

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$BuildDir = "build"
$Output = Join-Path $BuildDir "azc.exe"
$Source = "./cmd/azc"

New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null

Write-Host "Building Azin compiler..."


go build `
    -trimpath `
    -o "$Output" `
    "$Source"

Write-Host "Done: $Output"