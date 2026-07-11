Write-Host "Generating API documentation..."

doc2go `
    -internal `
    -out docs/api `
    ./...

Write-Host "Done!"