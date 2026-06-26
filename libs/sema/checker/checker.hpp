#pragma once
#include "ast/nodes/nodes.hpp"
#include "frontend/lexer/lexer.hpp"
#include "frontend/parser/parser.hpp"
#include <azin/support/fs/module_path.hpp>

#include <cstdint>
#include <filesystem>
#include <format>
#include <fstream>
#include <iterator>
#include <optional>
#include <stdexcept>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>

namespace azin::sema {

struct SemaError : std::runtime_error {
    frontend::Span span;
    SemaError(std::string const &msg, frontend::Span s)
        : std::runtime_error(msg)
        , span(s) {}
};

struct FnInfo {
    bool        zero_param  = false;
    std::string return_type;
    std::string module_name;
};

enum class NameResolution : uint8_t { Call, Reference };

struct ResolvedName {
    std::string    name;
    NameResolution resolution;
    frontend::Span span;
};

struct ExpectedType {
    enum class Kind : uint8_t { None, Named, FnRef } kind = Kind::None;
    std::string named;
};

class Checker {
public:
    auto declare_module(ast::Module const &mod) -> void {
        for (auto const &fn : mod.fns) {
            register_fn(fn, "");
        }
    }

    auto resolve_imports(ast::Module const &mod) -> void {
        for (auto const &imp : mod.imports) {
            load_module(imp.path, imp.span);
        }
    }

    auto ensure_std_namespace(std::vector<std::string> const &path, frontend::Span span) -> void {
        if (path.empty() || path.front() != "std") { return; }
        if (path.size() < 2) { return; }
        std::vector<std::string> const mod_path(path.begin(), path.begin() + static_cast<std::ptrdiff_t>(path.size() - 1));
        load_module(mod_path, span);
    }

    [[nodiscard]] auto fn_table() const -> std::unordered_map<std::string, FnInfo> const & {
        return fns_;
    }

    auto check_module(ast::Module const &mod) -> std::vector<ResolvedName> {
        std::vector<ResolvedName> resolved;
        for (auto const &fn : mod.fns) {
            check_fn(fn, resolved);
        }
        return resolved;
    }

    auto check_qualified_call(ast::QualifiedCallExpr const &qc) -> void {
        if (qc.path.size() < 2) { return; }

        std::vector<std::string> const mod_path(qc.path.begin(), qc.path.end() - 1);
        std::string const &fn_name = qc.path.back();
        std::string const  mod_key = join(mod_path, "::");

        ensure_std_namespace(qc.path, qc.span);

        auto it = module_fns_.find(mod_key);
        if (it == module_fns_.end()) {
            throw SemaError{
                std::format("module '{}' was not imported", mod_key),
                qc.span};
        }

        if (!it->second.contains(fn_name)) {
            throw SemaError{
                std::format("function '{}' not found in module '{}'", fn_name, mod_key),
                qc.span};
        }
    }

private:
    std::unordered_map<std::string, FnInfo>                          fns_;
    std::unordered_map<std::string, std::unordered_set<std::string>> module_fns_;
    std::unordered_set<std::string>                                  loaded_;

    auto register_fn(ast::FnDecl const &fn, std::string const &mod) -> void {
        std::string ret = "unit";
        if (fn.return_type != nullptr) {
            if (auto const *n = dynamic_cast<ast::NamedType const *>(fn.return_type.get())) {
                ret = n->name;
            }
        }
        fns_[fn.name] = FnInfo{.zero_param = fn.params.empty(), .return_type = ret, .module_name = mod};
    }

    static auto join(std::vector<std::string> const &parts, std::string_view sep) -> std::string {
        std::string out;
        for (std::size_t i = 0; i < parts.size(); ++i) {
            if (i > 0) { out += sep; }
            out += parts.at(i);
        }
        return out;
    }

    // NOLINTNEXTLINE(misc-no-recursion)
    auto load_module(std::vector<std::string> const &path, frontend::Span span) -> void {
        std::string const key = join(path, "::");
        if (loaded_.contains(key)) { return; }
        loaded_.insert(key);

        auto const file_path = support::fs::module_path(path);
        if (!std::filesystem::exists(file_path)) {
            throw SemaError{
                std::format("module '{}' not found (looked in '{}')",
                            key, file_path.string()),
                span};
        }

        std::ifstream f(file_path, std::ios::binary);
        if (!f) {
            throw SemaError{
                std::format("could not open module file '{}'", file_path.string()),
                span};
        }

        std::string const src{std::istreambuf_iterator<char>(f), std::istreambuf_iterator<char>()};
        std::string const sf_name = file_path.string();
        frontend::SourceFile const sf{.name = sf_name, .text = src};
        auto tokens = frontend::Lexer{sf}.tokenize();
        auto mod    = frontend::Parser{std::move(tokens)}.parse_module();

        std::unordered_set<std::string> exported;
        for (auto const &fn : mod.fns) {
            register_fn(fn, key);
            exported.insert(fn.name);
        }
        module_fns_[key] = std::move(exported);

        for (auto const &imp : mod.imports) {
            load_module(imp.path, imp.span);
        }
    }

    auto check_fn(ast::FnDecl const &fn, std::vector<ResolvedName> &out) -> void {
        for (auto const &stmt : fn.body) {
            check_stmt(*stmt, ExpectedType{}, out);
        }
    }

    // NOLINTNEXTLINE(misc-no-recursion)
    auto check_expr(ast::Expr const &expr, ExpectedType const &ctx, std::vector<ResolvedName> &out) -> void {
        if (auto const *ne = dynamic_cast<ast::NameExpr const *>(&expr)) {
            resolve_name(*ne, ctx, out);
        } else if (auto const *ce = dynamic_cast<ast::CallExpr const *>(&expr)) {
            for (auto const &arg : ce->args) { check_expr(*arg, ExpectedType{}, out); }
        } else if (auto const *qc = dynamic_cast<ast::QualifiedCallExpr const *>(&expr)) {
            check_qualified_call(*qc);
            for (auto const &arg : qc->args) { check_expr(*arg, ExpectedType{}, out); }
        }
    }

    auto check_stmt(ast::Stmt const &stmt, ExpectedType const & /*ctx*/, std::vector<ResolvedName> &out) -> void {
        if (auto const *ret = dynamic_cast<ast::ReturnStmt const *>(&stmt)) {
            if (ret->value != nullptr) { check_expr(*ret->value, ExpectedType{}, out); }
        } else if (auto const *es = dynamic_cast<ast::ExprStmt const *>(&stmt)) {
            check_expr(*es->expr, ExpectedType{}, out);
        } else if (auto const *vd = dynamic_cast<ast::VarDecl const *>(&stmt)) {
            check_var(*vd, out);
        }
    }

    auto check_var(ast::VarDecl const &vd, std::vector<ResolvedName> &out) -> void {
        if (vd.init == nullptr) { return; }

        ExpectedType ctx;
        if (vd.type_annotation != nullptr) {
            if (dynamic_cast<ast::FnType const *>(vd.type_annotation.get()) != nullptr) {
                ctx = ExpectedType{.kind = ExpectedType::Kind::FnRef, .named = {}};
            } else if (auto const *n = dynamic_cast<ast::NamedType const *>(vd.type_annotation.get())) {
                ctx = ExpectedType{.kind = ExpectedType::Kind::Named, .named = n->name};
            }
        }
        check_expr(*vd.init, ctx, out);
    }

    auto resolve_name(ast::NameExpr const &ne, ExpectedType const &ctx, std::vector<ResolvedName> &out) -> void {
        auto const it = fns_.find(ne.name);
        if (it == fns_.end()) { return; }

        FnInfo const &info = it->second;

        if (!info.zero_param) {
            out.push_back({.name = ne.name, .resolution = NameResolution::Reference, .span = ne.span});
            return;
        }

        switch (ctx.kind) {
        case ExpectedType::Kind::FnRef:
            out.push_back({.name = ne.name, .resolution = NameResolution::Reference, .span = ne.span});
            return;
        case ExpectedType::Kind::Named:
        case ExpectedType::Kind::None:
            out.push_back({.name = ne.name, .resolution = NameResolution::Call, .span = ne.span});
            return;
        }
    }
};

inline auto check_ambiguous_var(ast::VarDecl const &vd, std::unordered_map<std::string, FnInfo> const &fns) -> std::optional<SemaError> {
    if (vd.type_annotation != nullptr || vd.init == nullptr) { return std::nullopt; }

    auto const *ne = dynamic_cast<ast::NameExpr const *>(vd.init.get());
    if (ne == nullptr) { return std::nullopt; }

    auto const it = fns.find(ne->name);
    if (it == fns.end() || !it->second.zero_param) { return std::nullopt; }

    return SemaError{
        std::format("ambiguous use of zero-parameter function '{}': "
                    "add a type annotation (e.g. 'let {}: {} = {};' to call, "
                    "'let {}: fn = {};' to reference)",
                    ne->name,
                    vd.name, it->second.return_type, ne->name,
                    vd.name, ne->name),
        ne->span};
}

} // namespace azin::sema
