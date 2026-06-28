#include "apm-install.hpp"
#include <cstdio>
#include <iostream>

void print_help() {
    std::cout << R"(╭────────────────────────────────────────────────────────────╮
│                    AZIN PACKAGE MANAGER                    │
│                                                            │
│  https://github.com/azin-lang/Azin                         │
│                                                            │
│  Usage:                                                    │
│    apm [options]                                           │
│                                                            │
│  Commands:                                                 │
│    install      Install a package                          │
│    upgrade      Upgrade installed packages                 │
│    remove       Remove a package                           │
│                                                            │
╰────────────────────────────────────────────────────────────╯
)";
}

int main (int argc, char *argv[]) {
    if (argc < 3) {
        print_help();    
        return 1;
    }
    
    Installer apm_install;
    
    std::string command = argv[1];
    std::string package_name = argv[2];

    if (command == "install") {
        std::string repo_url = "https://github.com/" + package_name;
        apm_install.apm_download(repo_url);
    }
    
    return 0;
}
