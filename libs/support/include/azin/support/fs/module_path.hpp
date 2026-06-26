#pragma once
#include <filesystem>
#include <string>
#include <vector>

namespace azin::support::fs {

// Returns ~/.cache/azin/lib on linux/mac, %USERPROFILE%\.cache\azin\lib on windows
auto azin_lib_root() -> std::filesystem::path;

// Converts ["std", "fmt"] -> azin_lib_root() / "std" / "fmt.az"
auto module_path(std::vector<std::string> const &path_parts) -> std::filesystem::path;

} // namespace azin::support::fs
