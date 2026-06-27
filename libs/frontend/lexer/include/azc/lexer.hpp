#pragma once

#include <string_view>
#include <vector>

#include <azc/token.hpp>

namespace azc::frontend {

    class lexer {
    public:
        explicit lexer(std::string_view source);

        [[nodiscard]]
        auto tokenize() -> std::vector<token>;

    private:
        std::string_view m_source;

        std::size_t m_position{0};
        std::size_t m_line{1};
        std::size_t m_column{1};

        [[nodiscard]]
        auto is_at_end() const noexcept -> bool;

        [[nodiscard]]
        auto peek() const noexcept -> char;

        [[nodiscard]]
        auto peek_next() const noexcept -> char;

        auto advance() noexcept -> char;

        auto match(char expected) noexcept -> bool;

        auto skip_whitespace() noexcept -> void;

        auto scan_token(std::vector<token>& tokens) -> void;

        auto identifier(std::vector<token>& tokens) -> void;

        auto number(std::vector<token>& tokens) -> void;

        auto string(std::vector<token>& tokens) -> void;

        [[nodiscard]]
        auto make_token(
            token_kind kind,
            std::size_t start,
            std::size_t line,
            std::size_t column
        ) const -> token;

        [[nodiscard]]
        static auto identifier_kind(std::string_view lexeme) noexcept
            -> token_kind;
    };

} // namespace azc::frontend