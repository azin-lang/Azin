#pragma once
#include <expected>
#include <filesystem>
#include <fstream>

namespace azin::filesystem {
namespace fs = std::filesystem;

struct FileError {
    std::string message;
};

using Result = std::expected<void, FileError>;
using FileResult = std::expected<std::ifstream, FileError>;

auto check_file_exists(fs::path const &source_path) -> Result;
auto check_extension(fs::path const &source_path) -> Result;
auto open_source_file(fs::path const &source_path) -> FileResult;

} // namespace azin::filesystem
