#include <azin/helpers.hpp>
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


int azin::filesystem::checkFileExists(const fs::path& source_path) {
    if (!fs::exists(source_path)) {
        std::cerr << azin::ux::color::red << "Source file does not exist: " << source_path << azin::ux::color::reset << "\n";
        return 1;
    }
    return 0;
}

int azin::filesystem::checkExtension(const fs::path& source_path) {
    if (source_path.extension() != ".az") {
        std::cerr << azin::ux::color::red << "Invalid source file extension!\nFile is not an azin source file " << source_path.extension() << azin::ux::color::reset << "\n";
        return 1;
    }
    return 0;
}

std::ifstream azin::filesystem::openSourceFile(const fs::path& source_path) {
    std::ifstream source_file(source_path);
    if (!source_file.is_open()) {
        std::cerr << azin::ux::color::red << "Failed to open source file: " << source_path << azin::ux::color::reset << "\n";
        exit(1);
    }
    return source_file;
}
