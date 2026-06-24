#include <azin/colors.hpp>
#include <azin/version.hpp>

#include <iostream>
#include <span>
#include <string_view>

[[maybe_unused]] auto version_command(std::span<std::string_view const> /* unused */) -> int {
    std::cout << azin::ux::color::cyan << "Azin 0.0.1" << azin::ux::color::reset << "\n";
    return 0;
}
