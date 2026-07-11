$AZC_SOURCE = ".\cmd\compiler"
$BUILD_DIR = ".\build"

New-Item -ItemType Directory -Path $BUILD_DIR -Force | Out-Null

go build -o "$BUILD_DIR\azc.exe" $AZC_SOURCE