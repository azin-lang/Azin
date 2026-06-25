#pragma once
#include <span>
#include <string_view>

auto new_command(std::span<std::string_view const> args) -> int;
