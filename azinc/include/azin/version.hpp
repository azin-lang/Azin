#pragma once
#include <span>
#include <string_view>

auto version_command(std::span<std::string_view const> args) -> int;
