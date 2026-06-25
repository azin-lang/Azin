#include <azin/colors.hpp>
#include <azin/new.hpp>

#include <filesystem>
#include <format>
#include <fstream>
#include <iostream>
#include <span>
#include <string>
#include <string_view>
#include <system_error>

namespace fs = std::filesystem;

namespace {

auto write_file(fs::path const &path, std::string_view const content) -> bool {
    std::ofstream file(path, std::ios::binary);
    if (!file) {
        return false;
    }

    file.write(content.data(), static_cast<std::streamsize>(content.size()));
    return file.good();
}

auto create_dir(fs::path const &path) -> bool {
    std::error_code code;
    fs::create_directories(path, code);
    return !code;
}

auto create_project(std::string_view const name) -> int {
    fs::path const root = fs::current_path() / name;

    if (fs::exists(root)) {
        std::cout << azin::ux::color::red << "error: directory '" << name << "' already exists"
                  << azin::ux::color::reset << "\n";
        return 1;
    }

    if (!create_dir(root / "src") || !create_dir(root / "bin")) {
        std::cout << azin::ux::color::red << "error: failed to create project directories"
                  << azin::ux::color::reset << "\n";
        return 1;
    }

    std::string const toml = std::format(R"([project]
name = "{}"
version = "0.1.0"
azin-version = "0.0.1"
)",
                                         name);

    constexpr std::string_view main_src = R"(fn main() {
    // your azin code here
}
)";

    if (!write_file(root / "azin.toml", toml)) {
        std::cout << azin::ux::color::red << "error: failed to create azin.toml"
                  << azin::ux::color::reset << "\n";
        return 1;
    }

    if (!write_file(root / "src" / "main.az", main_src)) {
        std::cout << azin::ux::color::red << "error: failed to create src/main.az"
                  << azin::ux::color::reset << "\n";
        return 1;
    }

    std::cout << azin::ux::color::green << "Created new project '" << name << "'"
              << azin::ux::color::reset << "\n";

    return 0;
}

} // namespace

auto new_command(std::span<std::string_view const> const args) -> int {
    if (args.empty()) {
        std::cout << azin::ux::color::green << "Usage: azin new <project-name>"
                  << azin::ux::color::reset << "\n";
        return 1;
    }

    return create_project(args.front());
}
