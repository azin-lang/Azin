#pragma once

#include <string_view>
#include <vector>

#include <azc/token.hpp>
#include <azc/diagnostic_engine.hpp>

namespace azc::frontend {

    class lexer {
    public:
        explicit lexer(std::string_view source, std::string_view filename, diagnostic_engine& diagnostics);

        [[nodiscard]]
        auto tokenize() -> std::vector<token>;

    private:
        std::string_view m_source;
        std::string_view m_filename;
        diagnostic_engine& m_diagnostics; // NOLINT(cppcoreguidelines-avoid-const-or-ref-data-members)

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

        auto character(std::vector<token>& tokens) -> void;

        auto string(std::vector<token>& tokens) -> void;

        auto recover_to(char delimiter) noexcept -> void;

        [[nodiscard]]
        auto make_token(
            token_kind kind,
            std::size_t start,
            std::size_t line,
            std::size_t column
        ) const -> token;

        auto emit(
            std::vector<token>& tokens,
            token_kind kind,
            std::size_t start,
            std::size_t line,
            std::size_t column
        ) const -> void;

        [[nodiscard]]
        static auto identifier_kind(std::string_view lexeme) noexcept
            -> token_kind;
    };

} // namespace azc::frontend