Name:           azin
Version:        0.2.2
Release:        1%{?dist}
Summary:        A modern, high-performance systems programming language compiler

License:        MIT
URL:            https://azin-lang.org/

Source0:        https://github.com/azin-lang/Azin/archive/refs/tags/v%{version}.tar.gz

BuildRequires:  golang >= 1.22
Requires:       gcc
Recommends:     clang

%description
Azin is a modern, high-performance systems programming language compiler engineered for structural clarity, low-level execution control, and human readability.

%prep
%autosetup -n Azin-%{version}

%build
export CGO_ENABLED=0
go build -trimpath -ldflags="-s -w -X main.Version=%{version}" -o bin/azc ./cmd/azc

%install
install -D -m 0755 bin/azc %{buildroot}%{_bindir}/azc

%check
export CGO_ENABLED=0
export AZC=$(pwd)/bin/azc
go test ./...

%files
%license LICENSE
%doc README.md VERSIONING.md
%{_bindir}/azc

%changelog
* Fri Jul 31 2026 Azin Contributors <https://github.com/azin-lang/Azin> - 0.2.2-1
- Initial release for Copr
