#include <azin/version.hpp>
#include <azin/helpers.hpp>
#include <azin/colors.hpp>
#include <filesystem>
#include <functional>
#include <algorithm>
#include <iostream>
#include <sstream>
#include <fstream>
#include <cstdlib>
#include <vector>
#include <string>
#include <memory>

namespace fs = std::filesystem;

/// dumbfuck if you cant understand ts
void handle_error(const std::runtime_error& error) {
    std::cout << azin::ux::color::red << error.what() << azin::ux::color::reset << std::endl;
}

/// ok so this one might need an explanation but it's pretty self-explanatory for amateurs B)
/// for dumbasses: it checks the amount of arguments against the minimum required
void check_arguments(int argc, int min_args, const std::string& usage) {
    if (argc < min_args) {
        std::cout << azin::ux::color::green << usage << azin::ux::color::reset << "\n";
        exit(1);
    }
}

int main(int argc, char* argv[]) {
    check_arguments(argc, 2, "Usage: azinc <source> <args>");

    fs::path source_path = fs::path(argv[1]);
    if (azin::filesystem::checkFileExists(source_path) != 0) {
        return 1;
    }

    if (azin::filesystem::checkExtension(source_path) != 0) {
        return 1;
    }

    std::ifstream source_file = azin::filesystem::openSourceFile(source_path);

    std::string source_code((std::istreambuf_iterator<char>(source_file)), std::istreambuf_iterator<char>());
    std::cout << source_code << std::endl;

    return 0;
}