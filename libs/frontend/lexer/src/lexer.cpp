#include <cctype>
#include <unordered_map>
#include <azin/token.hpp>
#include <azin/lexer.hpp>

#include <algorithm>
#include <cppcoro/generator.hpp>
#include <utility>

namespace azc::frontend {

lexer::lexer(std::string_view const source, std::string_view const filename,
             diagnostic_engine &diagnostics)
    : m_source(source)
    , m_filename(filename)
    , m_diagnostics(diagnostics) {
}

auto lexer::tokens() -> cppcoro::generator<token> {
    while (true) {
        skip_whitespace();

        if (is_at_end()) {
            break;
        }

        if (auto token = scan_token()) {
            co_yield *token;
        }
    }

    co_yield make_token(token_kind::eof, mark());
}

auto lexer::get_lexeme(token const &token) const noexcept -> std::string_view {
    return m_source.substr(token.offset, token.length);
}

auto lexer::mark() const noexcept -> token_start {
    return {
        .offset = m_position,
        .line = m_line,
        .column = m_column,
    };
}

auto lexer::make_token(token_kind kind, token_start start) const noexcept -> token {
    return {
        .kind = kind,
        .length = m_position - start.offset,
        .offset = start.offset,
        .line = start.line,
        .column = start.column,
    };
}

auto lexer::is_at_end() const noexcept -> bool {
    return m_position >= m_source.size();
}

auto lexer::peek() const noexcept -> char {
    return is_at_end() ? '\0' : m_source[m_position];
}

auto lexer::peek_next() const noexcept -> char {
    return (m_position + 1 < m_source.size()) ? m_source[m_position + 1] : '\0';
}

auto lexer::advance() noexcept -> char {
    char const c = m_source[m_position++];

    if (c == '\n') {
        ++m_line;
        m_column = 1;
    }
    else {
        ++m_column;
    }

    return c;
}

auto lexer::match(char expected) noexcept -> bool {
    if (is_at_end() || peek() != expected) {
        return false;
    }

    advance();
    return true;
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

auto lexer::recover_to(char delimiter) noexcept -> void {
    while (!is_at_end()) {
        if (peek() == delimiter) {
            advance();
            return;
        }

        if (peek() == '\n') {
            return;
        }

        advance();
    }
}

auto lexer::scan_token() -> std::optional<token> {
    auto start = mark();
    char const c = advance();

    if (std::isalpha(static_cast<unsigned char>(c)) != 0 || c == '_') {
        return scan_identifier(start);
    }

    if (std::isdigit(static_cast<unsigned char>(c)) != 0) {
        return scan_number(start);
    }

    switch (c) {
    case '\'':
        return scan_character(start);

    case '"':
        return scan_string(start);

    case '/':
        return scan_slash(start);

    default:
        return scan_operator(c, start);
    }
}

auto lexer::scan_identifier(token_start start) -> token {
    while (std::isalnum(static_cast<unsigned char>(peek())) != 0 || peek() == '_') {
        advance();
    }

    return make_token(identifier_kind(m_source.substr(start.offset, m_position - start.offset)),
                      start);
}

auto lexer::scan_number(token_start start) -> token {
    while (std::isdigit(static_cast<unsigned char>(peek())) != 0) {
        advance();
    }

    if (peek() == '.' && std::isdigit(static_cast<unsigned char>(peek_next())) != 0) {
        advance();

        while (std::isdigit(static_cast<unsigned char>(peek())) != 0) {
            advance();
        }

        return make_token(token_kind::float_literal, start);
    }

    return make_token(token_kind::integer_literal, start);
}

auto lexer::scan_character(token_start start) -> std::optional<token> {
    if (is_at_end()) {
        report_error(start.line, start.column, "Unterminated character literal.");
        return std::nullopt;
    }

    advance();

    if (!match('\'')) {
        report_error(start.line, start.column,
                     "Character literal must contain exactly one character.");
        recover_to('\'');
        return std::nullopt;
    }

    return make_token(token_kind::character_literal, start);
}

auto lexer::scan_string(token_start start) -> std::optional<token> {
    while (!is_at_end() && peek() != '"' && peek() != '\n') {
        advance();
    }

    if (is_at_end() || peek() == '\n') {
        report_error(start.line, start.column, "Unterminated string literal.");
        recover_to('"');
        return std::nullopt;
    }

    advance();

    return make_token(token_kind::string_literal, start);
}

auto lexer::scan_operator(char c, token_start start) -> std::optional<token> {
    switch (c) {
    case '+':
        if (match('=')) {
            return make_token(token_kind::plus_equal, start);
        }
        if (match('+')) {
            return make_token(token_kind::plus_plus, start);
        }
        return make_token(token_kind::plus, start);

    case '-':
        if (match('>')) {
            return make_token(token_kind::arrow, start);
        }
        if (match('=')) {
            return make_token(token_kind::minus_equal, start);
        }
        if (match('-')) {
            return make_token(token_kind::minus_minus, start);
        }
        return make_token(token_kind::minus, start);

    case '*':
        return make_token(match('=') ? token_kind::star_equal : token_kind::star, start);

    case '=':
        return make_token(match('=') ? token_kind::equal_equal : token_kind::equal, start);

    case '!':
        return make_token(match('=') ? token_kind::bang_equal : token_kind::bang, start);

    case '<':
        if (match('=')) {
            return make_token(token_kind::less_equal, start);
        }
        if (match('<')) {
            return make_token(token_kind::less_less, start);
        }
        return make_token(token_kind::less, start);

    case '>':
        if (match('=')) {
            return make_token(token_kind::greater_equal, start);
        }
        if (match('>')) {
            return make_token(token_kind::greater_greater, start);
        }
        return make_token(token_kind::greater, start);

    case '|':
        if (match('|')) {
            return make_token(token_kind::logical_or, start);
        }
        if (match('=')) {
            return make_token(token_kind::pipe_equal, start);
        }
        return make_token(token_kind::pipe, start);

    case '&':
        if (match('&')) {
            return make_token(token_kind::logical_and, start);
        }
        if (match('=')) {
            return make_token(token_kind::ampersand_equal, start);
        }
        return make_token(token_kind::ampersand, start);

    case '(':
        return make_token(token_kind::left_paren, start);

    case ')':
        return make_token(token_kind::right_paren, start);

    case '{':
        return make_token(token_kind::left_brace, start);

    case '}':
        return make_token(token_kind::right_brace, start);

    case '[':
        return make_token(token_kind::left_bracket, start);

    case ']':
        return make_token(token_kind::right_bracket, start);

    case ',':
        return make_token(token_kind::comma, start);

    case ';':
        return make_token(token_kind::semicolon, start);

    case ':':
        return make_token(token_kind::colon, start);

    case '.':
        return make_token(token_kind::dot, start);

    case '%':
        return make_token(match('=') ? token_kind::modulo_equal : token_kind::modulo, start);

    case '^':
        return make_token(match('=') ? token_kind::caret_equal : token_kind::caret, start);

    case '~':
        return make_token(token_kind::tilde, start);

    default:
        report_error(start.line, start.column, "Unexpected character '{}'.", c);
        return std::nullopt;
    }
}
  
auto lexer::scan_slash(token_start start) -> std::optional<token> {
    if (match('/')) {
        skip_line_comment();
        return std::nullopt;
    }

    if (match('*')) {
        skip_block_comment(start);
        return std::nullopt;
    }

    if (match('=')) {
        return make_token(token_kind::slash_equal, start);
    }

    return make_token(token_kind::slash, start);
}

auto lexer::skip_line_comment() noexcept -> void {
    while (!is_at_end() && peek() != '\n') {
        advance();
    }
}
  
auto lexer::skip_block_comment(token_start start) noexcept -> void {
    while (!is_at_end()) {
        if (peek() == '*' && peek_next() == '/') {
            advance();
            advance();
            return;
        }

        advance();
    }

    report_error(start.line, start.column, "Unterminated block comment.");
}

auto lexer::identifier_kind(std::string_view lexeme) noexcept -> token_kind {
    static constexpr std::pair<std::string_view, token_kind> keywords[] = {
        {"char", token_kind::kw_char},     {"end", token_kind::kw_end},
        {"float", token_kind::kw_float},   {"fn", token_kind::kw_fn},
        {"int", token_kind::kw_int},       {"return", token_kind::kw_return},
        {"string", token_kind::kw_string}, {"var", token_kind::kw_var},
    };

    auto const *it = std::lower_bound(
        std::begin(keywords), std::end(keywords), lexeme,
        [](auto const &keyword, std::string_view text) { return keyword.first < text; });

    if (it != std::end(keywords) && it->first == lexeme) {
        return it->second;
    }

    return token_kind::identifier;
}

} // namespace azc::frontend
