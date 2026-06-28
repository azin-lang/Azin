#include <CLI/CLI.hpp>
#include <azc/cli.hpp>
#include <azc/source.hpp>
#include <cstdio> // NOLINT
#include <filesystem>
#include <azc/lexer.hpp>
#include <azc/token.hpp>
#include <azc/diagnostic.hpp>
#include <azc/diagnostic_engine.hpp>

// Disable the unreachable code warning for MSVC
#if defined(_MSC_VER) && !defined(__llvm__)
    #pragma warning(push)
    #pragma warning(disable : 4702)
#endif

#include <fmt/base.h>
#include <fmt/color.h>
#include <span>

#if defined(_MSC_VER) && !defined(__llvm__)
    #pragma warning(pop)
#endif

#include <string_view>

namespace cli = azc::cli;

namespace {

template <typename... Args>
void errorf(fmt::format_string<Args...> fmt_str, Args &&...args) {
    fmt::print(stderr, fg(fmt::color::red), fmt_str, std::forward<Args>(args)...);
    fmt::print(stderr, "\n");
}

void errprintln(std::string_view const msg) {
    fmt::print(stderr, fg(fmt::color::red), "error: {}\n", msg);
}

auto print_token_stream(std::span<azc::frontend::token const> tokens) -> void {
    for (const auto& token : tokens) {
        fmt::println(
            "{} \"{}\" ({}:{})",
            azc::frontend::token_kind_to_string(token.kind),
            token.lexeme,
            token.line,
            token.column
        );
    }
}

auto print_diagnostics(std::span<azc::frontend::diagnostic const> diagnostics) -> void {
    for (auto const& diagnostic : diagnostics) {
        errorf("{}", diagnostic.message);
    }
}

} // namespace

auto cli::run(int const argc, char const *const *argv) -> int {
    CLI::App app{"Azin Compiler"};

    bool version{false};
    bool print_tokens{false};
    std::filesystem::path input;

    app.add_flag("--version", version, "Display the compiler's version");
    app.add_option("input", input, "Source file to compile");
    app.add_flag("--print-tokens", print_tokens, "Print the tokens");

    CLI11_PARSE(app, argc, argv);

    if (version) {
        fmt::println("azc {}", AZIN_COMPILER_VERSION);
        return 0;
    }

    if (input.empty()) {
        errprintln("no input file specified\nUsage: azc <source>");
        return 1;
    }

    source::Manager source{input};

    if (!source.load()) {
        errorf("failed to open '{}'", input.string());
        return 1;
    }

    fmt::println("Loaded {} bytes", source.text().size());

    frontend::diagnostic_engine diagnostics;

    frontend::lexer lexer{source.text(), source.file_name(), diagnostics};
    auto tokens = lexer.tokenize();

    print_diagnostics(diagnostics.diagnostics());

    if (print_tokens) {
        print_token_stream(tokens);
    }

    return 0;
}
