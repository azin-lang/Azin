# Usage: .\scripts\build\build.ps1

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$BuildDir = "build"
$Output = Join-Path $BuildDir "azin_lsp.exe"
$Source = ".\modules\lsp\cmd\lsp"

New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null

Write-Host "Building Azin compiler..."


go build `
    -trimpath `
    -o "$Output" `
    "$Source"

Write-Host "Done: $Output"