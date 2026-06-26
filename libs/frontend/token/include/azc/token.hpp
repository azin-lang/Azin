#pragma once

#include <cstddef>
#include <string>

namespace azc::frontend {

    enum class token_kind {
        // Identifiers & literals
        identifier,
        integer_literal,
        string_literal,

        // Keywords
        kw_fn,
        kw_var,
        kw_int,

        // Arithmetic
        plus,
        minus,
        star,
        slash,

        // Delimiters
        left_paren,
        right_paren,
        left_brace,
        right_brace,
        comma,
        semicolon,
        colon,

        // Operators
        equal,
        equal_equal,
        bang,
        bang_equal,
        less,
        less_equal,
        greater,
        greater_equal,
        arrow,

        eof,
    };

    struct token {
        token_kind kind;
        std::string lexeme;

        std::size_t offset;
        std::size_t line;
        std::size_t column;
    };

} // namespace azc::frontend