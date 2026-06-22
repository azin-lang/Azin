#include "colors.hpp"
#include <algorithm>
#include <iostream>
#include <sstream>
#include <cstdlib>
#include <vector>
#include <string>
#include <functional>
#include <memory>
#include <filesystem>

/// dumbfuck if you cant understand ts
void handle_error(const std::runtime_error& e) {
    std::cout << azin::ux::color::red << e.what() << azin::ux::color::reset << std::endl;
}

/// ok so this one might need an explanation but it's pretty self-explanatory for amateurs B)
/// for dumbasses: it checks the amount of arguments against the minimum required
void check_arguments(int argc, int min_args, const std::string& usage) {
    if (argc < min_args) {
        std::cout << azin::ux::color::green << usage << azin::ux::color::reset << "\n";
        exit(1);
    }
}

struct Command {
    std::string name;
    std::string description;

    //for beginners this just gives you a way to run a function with arguments
    std::function<int(int, char**)> execute;

    Command(std::string name, std::string desc, std::function<int(int, char**)> fn) :
        name(std::move(name)),
        description(std::move(desc)),
        execute(std::move(fn)) {}
};

/// registry for managing commands (WOOAH I DIDNT KNOW THAT!)
struct CommandRegistry {
    std::vector<Command> commands;

    /// i think you should go to a doctor if you dont understand what this does
    void addCommand(const Command& command) {
        commands.push_back(command);
    }

    /// same goes here lol
    // In CommandRegistry class
    void executeCommand(const std::string& name, const std::vector<std::string>& args) {
        for (const auto& cmd : commands) {
            if (cmd.name == name) {
                // Build argv array
                std::vector<char*> argv;
                argv.push_back(const_cast<char*>(name.c_str()));
                for (const auto& arg : args) {
                    argv.push_back(const_cast<char*>(arg.c_str()));
                }
                argv.push_back(nullptr);
                cmd.execute((int)argv.size() - 1, argv.data());
                return;
            }
        }
        std::cout << "Command not found: " << name << std::endl;
    }
};

CommandRegistry registry;

int TestCommand(int argc, char* argv[]) {
    std::cout << "Command!" << std::endl;
    return 0;
}

int buildCommand(int argc, char* argv[]) {
    std::cout << "Building... (STUB)" << std::endl;

    return 0;
}

/// version 0.0.1 (wow so groundbreaking)
int versionCommand(int argc, char* argv[]) {
    (void)argc; (void)argv;
    std::cout << azin::ux::color::cyan << "Azin 0.0.1" << azin::ux::color::reset << "\n";
    return 0;
}

/// helpp please!!!
int helpCommand(int argc, char* argv[]) {
    (void)argc; (void)argv;
    std::cout << azin::ux::color::green << "Usage: azin <command> <args>\n\n" << azin::ux::color::reset;
    std::cout << azin::ux::color::cyan << "Available commands:\n" << azin::ux::color::reset;

    if (registry.commands.empty()) {
        std::cout << "   " << azin::ux::color::red << "(no commands registered)" << azin::ux::color::reset << "\n";
        return 1;
    }

    // this comment was nuked by the turtle
    size_t max_len = 0;
    for (const auto& cmd : registry.commands)
        max_len = std::max(max_len, cmd.name.size());

    /// I CHOOSE DEATH!
    for (const Command& cmd : registry.commands) {
        std::string padding(max_len - cmd.name.size() + 2, ' ');
        std::string indent = "   " + std::string(max_len + 5, ' '); // For subsequent lines

        // Split description into lines
        std::vector<std::string> lines;
        std::stringstream ss(cmd.description);
        std::string line;
        while (std::getline(ss, line)) {
            lines.push_back(line);
        }

        // Print first line with name and padding
        std::cout << "   " << azin::ux::color::yellow << cmd.name << azin::ux::color::reset
                << padding << "- " << lines[0] << "\n";

        // Print remaining lines with indentation
        for (size_t i = 1; i < lines.size(); ++i) {
            std::cout << indent << lines[i] << "\n";
        }
    }

    std::cout << "\n";
    return 0;
};

int initialize() {
    registry.addCommand(Command("build", "Build the project", buildCommand));
    registry.addCommand(Command("test", "Testing command", TestCommand));
    registry.addCommand(Command("version", "Display version information", versionCommand));
    registry.addCommand(Command("help", "Display help information", helpCommand));

    return 0;
}

int main(int argc, char* argv[]) {
    int success = initialize();

    if (argc < 2) {
        helpCommand(argc, argv);
        return 1;
    }

    if (success != 0) {
        std::cerr << azin::ux::color::red << "Failed to initialize commands" << azin::ux::color::reset << "\n";
        return 1;
    }

    registry.executeCommand(argv[1], std::vector<std::string>(argv + 2, argv + argc));
}
