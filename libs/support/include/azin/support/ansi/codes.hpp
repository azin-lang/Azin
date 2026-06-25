#pragma once
#include <string_view>
using namespace std::string_view_literals;

namespace azin::support::ansi::code {

inline constexpr auto reset = "\x1b[0m"sv;

// Styles
inline constexpr auto bold = "\x1b[1m"sv;
inline constexpr auto dim = "\x1b[2m"sv;
inline constexpr auto italic = "\x1b[3m"sv;
inline constexpr auto underline = "\x1b[4m"sv;
inline constexpr auto reverse = "\x1b[7m"sv;
inline constexpr auto hidden = "\x1b[8m"sv;

// Foreground colors
inline constexpr auto black = "\x1b[30m"sv;
inline constexpr auto red = "\x1b[31m"sv;
inline constexpr auto green = "\x1b[32m"sv;
inline constexpr auto yellow = "\x1b[33m"sv;
inline constexpr auto blue = "\x1b[34m"sv;
inline constexpr auto magenta = "\x1b[35m"sv;
inline constexpr auto cyan = "\x1b[36m"sv;
inline constexpr auto white = "\x1b[37m"sv;

// Bright foreground colors
inline constexpr auto bright_black = "\x1b[90m"sv;
inline constexpr auto bright_red = "\x1b[91m"sv;
inline constexpr auto bright_green = "\x1b[92m"sv;
inline constexpr auto bright_yellow = "\x1b[93m"sv;
inline constexpr auto bright_blue = "\x1b[94m"sv;
inline constexpr auto bright_magenta = "\x1b[95m"sv;
inline constexpr auto bright_cyan = "\x1b[96m"sv;
inline constexpr auto bright_white = "\x1b[97m"sv;

// Background colors
inline constexpr auto bg_black = "\x1b[40m"sv;
inline constexpr auto bg_red = "\x1b[41m"sv;
inline constexpr auto bg_green = "\x1b[42m"sv;
inline constexpr auto bg_yellow = "\x1b[43m"sv;
inline constexpr auto bg_blue = "\x1b[44m"sv;
inline constexpr auto bg_magenta = "\x1b[45m"sv;
inline constexpr auto bg_cyan = "\x1b[46m"sv;
inline constexpr auto bg_white = "\x1b[47m"sv;
} // namespace azin::support::ansi::code

namespace azin::support::ansi::cursor {
inline constexpr auto hide = "\x1b[?25l"sv;
inline constexpr auto show = "\x1b[?25h"sv;
inline constexpr auto up = "\x1b[1A"sv;
inline constexpr auto down = "\x1b[1B"sv;
inline constexpr auto left = "\x1b[1D"sv;
inline constexpr auto right = "\x1b[1C"sv;
} // namespace azin::support::ansi::cursor

namespace azin::support::ansi::screen {
inline constexpr auto clear = "\x1b[2J"sv;
inline constexpr auto clear_line = "\x1b[2K"sv;
inline constexpr auto home = "\x1b[H"sv;
} // namespace azin::support::ansi::screen
