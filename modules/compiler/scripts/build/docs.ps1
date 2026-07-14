$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$OutputDir = "docs\compiler\api"

Write-Host "Generating API documentation..."

New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null

doc2go `
    -internal `
    -out $OutputDir `
    ./...


Write-Host "Documentation generated successfully."
Write-Host "Output: $OutputDir"
