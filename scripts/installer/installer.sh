#!/bin/sh
set -eu

REPO="https://github.com/azin-lang/Azin.git"
BIN="azc"

check_dep() {
    command -v "$1" >/dev/null 2>&1 || {
        echo "$1 is required."
        exit 1
    }
}

case "$(uname -s)" in
    Linux|Darwin|FreeBSD|OpenBSD|NetBSD)
        INSTALL_DIR="/usr/local/bin"
        ;;
    *)
        echo "Unsupported OS"
        exit 1
        ;;
esac

check_dep git
check_dep go
check_dep gcc

CLONE_DIR=$(mktemp -d)
trap 'rm -rf "$CLONE_DIR"' EXIT

git clone "$REPO" "$CLONE_DIR"

cd "$CLONE_DIR"

sh scripts/build/build.sh

if [ ! -f "build/$BIN" ]; then
    echo "Build failed: build/$BIN not found."
    exit 1
fi

sudo install -m 755 "build/$BIN" "$INSTALL_DIR/$BIN"

echo "Installed $BIN to $INSTALL_DIR/$BIN"
echo "Restart your shell or run 'hash -r' if your shell doesn't find it."