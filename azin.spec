#
# spec file for package azin
#
# Copyright (c) 2026 Azin Contributors
#
# All modifications and additions to the file contributed by third parties
# remain the property of their respective owners, unless otherwise agreed
# upon. The license for this file, and modifications and additions to the
# file, is the same license as for the pristine package itself (MIT).
#

Name:           azin
Version:        0.2.2
Release:        1%{?dist}
Summary:        A modern, high-performance systems programming language compiler
License:        MIT
Group:          Development/Languages/Other
URL:            https://azin-lang.org
Source0:        %{name}-%{version}.tar.gz
# Source1 generated via OBS `obs-service-go_modules` or `go mod vendor`
Source1:        vendor.tar.gz

# Supported architectures for Go compiler (ARM 64-bit aarch64, ARM 32-bit arm, x86_64, etc.)
%if 0%{?fedora} || 0%{?rhel}
ExclusiveArch:  %{go_arches}
%else
ExclusiveArch:  aarch64 %{arm} x86_64 ppc64le s390x riscv64
%endif

BuildRequires:  golang >= 1.22

# Azin compiles source code by transpiling to C11 and invoking a native C compiler.
# GCC or Clang is required at runtime to generate native binaries.
Requires:       gcc
Recommends:     clang
Recommends:     glibc-devel

%description
Azin is a modern, high-performance systems programming language engineered for
structural clarity, low-level execution control, and human readability.

Key features:
- Explicit block scoping: Replaces traditional brace nesting ({}) with clean do / end blocks.
- Static typing: Compiler-enforced type safety with no runtime type checks or garbage collection.
- Minimalist punctuation: Eliminates unnecessary syntax where program structure is explicit.
- Systems-first: Designed for direct native compilation via C11 transpilation.

%prep
%autosetup -n %{name}-%{version} -a 1

%build
# Disable cgo for static compiler binary, set trimpath for reproducible builds
export CGO_ENABLED=0
export GOFLAGS="-buildmode=pie -trimpath -mod=vendor"

go build \
    -ldflags="-s -w -X main.Version=%{version}" \
    -o bin/azc \
    ./cmd/azc

%install
install -D -m 0755 bin/azc %{buildroot}%{_bindir}/azc

%check
# Run internal unit tests during build
export CGO_ENABLED=0
export GOFLAGS="-mod=vendor"
go test ./...

%files
%license LICENSE
%doc README.md VERSIONING.md
%{_bindir}/azc

%changelog
* Fri Jul 31 2026 Azin Packaging Team <packaging@azin-lang.org> - 0.2.2-1
- Initial RPM package spec for Open Build Service (OBS) with ARM (aarch64/armv7hl) support
