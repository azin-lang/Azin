/**
 * @file lexer.hpp
 * @brief Declares the Azin lexical analyzer.
 */

#pragma once

#include <string_view>
#include <vector>

#include <azin/token.hpp>
#include <azin/diagnostic_engine.hpp>

namespace azc::frontend {

    /**
     * @brief Converts Azin source code into a sequence of lexical tokens.
     *
     * The lexer performs lexical analysis by reading the source text,
     * recognizing keywords, identifiers, literals, operators, and
     * punctuation while reporting lexical errors through a diagnostic
     * engine.
     */
    class lexer {
    public:
        /**
         * @brief Constructs a lexer.
         *
         * @param source Source code to tokenize.
         * @param filename Name of the source file used in diagnostics.
         * @param diagnostics Diagnostic engine used to report lexer errors.
         */
        explicit lexer(
            std::string_view source,
            std::string_view filename,
            diagnostic_engine& diagnostics
        );

        /**
         * @brief Tokenizes the entire source file.
         *
         * Scans the input until the end of the source is reached and
         * appends an EOF token to the resulting token stream.
         *
         * @return A vector containing all generated tokens.
         */
        [[nodiscard]]
        auto tokenize() -> std::vector<token>;

    private:
        std::string_view m_source;
        std::string_view m_filename;
        diagnostic_engine& m_diagnostics; // NOLINT(cppcoreguidelines-avoid-const-or-ref-data-members)
        std::size_t m_position{0};
        std::size_t m_line{1};
        std::size_t m_column{1};

        /**
         * @brief Returns whether the lexer reached the end of the source.
         */
        [[nodiscard]]
        auto is_at_end() const noexcept -> bool;

        /**
         * @brief Returns the current character without consuming it.
         */
        [[nodiscard]]
        auto peek() const noexcept -> char;

        /**
         * @brief Returns the next character without consuming it.
         */
        [[nodiscard]]
        auto peek_next() const noexcept -> char;

        /**
         * @brief Consumes and returns the current character.
         *
         * Updates the current line and column counters.
         */
        auto advance() noexcept -> char;

        /**
         * @brief Consumes a character if it matches the expected value.
         *
         * @param expected Character to match.
         * @return true if the character matched.
         * @return false otherwise.
         */
        auto match(char expected) noexcept -> bool;

        /**
         * @brief Skips whitespace characters.
         */
        auto skip_whitespace() noexcept -> void;

        /**
         * @brief Scans a single token.
         *
         * @param tokens Token list receiving the scanned token.
         */
        auto scan_token(std::vector<token>& tokens) -> void;

        /**
         * @brief Lexes an identifier or keyword.
         *
         * @param tokens Token list receiving the scanned token.
         */
        auto identifier(std::vector<token>& tokens) -> void;

        /**
         * @brief Lexes an integer or floating-point literal.
         *
         * @param tokens Token list receiving the scanned token.
         */
        auto number(std::vector<token>& tokens) -> void;

        /**
         * @brief Lexes a character literal.
         *
         * @param tokens Token list receiving the scanned token.
         */
        auto character(std::vector<token>& tokens) -> void;

        /**
         * @brief Lexes a string literal.
         *
         * @param tokens Token list receiving the scanned token.
         */
        auto string(std::vector<token>& tokens) -> void;

        /**
         * @brief Advances until a delimiter or the end of the current line.
         *
         * Used for recovering after lexical errors.
         *
         * @param delimiter Character marking the end of recovery.
         */
        auto recover_to(char delimiter) noexcept -> void;

        /**
         * @brief Creates a token spanning the current lexeme.
         *
         * @param kind Token kind.
         * @param start Starting byte offset.
         * @param line Starting line.
         * @param column Starting column.
         *
         * @return Constructed token.
         */
        [[nodiscard]]
        auto make_token(
            token_kind kind,
            std::size_t start,
            std::size_t line,
            std::size_t column
        ) const -> token;

        /**
         * @brief Creates and appends a token.
         *
         * @param tokens Destination token list.
         * @param kind Token kind.
         * @param start Starting byte offset.
         * @param line Starting line.
         * @param column Starting column.
         */
        auto emit(
            std::vector<token>& tokens,
            token_kind kind,
            std::size_t start,
            std::size_t line,
            std::size_t column
        ) const -> void;

        /**
         * @brief Determines whether an identifier is a keyword.
         *
         * @param lexeme Identifier text.
         * @return Corresponding token kind.
         */
        [[nodiscard]]
        static auto identifier_kind(std::string_view lexeme) noexcept
            -> token_kind;
    };

} // namespace azc::frontend