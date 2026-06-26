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
        kw_return,
        kw_char,
        kw_int,

        // Operators
        plus,
        minus,
        star,
        slash,
        equal,
        equal_equal,
        bang,
        bang_equal,
        less,
        less_equal,
        greater,
        greater_equal,
        arrow,

        // Delimiters
        left_paren,
        right_paren,
        left_brace,
        right_brace,
        comma,
        semicolon,
        colon,

        eof,
    };

    [[nodiscard]]
    constexpr std::string_view token_kind_to_string(token_kind kind) noexcept {
        switch (kind) {
            // Identifiers & literals
            case token_kind::identifier:       return "identifier";
            case token_kind::integer_literal:  return "integer_literal";
            case token_kind::string_literal:   return "string_literal";

            // Keywords
            case token_kind::kw_fn:            return "kw_fn";
            case token_kind::kw_var:           return "kw_var";
            case token_kind::kw_return:        return "kw_return";
            case token_kind::kw_char:          return "kw_char";
            case token_kind::kw_int:           return "kw_int";

            // Operators
            case token_kind::plus:             return "plus";
            case token_kind::minus:            return "minus";
            case token_kind::star:             return "star";
            case token_kind::slash:            return "slash";
            case token_kind::equal:            return "equal";
            case token_kind::equal_equal:      return "equal_equal";
            case token_kind::bang:             return "bang";
            case token_kind::bang_equal:       return "bang_equal";
            case token_kind::less:             return "less";
            case token_kind::less_equal:       return "less_equal";
            case token_kind::greater:          return "greater";
            case token_kind::greater_equal:    return "greater_equal";
            case token_kind::arrow:            return "arrow";

            // Delimiters
            case token_kind::left_paren:       return "left_paren";
            case token_kind::right_paren:      return "right_paren";
            case token_kind::left_brace:       return "left_brace";
            case token_kind::right_brace:      return "right_brace";
            case token_kind::comma:            return "comma";
            case token_kind::semicolon:        return "semicolon";
            case token_kind::colon:            return "colon";

            case token_kind::eof:              return "eof";
        }

        return "unknown";
    }


    struct token {
        token_kind kind;
        std::string lexeme;

        std::size_t offset;
        std::size_t line;
        std::size_t column;
    };

} // namespace azc::frontend