#include "azin/version.hpp"

#include <azin/support/ansi/styled_view.hpp>

#include <iostream>
#include <span>
#include <string_view>

namespace ansi = azin::support::ansi;

auto version_command(std::span<std::string_view const> /* unused */) -> int {
    std::cout << ansi::cyan("Azin 0.0.1\n");
    return 0;
}
