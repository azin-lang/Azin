#pragma once
#include <algorithm>
#include <array>
#include <format>
#include <string>
#include <string_view>

#include "codes.hpp"

namespace azin::support::ansi {
template <std::size_t N>
struct [[nodiscard]] styled_view {
    std::string_view text;
    std::array<std::string_view, N> codes;

    constexpr styled_view(std::string_view const t,
                          std::array<std::string_view, N> const c) noexcept
        : text{t}
        , codes{c} {
    }

    [[nodiscard]] constexpr auto operator+(std::string_view const new_code) const noexcept
        -> styled_view<N + 1> {
        std::array<std::string_view, N + 1> extended_codes{};
        std::copy(codes.begin(), codes.end(), extended_codes.begin());
        extended_codes[N] = new_code;
        return styled_view<N + 1>{text, extended_codes};
    }

    [[nodiscard]] auto to_string() const -> std::string {
        std::size_t total_size = text.size() + code::reset.size();
        for (auto code : codes) {
            total_size += code.size();
        }

        std::string out;
        out.reserve(total_size);
        for (auto code : codes) {
            out.append(code);
        }
        out.append(text);
        out.append(code::reset);
        return out;
    }
};

template <std::size_t N>
styled_view(std::string_view, std::array<std::string_view, N>) -> styled_view<N>;

template <std::convertible_to<std::string_view>... Codes>
[[nodiscard]] auto compose(Codes... codes) -> std::string {
    std::string out;
    out.reserve((std::string_view{codes}.size() + ... + 0));
    (out.append(codes), ...);
    return out;
}

template <std::convertible_to<std::string_view>... Codes>
[[nodiscard]] auto wrap(std::string_view const text, Codes... codes) -> std::string {
    std::string out;
    out.reserve(text.size() + code::reset.size() + (std::string_view{codes}.size() + ... + 0));

    (out.append(codes), ...);
    out.append(text);
    out.append(code::reset);
    return out;
}

template <typename T>
concept StringViewable = std::convertible_to<T, std::string_view>;

[[nodiscard]] constexpr auto red(std::string_view const t) noexcept {
    return styled_view{t, std::array{code::red}};
}

template <std::size_t N>
[[nodiscard]] constexpr auto red(styled_view<N> const &v) noexcept {
    return v + code::red;
}

[[nodiscard]] constexpr auto green(std::string_view const t) noexcept {
    return styled_view{t, std::array{code::green}};
}

template <std::size_t N>
[[nodiscard]] constexpr auto green(styled_view<N> const &v) noexcept {
    return v + code::green;
}

[[nodiscard]] constexpr auto cyan(std::string_view const t) noexcept {
    return styled_view{t, std::array{code::cyan}};
}

template <std::size_t N>
[[nodiscard]] constexpr auto cyan(styled_view<N> const &v) noexcept {
    return v + code::cyan;
}

[[nodiscard]] constexpr auto bold(std::string_view const t) noexcept {
    return styled_view{t, std::array{code::bold}};
}

template <std::size_t N>
[[nodiscard]] constexpr auto bold(styled_view<N> const &v) noexcept {
    return v + code::bold;
}

template <std::size_t N>
auto operator<<(std::ostream &os, styled_view<N> const &view) -> std::ostream & {
    for (auto code : view.codes) {
        os << code;
    }
    return os << view.text << code::reset;
}
} // namespace azin::support::ansi

// NOLINTBEGIN(*-std-namespace-modification,cert-dcl58-cpp)
template <std::size_t N>
struct std::formatter<azin::support::ansi::styled_view<N>> : std::formatter<std::string_view> {
    auto format(azin::support::ansi::styled_view<N> const &view, auto &ctx) const {
        auto out = ctx.out();
        for (auto code : view.codes) {
            out = std::format_to(out, "{}", code);
        }

        // Synchronize the context's internal iterator state before rendering the text payload.
        // This ensures formatting options like width padding ({:<20}) function correctly.
        ctx.advance_to(out);
        out = std::formatter<std::string_view>::format(view.text, ctx);

        return std::format_to(out, "{}", azin::support::ansi::code::reset);
    }
};
// NOLINTEND(*-std-namespace-modification,cert-dcl58-cpp)