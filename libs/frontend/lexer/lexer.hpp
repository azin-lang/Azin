#pragma once
#include "frontend/source/include/azc/source.hpp"

#include <cstdint>
#include <string>
#include <string_view>
#include <vector>

namespace azin::frontend {

enum class TK : uint8_t {
    Fn,     
    End,
    Return, 
    Let,       // Note for Alex & Stefan here: I couldn't remember what we were doing for variables at the time; change to var if needed
    Import,

    LParen,     // (
    RParen,     // )
    Comma,  
    Colon,
    ColonColon, // :: 
    Semicolon,
    Dot,
    Equals,

    Ident,
    IntLit,
    StrLit,

    Eof,
    Error,
};

struct Token {
    TK          kind;
    Span        span;
    std::string text;
};

class Lexer {
public:
    explicit Lexer(SourceFile src)
        : src_(src) {}

    auto tokenize() -> std::vector<Token> {
        std::vector<Token> out;
        while (true) {
            skip_whitespace_and_comments();
            if (at_end()) {
                out.push_back({.kind = TK::Eof, .span = {.begin = loc_, .end = loc_}, .text = {}});
                break;
            }
            out.push_back(next_token());
        }
        return out;
    }

private:
    SourceFile src_;
    uint32_t   pos_ = 0;
    Loc        loc_;

    [[nodiscard]] auto at_end() const -> bool { return pos_ >= src_.text.size(); }
    [[nodiscard]] auto peek()   const -> char { return src_.text.at(pos_); }
    [[nodiscard]] auto peek2()  const -> char {
        return pos_ + 1 < src_.text.size() ? src_.text.at(pos_ + 1) : '\0';
    }

    auto advance() -> char {
        char const c = src_.text.at(pos_++);
        if (c == '\n') { ++loc_.line; loc_.col = 1; }
        else           { ++loc_.col; }
        loc_.offset = pos_;
        return c;
    }

    auto skip_whitespace_and_comments() -> void {
        while (!at_end()) {
            char const c = peek();
            if (c == '/' && peek2() == '/') {
                while (!at_end() && peek() != '\n') { advance(); }
            } else if (c == ' ' || c == '\t' || c == '\r' || c == '\n') {
                advance();
            } else {
                break;
            }
        }
    }

    auto make_token(TK kind, Loc begin, std::string text = {}) -> Token {
        return Token{.kind = kind, .span = {.begin = begin, .end = loc_}, .text = std::move(text)};
    }

    auto next_token() -> Token {
        Loc  const begin = loc_;
        char const c     = advance();

        switch (c) {
        case '(': return make_token(TK::LParen,    begin);
        case ')': return make_token(TK::RParen,    begin);
        case ',': return make_token(TK::Comma,     begin);
        case ';': return make_token(TK::Semicolon, begin);
        case '.': return make_token(TK::Dot,       begin);
        case '=': return make_token(TK::Equals,    begin);
        case ':':
            if (peek() == ':') { advance(); return make_token(TK::ColonColon, begin); }
            return make_token(TK::Colon, begin);
        case '"': return lex_string(begin);
        default:
            if (is_digit(c)) { return lex_int(begin, c); }
            if (is_alpha(c)) { return lex_ident(begin, c); }
            return make_token(TK::Error, begin, std::string(1, c));
        }
    }

    auto lex_ident(Loc begin, char first) -> Token {
        std::string text(1, first);
        while (!at_end() && is_ident(peek())) { text += advance(); }

        TK kind = TK::Ident;
        if      (text == "fn")     { kind = TK::Fn; }
        else if (text == "end")    { kind = TK::End; }
        else if (text == "return") { kind = TK::Return; }
        else if (text == "let")    { kind = TK::Let; }
        else if (text == "import") { kind = TK::Import; }

        return make_token(kind, begin, kind == TK::Ident ? text : std::string{});
    }

    auto lex_int(Loc begin, char first) -> Token {
        std::string text(1, first);
        while (!at_end() && is_digit(peek())) { text += advance(); }
        return make_token(TK::IntLit, begin, std::move(text));
    }

    auto lex_string(Loc begin) -> Token {
        std::string text;
        while (!at_end() && peek() != '"') {
            char const c = advance();
            if (c == '\\' && !at_end()) {
                char const esc = advance();
                switch (esc) {
                case 'n':  text += '\n'; break;
                case 't':  text += '\t'; break;
                case '"':  text += '"';  break;
                case '\\': text += '\\'; break;
                default:   text += esc;  break;
                }
            } else {
                text += c;
            }
        }
        if (!at_end()) { advance(); }
        return make_token(TK::StrLit, begin, std::move(text));
    }

    [[nodiscard]] static auto is_digit(char c) -> bool { return c >= '0' && c <= '9'; }
    [[nodiscard]] static auto is_alpha(char c) -> bool {
        return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_';
    }
    [[nodiscard]] static auto is_ident(char c) -> bool { return is_alpha(c) || is_digit(c); }
};

} // namespace azin::frontend
