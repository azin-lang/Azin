#include <cctype>
#include <azc/token.hpp>
#include <azc/lexer.hpp>
#include <fmt/format.h>
#include <string_view>
#include <azc/diagnostic.hpp>
#include <azc/diagnostic_engine.hpp>

namespace azc::frontend {

lexer::lexer(std::string_view source, std::string_view filename, diagnostic_engine& diagnostics)
    : m_source(source),
      m_filename(filename),
      m_diagnostics(diagnostics) {}

auto lexer::tokenize() -> std::vector<token> {
    std::vector<token> tokens;

    while (!is_at_end()) {
        skip_whitespace();

        if (is_at_end()) {
            break;
        }

        scan_token(tokens);
    }

    tokens.emplace_back(token{
        .kind = token_kind::eof,
        .lexeme = "",
        .offset = m_position,
        .line = m_line,
        .column = m_column
    });

    return tokens;
}

auto lexer::is_at_end() const noexcept -> bool {
    return m_position >= m_source.size();
}

auto lexer::peek() const noexcept -> char {
    return is_at_end() ? '\0' : m_source[m_position];
}

auto lexer::peek_next() const noexcept -> char {
    return (m_position + 1 >= m_source.size())
        ? '\0'
        : m_source[m_position + 1];
}

auto lexer::advance() noexcept -> char {
    char const c = m_source[m_position++];

    if (c == '\n') {
        ++m_line;
        m_column = 1;
    } else {
        ++m_column;
    }

    return c;
}

auto lexer::match(char expected) noexcept -> bool {
    if (is_at_end())
        return false;

    if (peek() != expected)
        return false;

    advance();
    return true;
}

auto lexer::recover_to(char delimiter) noexcept -> void {
    while (!is_at_end()) {
        if (peek() == delimiter) {
            advance(); // consume delimiter
            return;
        }

        if (peek() == '\n') {
            return;
        }

        advance();
    }
}


auto lexer::skip_whitespace() noexcept -> void {
    while (!is_at_end()) {
        switch (peek()) {
            case ' ':
            case '\t':
            case '\r':
            case '\n':
                advance();
                break;

            default:
                return;
        }
    }
}

auto lexer::emit(
        std::vector<token>& tokens,
        token_kind kind,
        std::size_t start,
        std::size_t line,
        std::size_t column
    ) const -> void {
    tokens.emplace_back(make_token(kind, start, line, column));
}


auto lexer::scan_token(std::vector<token>& tokens) -> void {
    const auto start = m_position;
    const auto line = m_line;
    const auto column = m_column;

    char c = peek();

    if (std::isalpha(static_cast<unsigned char>(c)) || c == '_') {
        identifier(tokens);
        return;
    }

    if (std::isdigit(static_cast<unsigned char>(c))) {
        number(tokens);
        return;
    }

    if (c == '\'') {
        character(tokens);
        return;
    }

    if (c == '"') {
        string(tokens);
        return;
    }

    c = advance();

    switch (c) {
        // it doesn't support escape characters like '\n' (yet)
        case '+':
            emit(tokens, token_kind::plus, start, line, column);
            break;

        case '-':
            if (match('>')) {
                emit(tokens, token_kind::arrow, start, line, column);
            } else {
                emit(tokens, token_kind::minus, start, line, column);
            }
            break;

        case '*':
            emit(tokens, token_kind::star, start, line, column);
            break;

        case '/':
            if (match('/')) {
                while (!is_at_end() && peek() != '\n') {
                    advance();
                }
            } else if (match('*')) {
                while (!is_at_end()) {
                    if (peek() == '*' && peek_next() == '/') {
                        advance(); // *
                        advance(); // /
                        break;
                    }

                    advance();
                }

                if (is_at_end()) {
                    m_diagnostics.report({
                        .severity = diagnostic_severity::error,
                        .message = fmt::format(
                            "{}:{}:{}: Unterminated block comment.",
                            m_filename,
                            line,
                            column
                        )
                    });

                    return;
                }
            } else {
                emit(tokens, token_kind::slash, start, line, column);
            }
            break;

        case '=':
            if (match('=')) {
                emit(tokens, token_kind::equal_equal, start, line, column);
            } else {
                emit(tokens, token_kind::equal, start, line, column);
            }
            break;

        case '!':
            if (match('=')) {
                emit(tokens, token_kind::bang_equal, start, line, column);
            } else {
                emit(tokens, token_kind::bang, start, line, column);
            }
            break;

        case '<':
            if (match('=')) {
                emit(tokens, token_kind::less_equal, start, line, column);
            } else {
                emit(tokens, token_kind::less, start, line, column);
            }
            break;

        case '>':
            if (match('=')) {
                emit(tokens, token_kind::greater_equal, start, line, column);
            } else {
                emit(tokens, token_kind::greater, start, line, column);
            }
            break;

        case '(':
            emit(tokens, token_kind::left_paren, start, line, column);
            break;

        case ')':
            emit(tokens, token_kind::right_paren, start, line, column);
            break;

        case '{':
            emit(tokens, token_kind::left_brace, start, line, column);
            break;

        case '}':
            emit(tokens, token_kind::right_brace, start, line, column);
            break;

        case ',':
            emit(tokens, token_kind::comma, start, line, column);
            break;

        case ';':
            emit(tokens, token_kind::semicolon, start, line, column);
            break;

        case ':':
            emit(tokens, token_kind::colon, start, line, column);
            break;

        default:
            // temporary, it should continue lexing
            // with a diagnostic
            // in the future, we should do this:
            // main.az:2:11: error: Unexpected character '@'.
            //
            // 2 |     var x = 5 @ 10;
            //   |               ^
            m_diagnostics.report({
                .severity = diagnostic_severity::error,
                .message = fmt::format(
                    "{}:{}:{}: Unexpected character '{}'.",
                    m_filename,
                    line,
                    column,
                    c
                )
            });
    }
}

auto lexer::identifier(std::vector<token>& tokens) -> void {
    const auto start = m_position;
    const auto line = m_line;
    const auto column = m_column;

    while (std::isalnum(static_cast<unsigned char>(peek())) || peek() == '_') {
        advance();
    }

    std::string_view const text =
        m_source.substr(start, m_position - start);

    tokens.push_back(token{
        .kind = identifier_kind(text),
        .lexeme = text,
        .offset = start,
        .line = line,
        .column = column
    });
}

auto lexer::number(std::vector<token>& tokens) -> void {
    const auto start = m_position;
    const auto line = m_line;
    const auto column = m_column;

    while (std::isdigit(static_cast<unsigned char>(peek()))) {
        advance();
    }

    token_kind kind = token_kind::integer_literal;

    if (peek() == '.' && std::isdigit(static_cast<unsigned char>(peek_next()))) {

        kind = token_kind::float_literal;

        advance(); // consume '.'

        while (std::isdigit(static_cast<unsigned char>(peek()))) {
            advance();
        }
    }

    std::string_view const text =
        m_source.substr(start, m_position - start);

    tokens.push_back(token{
        .kind = kind,
        .lexeme = text,
        .offset = start,
        .line = line,
        .column = column
    });
}

    auto lexer::character(std::vector<token>& tokens) -> void {
    const auto start = m_position;
    const auto line = m_line;
    const auto column = m_column;

    advance(); // opening '

    if (is_at_end()) {
        m_diagnostics.report({
            .severity = diagnostic_severity::error,
            .message = fmt::format(
                "{}:{}:{}: Unterminated character literal.",
                m_filename,
                line,
                column
            )
        });

        return;
    }

    advance(); // character

    if (peek() != '\'') {
        m_diagnostics.report({
            .severity = diagnostic_severity::error,
            .message = fmt::format(
                "{}:{}:{}: Character literal must contain exactly one character.",
                m_filename,
                line,
                column
            )
        });

        recover_to('\'');
        return;
    }

    advance(); // closing '

    std::string_view const text =
        m_source.substr(start, m_position - start);

    tokens.push_back(token{
        .kind = token_kind::character_literal,
        .lexeme = text,
        .offset = start,
        .line = line,
        .column = column
    });
}

auto lexer::string(std::vector<token>& tokens) -> void {
    const auto start = m_position;
    const auto line = m_line;
    const auto column = m_column;

    advance(); // opening "

    while (!is_at_end() && peek() != '"' && peek() != '\n') {
        advance();
    }

    if (is_at_end()) {
        m_diagnostics.report({
            .severity = diagnostic_severity::error,
            .message = fmt::format(
                "{}:{}:{}: Unterminated string literal.",
                m_filename,
                line,
                column
            )
        });

        recover_to('"');
        return;
    }

    advance(); // closing "

    std::string_view const text = m_source.substr(start, m_position - start);

    tokens.push_back(token{
        .kind = token_kind::string_literal,
        .lexeme = text,
        .offset = start,
        .line = line,
        .column = column
    });
}

auto lexer::make_token(
    token_kind kind,
    std::size_t start,
    std::size_t line,
    std::size_t column
) const -> token {
    return {
        .kind = kind,
        .lexeme = m_source.substr(start, m_position - start),
        .offset = start,
        .line = line,
        .column = column
    };
}

auto lexer::identifier_kind(std::string_view lexeme) noexcept
    -> token_kind {

    if (lexeme == "fn")     return token_kind::kw_fn;
    if (lexeme == "var")    return token_kind::kw_var;
    if (lexeme == "return") return token_kind::kw_return;
    if (lexeme == "end")    return token_kind::kw_end;
    if (lexeme == "char")   return token_kind::kw_char;
    if (lexeme == "int")    return token_kind::kw_int;
    if (lexeme == "float")  return token_kind::kw_float;
    if (lexeme == "string") return token_kind::kw_string;

    return token_kind::identifier;
}

} // namespace azc::frontend