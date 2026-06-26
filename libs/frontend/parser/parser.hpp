#pragma once
#include "ast/nodes/nodes.hpp"
#include "frontend/lexer/lexer.hpp"
#include "frontend/source/include/azc/source.hpp"

#include <format>
#include <memory>
#include <stdexcept>
#include <string>
#include <string_view>
#include <vector>

namespace azin::frontend {

struct ParseError : std::runtime_error {
    Span span;
    ParseError(std::string const &msg, Span s)
        : std::runtime_error(msg)
        , span(s) {}
};

class Parser {
public:
    explicit Parser(std::vector<Token> tokens)
        : tokens_(std::move(tokens)) {}

    // NOLINTNEXTLINE(misc-no-recursion)
    auto parse_module() -> ast::Module {
        ast::Module mod;
        while (!at(TK::Eof)) {
            if (at(TK::Import)) {
                mod.imports.push_back(parse_import());
            } else if (at(TK::Fn)) {
                mod.fns.push_back(parse_fn());
            } else {
                throw ParseError(
                    std::format("expected 'import' or 'fn', got '{}'", current().text),
                    current().span);
            }
        }
        return mod;
    }

private:
    std::vector<Token> tokens_;
    uint32_t           pos_ = 0;

    [[nodiscard]] auto current() const -> Token const & { return tokens_.at(pos_); }
    [[nodiscard]] auto at(TK k)  const -> bool          { return current().kind == k; }

    auto advance() -> Token const & {
        Token const &t = tokens_.at(pos_);
        if (t.kind != TK::Eof) { ++pos_; }
        return t;
    }

    [[nodiscard]] auto prev_span_end() const -> Loc { return tokens_.at(pos_ - 1).span.end; }

    auto expect(TK kind, std::string_view what) -> Token const & {
        if (!at(kind)) {
            throw ParseError(
                std::format("expected {}, got '{}'", what, current().text),
                current().span);
        }
        return advance();
    }

    auto parse_import() -> ast::ImportDecl {
        Loc const begin = current().span.begin;
        expect(TK::Import, "'import'");

        ast::ImportDecl decl;
        decl.path.push_back(expect(TK::Ident, "module name").text);
        while (at(TK::ColonColon)) {
            advance();
            decl.path.push_back(expect(TK::Ident, "module name").text);
        }

        if (at(TK::Semicolon)) { advance(); }
        decl.span = {.begin = begin, .end = prev_span_end()};
        return decl;
    }

    auto parse_fn() -> ast::FnDecl {
        Loc const begin = current().span.begin;
        expect(TK::Fn, "'fn'");

        ast::FnDecl decl;
        decl.name = expect(TK::Ident, "function name").text;

        if (at(TK::LParen)) {
            advance();
            while (!at(TK::RParen) && !at(TK::Eof)) {
                decl.params.push_back(parse_param());
                if (at(TK::Comma)) { advance(); }
            }
            expect(TK::RParen, "')'");
        }

        if (at(TK::Colon)) {
            advance();
            decl.return_type = parse_type();
        }

        while (!at(TK::End) && !at(TK::Eof)) {
            decl.body.push_back(parse_stmt());
        }
        expect(TK::End, "'end'");

        decl.span = {.begin = begin, .end = prev_span_end()};
        return decl;
    }

    auto parse_param() -> ast::Param {
        ast::Param p;
        p.span.begin = current().span.begin;
        p.name = expect(TK::Ident, "parameter name").text;
        expect(TK::Colon, "':'");
        p.type = parse_type();
        p.span.end = prev_span_end();
        return p;
    }

    // NOLINTNEXTLINE(misc-no-recursion)
    auto parse_type() -> std::unique_ptr<ast::TypeExpr> {
        Loc const begin = current().span.begin;

        if (at(TK::Fn)) {
            advance();
            auto fn_type = std::make_unique<ast::FnType>();

            if (at(TK::LParen)) {
                advance();
                while (!at(TK::RParen) && !at(TK::Eof)) {
                    fn_type->param_types.push_back(parse_type());
                    if (at(TK::Comma)) { advance(); }
                }
                expect(TK::RParen, "')'");
            }

            if (at(TK::Colon)) {
                advance();
                fn_type->return_type = parse_type();
            }

            fn_type->span = {.begin = begin, .end = prev_span_end()};
            return fn_type;
        }

        auto named = std::make_unique<ast::NamedType>();
        named->name = expect(TK::Ident, "type name").text;
        named->span = {.begin = begin, .end = prev_span_end()};
        return named;
    }

    auto parse_stmt() -> std::unique_ptr<ast::Stmt> {
        Loc const begin = current().span.begin;

        if (at(TK::Return)) { return parse_return(begin); }
        if (at(TK::Let))    { return parse_let(begin); }
        if (at(TK::Import)) {
            return std::make_unique<ast::ImportDecl>(parse_import());
        }

        auto expr = parse_expr();
        Span const span = {.begin = begin, .end = prev_span_end()};
        if (at(TK::Semicolon)) { advance(); }

        auto stmt  = std::make_unique<ast::ExprStmt>();
        stmt->span = span;
        stmt->expr = std::move(expr);
        return stmt;
    }

    auto parse_return(Loc begin) -> std::unique_ptr<ast::ReturnStmt> {
        advance();
        auto stmt = std::make_unique<ast::ReturnStmt>();

        if (!at(TK::Semicolon) && !at(TK::End)) {
            stmt->value = parse_expr();
        }

        if (at(TK::Semicolon)) { advance(); }
        stmt->span = {.begin = begin, .end = prev_span_end()};
        return stmt;
    }

    auto parse_let(Loc begin) -> std::unique_ptr<ast::VarDecl> {
        advance();
        auto decl  = std::make_unique<ast::VarDecl>();
        decl->name = expect(TK::Ident, "variable name").text;

        if (at(TK::Colon)) {
            advance();
            decl->type_annotation = parse_type();
        }

        if (at(TK::Equals)) {
            advance();
            decl->init = parse_expr();
        }

        if (at(TK::Semicolon)) { advance(); }
        decl->span = {.begin = begin, .end = prev_span_end()};
        return decl;
    }

    // NOLINTNEXTLINE(misc-no-recursion,readability-function-cognitive-complexity)
    auto parse_expr() -> std::unique_ptr<ast::Expr> {
        Loc const begin = current().span.begin;

        if (at(TK::IntLit)) {
            auto lit   = std::make_unique<ast::IntLitExpr>();
            lit->value = std::stoll(current().text);
            lit->span  = current().span;
            advance();
            return lit;
        }

        if (at(TK::StrLit)) {
            auto lit   = std::make_unique<ast::StrLitExpr>();
            lit->value = current().text;
            lit->span  = current().span;
            advance();
            return lit;
        }

        if (at(TK::Ident)) {
            std::string const first     = current().text;
            Span        const name_span = current().span;
            advance();

            if (at(TK::ColonColon)) {
                std::vector<std::string> path;
                path.push_back(first);
                while (at(TK::ColonColon)) {
                    advance();
                    path.push_back(expect(TK::Ident, "identifier after '::'").text);
                }
                expect(TK::LParen, "'(' after qualified path");
                auto call  = std::make_unique<ast::QualifiedCallExpr>();
                call->path = std::move(path);
                while (!at(TK::RParen) && !at(TK::Eof)) {
                    call->args.push_back(parse_expr());
                    if (at(TK::Comma)) { advance(); }
                }
                expect(TK::RParen, "')'");
                call->span = {.begin = begin, .end = prev_span_end()};
                return call;
            }

            if (at(TK::LParen)) {
                advance();
                auto call    = std::make_unique<ast::CallExpr>();
                call->callee = first;
                while (!at(TK::RParen) && !at(TK::Eof)) {
                    call->args.push_back(parse_expr());
                    if (at(TK::Comma)) { advance(); }
                }
                expect(TK::RParen, "')'");
                call->span = {.begin = begin, .end = prev_span_end()};
                return call;
            }

            auto name_expr  = std::make_unique<ast::NameExpr>();
            name_expr->name = first;
            name_expr->span = name_span;
            return name_expr;
        }

        throw ParseError(
            std::format("unexpected token '{}' in expression", current().text),
            current().span);
    }
};

} // namespace azin::frontend
