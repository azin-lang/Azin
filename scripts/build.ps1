# Usage: ./build.ps1

$AZC_SOURCE = ".\cmd\compiler"

go build -o azc.exe $AZC_SOURCE
