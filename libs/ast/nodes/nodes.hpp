#pragma once
#include "frontend/source/include/azc/source.hpp"

#include <cstdint>
#include <memory>
#include <string>
#include <vector>

namespace azin::ast {

using azin::frontend::Span;

struct TypeExpr {
    TypeExpr() = default;
    TypeExpr(TypeExpr const &) = delete;
    TypeExpr(TypeExpr &&) = default;
    auto operator=(TypeExpr const &) -> TypeExpr & = delete;
    auto operator=(TypeExpr &&) -> TypeExpr & = default;
    virtual ~TypeExpr() = default;
};

struct NamedType : TypeExpr {
    std::string name;
    Span        span;
};

struct FnType : TypeExpr {
    std::vector<std::unique_ptr<TypeExpr>> param_types;
    std::unique_ptr<TypeExpr>              return_type;
    Span                                   span;
};

struct Expr {
    Span span;
    Expr() = default;
    Expr(Expr const &) = delete;
    Expr(Expr &&) = default;
    auto operator=(Expr const &) -> Expr & = delete;
    auto operator=(Expr &&) -> Expr & = default;
    virtual ~Expr() = default;
};

struct NameExpr : Expr {
    std::string name;
};

struct CallExpr : Expr {
    std::string                        callee;
    std::vector<std::unique_ptr<Expr>> args;
};

struct QualifiedCallExpr : Expr {
    std::vector<std::string>           path;
    std::vector<std::unique_ptr<Expr>> args;
};

struct IntLitExpr : Expr {
    int64_t value = 0;
};

struct StrLitExpr : Expr {
    std::string value;
};

struct Stmt {
    Span span;
    Stmt() = default;
    Stmt(Stmt const &) = delete;
    Stmt(Stmt &&) = default;
    auto operator=(Stmt const &) -> Stmt & = delete;
    auto operator=(Stmt &&) -> Stmt & = default;
    virtual ~Stmt() = default;
};

struct ReturnStmt : Stmt {
    std::unique_ptr<Expr> value;
};

struct ExprStmt : Stmt {
    std::unique_ptr<Expr> expr;
};

struct VarDecl : Stmt {
    std::string               name;
    std::unique_ptr<TypeExpr> type_annotation;
    std::unique_ptr<Expr>     init;
};

struct ImportDecl : Stmt {
    std::vector<std::string> path;
    Span                     span;
};

struct Param {
    std::string               name;
    std::unique_ptr<TypeExpr> type;
    Span                      span;
};

struct FnDecl {
    std::string                        name;
    std::vector<Param>                 params;
    std::unique_ptr<TypeExpr>          return_type;
    std::vector<std::unique_ptr<Stmt>> body;
    Span                               span;
};

struct Module {
    std::vector<ImportDecl> imports;
    std::vector<FnDecl>     fns;
};

} // namespace azin::ast
