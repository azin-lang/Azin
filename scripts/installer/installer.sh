#!/bin/sh
set -eu

REPO="https://github.com/azin-lang/Azin.git"
BIN="azc"
PREFIX="${PREFIX:-/usr/local}"
INSTALL_DIR="$PREFIX/bin"

require() {
    if ! command -v "$1" >/dev/null 2>&1; then
        printf 'Error: required command "%s" not found.\n' "$1" >&2
        exit 1
    fi
}

case "$(uname -s)" in
    Linux|Darwin|FreeBSD|OpenBSD|NetBSD|DragonFly|SunOS)
        ;;
    *)
        printf 'Error: unsupported operating system.\n' >&2
        exit 1
        ;;
esac

require git
require go

HAVE_CC=0
for cc in gcc clang cc; do
    if command -v "$cc" >/dev/null 2>&1; then
        HAVE_CC=1
        break
    fi
done

if [ "$HAVE_CC" -eq 0 ]; then
    printf 'Error: a C compiler (gcc, clang, or cc) is required.\n' >&2
    exit 1
fi

CLONE_DIR=$(mktemp -d)
trap 'rm -rf "$CLONE_DIR"' EXIT

printf 'Cloning %s...\n' "$REPO"
git clone --depth 1 "$REPO" "$CLONE_DIR"

cd "$CLONE_DIR"

sh scripts/build/build.sh

if [ ! -f "build/$BIN" ]; then
    printf 'Error: build failed, build/%s not found.\n' "$BIN" >&2
    exit 1
fi

find_privilege_command() {
    if [ "$(id -u)" -eq 0 ]; then
        return 0
    fi

    for cmd in doas sudo pkexec pfexec; do
        if command -v "$cmd" >/dev/null 2>&1; then
            PRIVCMD="$cmd"
            return 0
        fi
    done

    return 1
}

run_privileged() {
    if [ "$(id -u)" -eq 0 ]; then
        "$@"
        return
    fi

    if [ -n "${PRIVCMD:-}" ]; then
        "$PRIVCMD" "$@"
        return
    fi

    printf 'Error: Insufficient Permissions\n' >&2
    printf 'Install sudo, doas, pkexec or pfexec, or set PREFIX to a writable directory.\n' >&2
    exit 1
}

mkdir -p "$INSTALL_DIR" 2>/dev/null || true

if [ -w "$INSTALL_DIR" ]; then
    install -m755 "build/$BIN" "$INSTALL_DIR/$BIN"
else
    find_privilege_command || {
        printf 'Error: cannot write to %s.\n' "$INSTALL_DIR" >&2
        exit 1
    }
    run_privileged install -m755 "build/$BIN" "$INSTALL_DIR/$BIN"
fi

printf '\nAzin installed successfully.\n'
printf 'Executable: %s/%s\n' "$INSTALL_DIR" "$BIN"

if command -v "$INSTALL_DIR/$BIN" >/dev/null 2>&1; then
    "$INSTALL_DIR/$BIN" -version || true
else
    printf '\nNote: %s may not be in your PATH.\n' "$INSTALL_DIR"
    printf "Restart your shell or run 'hash -r'.\n"
fi