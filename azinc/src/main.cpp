#include <azin/version.hpp>
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
    if (!fs::exists(source_path)) {
        std::cerr << azin::ux::color::red << "Source file does not exist: " << source_path << azin::ux::color::reset << "\n";
        return 1;
    }

    if (source_path.extension() != ".az") {
        std::cerr << azin::ux::color::red << "Invalid source file extension!\nFile is not an azin source file " << source_path.extension() << azin::ux::color::reset << "\n";
        return 1;
    }


    std::ifstream source_file(source_path);
    if (!source_file.is_open()) {
        std::cerr << azin::ux::color::red << "Failed to open source file: " << source_path << azin::ux::color::reset << "\n";
        return 1;
    }

    std::string source_code((std::istreambuf_iterator<char>(source_file)), std::istreambuf_iterator<char>());
    std::cout << source_code << std::endl;


    return 0;
}