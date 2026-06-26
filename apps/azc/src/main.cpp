#include <CLI/CLI.hpp>
#include <iostream>
#include <exception>

#include <azc/cli.hpp>

auto main(int const argc, char const *const *argv) -> int {
    try {
        return azc::cli::run(argc, argv);
    }
    catch (CLI::ParseError const& e) {
        return e.get_exit_code();
    }
    catch (std::exception const& e) {
        std::cerr << e.what() << '\n';
        return 1;
    }
    catch (...) {
        return 1;
    }
}
