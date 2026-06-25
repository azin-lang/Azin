#include <azin/support/ansi/styled_view.hpp>
#include <azin/support/fs/filesystem.hpp>

#include <expected>
#include <filesystem>
#include <fstream>
#include <iostream>
#include <iterator>
#include <string>
#include <utility>

namespace ansi = azin::support::ansi;
namespace fs = azin::support::fs;

using FileError = fs::FileError;
using Result = std::expected<void, FileError>;

namespace {
[[nodiscard]]
auto print_usage(int const argc, char const *const *argv) -> bool {
    if (argc >= 2) {
        return true;
    }

    std::cout << ansi::green(
        std::format("Usage: {} <source>\n", argc > 0 ? argv[0] : "azinc")); // NOLINT

    return false;
}

auto print_error(FileError const &err) -> void {
    std::cerr << ansi::red(err.message) << '\n';
}

[[nodiscard]]
auto ok(Result const &res) -> bool {
    if (!res) {
        print_error(res.error());
        return false;
    }
    return true;
}

[[nodiscard]]
auto read_file(std::filesystem::path const &path) -> std::expected<std::string, FileError> {
    auto file_result = fs::open_source_file(path);
    if (!file_result) {
        return std::unexpected(file_result.error());
    }

    std::ifstream file = std::move(*file_result);

    return std::string{std::istreambuf_iterator(file), std::istreambuf_iterator<char>()};
}

[[nodiscard]]
auto run(int const argc, char const *const *argv) -> int {
    if (!print_usage(argc, argv)) {
        return 1;
    }

    std::filesystem::path const path{argv[1]}; // NOLINT

    if (!ok(fs::check_file_exists(path))) {
        return 1;
    }

    if (!ok(fs::check_extension(path))) {
        return 1;
    }

    auto content = read_file(path);
    if (!content) {
        print_error(content.error());
        return 1;
    }

    std::cout << *content << '\n';
    return 0;
}

} // namespace

auto main(int const argc, char const *const *argv) -> int {
    try {
        return run(argc, argv);
    }
    catch (...) {
        return 1;
    }
}
