#include <azin/support/fs/filesystem.hpp>

#include <filesystem>
#include <expected>
#include <fstream>
#include <ios>

namespace azin::support::fs {
auto check_file_exists(std::filesystem::path const &source_path) -> Result {
    if (!std::filesystem::exists(source_path)) {
        return std::unexpected(FileError{"Source file does not exist: " + source_path.string()});
    }

    return {};
}

auto check_extension(std::filesystem::path const &source_path) -> Result {
    if (source_path.extension() != ".az") {
        return std::unexpected(
            FileError{"Invalid source file extension: " + source_path.extension().string()});
    }

    return {};
}

auto open_source_file(std::filesystem::path const &source_path) -> FileResult {
    std::ifstream file(source_path, std::ios::binary);

    if (!file) {
        return std::unexpected(FileError{"Failed to open source file: " + source_path.string()});
    }

    return file;
}
} // namespace azin::support::fs
