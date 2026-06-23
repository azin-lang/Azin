#include <azin/helpers.hpp>
#include <expected>
#include <filesystem>
#include <fstream>
#include <ios>
#include <string>

namespace azin::filesystem {

auto check_file_exists(fs::path const &source_path) -> Result {
    if (!fs::exists(source_path)) {
        return std::unexpected(FileError{"Source file does not exist: " + source_path.string()});
    }

    return {};
}

auto check_extension(fs::path const &source_path) -> Result {
    if (source_path.extension() != ".az") {
        return std::unexpected(
            FileError{"Invalid source file extension: " + source_path.extension().string()});
    }

    return {};
}

auto open_source_file(fs::path const &source_path) -> FileResult {
    std::ifstream file(source_path, std::ios::binary);

    if (!file) {
        return std::unexpected(FileError{"Failed to open source file: " + source_path.string()});
    }

    return file;
}

} // namespace azin::filesystem
