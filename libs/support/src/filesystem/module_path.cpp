#include <azin/support/fs/module_path.hpp>

#include <cstdlib>
#include <filesystem>
#include <stdexcept>
#include <string>
#include <vector>

namespace azin::support::fs {

auto azin_lib_root() -> std::filesystem::path {
#ifdef _WIN32
    // For some reason, getenv causes errors here.
    // This is why I hate MSVC
    #pragma warning(suppress: 4996)
    char const *home = std::getenv("USERPROFILE"); // NOLINT(concurrency-mt-unsafe)
#else
    char const *home = std::getenv("HOME");        // NOLINT(concurrency-mt-unsafe)
#endif
    if (home == nullptr || std::string_view{home}.empty()) {
        throw std::runtime_error("cannot determine home directory (HOME/USERPROFILE not set)");
    }
    return std::filesystem::path{home} / ".cache" / "azin" / "lib";
}

auto module_path(std::vector<std::string> const &parts) -> std::filesystem::path {
    // ["std", "fmt"] -> <root>/std/fmt/fmt.az
    auto root = azin_lib_root();
    for (auto const &p : parts) {
        root /= p;
    }
    root /= parts.back() + ".az";
    return root;
}

} // namespace azin::support::fs
