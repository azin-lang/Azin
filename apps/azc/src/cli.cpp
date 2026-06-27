#include <azin/support/ansi/styled_view.hpp>

#include <CLI/CLI.hpp>
#include <azc/cli.hpp>
#include <azc/source.hpp>
#include <cstdio> // NOLINT
#include <filesystem>
#include <format>
#include <print>

namespace ansi = azin::support::ansi;
namespace cli = azc::cli;

auto cli::run(int const argc, char const *const *argv) -> int {
    CLI::App app{"Azin Compiler"};

    bool version{false};
    std::filesystem::path input;

    app.add_flag("--version", version, "Display the compiler's version");

    app.add_option("input", input, "Source file to compile");

    CLI11_PARSE(app, argc, argv);

    if (version) {
        std::println("azc {}", AZIN_COMPILER_VERSION);
        return 0;
    }

    if (input.empty()) {
        std::println(stderr, "{}", ansi::red("error: no input file specified"));
        std::println(stderr, "{}", ansi::red("Usage: azc <source>"));
        return 1;
    }

    source::Manager source{input};

    if (!source.load()) {
        std::println("{}", ansi::red(std::format("error: failed to open '{}'", input.string())));
        return 1;
    }

    // TODO: lexer parser etc etc
    std::println("Loaded {} bytes", source.text().size());

    // remove it later
    std::println("File Content:\n{}", source.text());

    return 0;
}
