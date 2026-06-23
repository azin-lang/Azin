#include <azin/colors.hpp>
#include <azin/helpers.hpp>
#include <expected>
#include <filesystem>
#include <fstream>
#include <iostream>
#include <iterator>
#include <string>
#include <utility>

namespace fs = std::filesystem;

using FileError = azin::filesystem::FileError;
using Result = std::expected<void, FileError>;

namespace {
[[nodiscard]]
auto print_usage(int const argc, char const *const *argv) -> bool {
    if (argc >= 2) {
        return true;
    }

    std::cout << azin::ux::color::green << "Usage: " << (argc > 0 ? argv[0] : "azinc") // NOLINT
              << " <source>" << azin::ux::color::reset << '\n';

    return false;
}

auto print_error(FileError const &err) -> void {
    std::cerr << azin::ux::color::red << err.message << azin::ux::color::reset << '\n';
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
auto read_file(fs::path const &path) -> std::expected<std::string, FileError> {
    auto file_result = azin::filesystem::open_source_file(path);
    if (!file_result) {
        return std::unexpected(file_result.error());
    }

    std::ifstream file = std::move(*file_result);

    return std::string{std::istreambuf_iterator<char>(file), std::istreambuf_iterator<char>()};
}

[[nodiscard]]
auto run(int const argc, char const *const *argv) -> int {
    if (!print_usage(argc, argv)) {
        return 1;
    }

    fs::path const path{argv[1]}; // NOLINT

    if (!ok(azin::filesystem::check_file_exists(path))) {
        return 1;
    }

    if (!ok(azin::filesystem::check_extension(path))) {
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
