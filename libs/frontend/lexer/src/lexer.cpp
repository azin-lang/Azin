#include <cctype>
#include <azc/token.hpp>
#include <azc/lexer.hpp>

namespace azc::frontend {

lexer::lexer(std::string_view source)
    : m_source(source) {}

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

auto lexer::scan_token(std::vector<token>& tokens) -> void {
    const auto start = m_position;
    const auto line = m_line;
    const auto column = m_column;

    char c = advance();

    if (std::isalpha(static_cast<unsigned char>(c)) || c == '_') {
        --m_position;
        --m_column;
        identifier(tokens);
        return;
    }

    if (std::isdigit(static_cast<unsigned char>(c))) {
        --m_position;
        --m_column;
        number(tokens);
        return;
    }

    switch (c) {
        case '"':
            --m_position;
            --m_column;
            string(tokens);
            return;

        case '+':
            tokens.push_back(make_token(token_kind::plus, start, line, column));
            break;

        case '-':
            if (match('>'))
                tokens.push_back(make_token(token_kind::arrow, start, line, column));
            else
                tokens.push_back(make_token(token_kind::minus, start, line, column));
            break;

        case '*':
            tokens.push_back(make_token(token_kind::star, start, line, column));
            break;

        case '/':
            tokens.push_back(make_token(token_kind::slash, start, line, column));
            break;

        case '=':
            if (match('='))
                tokens.push_back(make_token(token_kind::equal_equal, start, line, column));
            else
                tokens.push_back(make_token(token_kind::equal, start, line, column));
            break;

        case '!':
            if (match('='))
                tokens.push_back(make_token(token_kind::bang_equal, start, line, column));
            else
                tokens.push_back(make_token(token_kind::bang, start, line, column));
            break;

        case '<':
            if (match('='))
                tokens.push_back(make_token(token_kind::less_equal, start, line, column));
            else
                tokens.push_back(make_token(token_kind::less, start, line, column));
            break;

        case '>':
            if (match('='))
                tokens.push_back(make_token(token_kind::greater_equal, start, line, column));
            else
                tokens.push_back(make_token(token_kind::greater, start, line, column));
            break;

        case '(':
            tokens.push_back(make_token(token_kind::left_paren, start, line, column));
            break;

        case ')':
            tokens.push_back(make_token(token_kind::right_paren, start, line, column));
            break;

        case '{':
            tokens.push_back(make_token(token_kind::left_brace, start, line, column));
            break;

        case '}':
            tokens.push_back(make_token(token_kind::right_brace, start, line, column));
            break;

        case ',':
            tokens.push_back(make_token(token_kind::comma, start, line, column));
            break;

        case ';':
            tokens.push_back(make_token(token_kind::semicolon, start, line, column));
            break;

        case ':':
            tokens.push_back(make_token(token_kind::colon, start, line, column));
            break;

        default:
            break;
    }
}

auto lexer::identifier(std::vector<token>& tokens) -> void {
    const auto start = m_position;
    const auto line = m_line;
    const auto column = m_column;

    while (std::isalnum(static_cast<unsigned char>(peek())) || peek() == '_')
        advance();

    std::string_view text = m_source.substr(start, m_position - start);

    tokens.push_back(token{
        .kind = identifier_kind(text),
        .lexeme = std::string(text),
        .offset = start,
        .line = line,
        .column = column
    });
}

auto lexer::number(std::vector<token>& tokens) -> void {
    const auto start = m_position;
    const auto line = m_line;
    const auto column = m_column;

    while (std::isdigit(static_cast<unsigned char>(peek())))
        advance();

    std::string_view text = m_source.substr(start, m_position - start);

    tokens.push_back(token{
        .kind = token_kind::integer_literal,
        .lexeme = std::string(text),
        .offset = start,
        .line = line,
        .column = column
    });
}

auto lexer::string(std::vector<token>& tokens) -> void {
    const auto start = m_position;
    const auto line = m_line;
    const auto column = m_column;

    advance();

    while (!is_at_end() && peek() != '"')
        advance();

    if (!is_at_end())
        advance();

    std::string_view text = m_source.substr(start, m_position - start);

    tokens.push_back(token{
        .kind = token_kind::string_literal,
        .lexeme = std::string(text),
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
        .lexeme = std::string(m_source.substr(start, m_position - start)),
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

    return token_kind::identifier;
}

} // namespace azc::frontend