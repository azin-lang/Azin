#include <algorithm>
#include <azin/colors.hpp>
#include <azin/new.hpp>
#include <azin/version.hpp>
#include <cstdio>
#include <exception>
#include <functional>
#include <print>
#include <span>
#include <string>
#include <string_view>
#include <utility>
#include <vector>

using Args = std::span<std::string_view const>;
using CommandFn = std::function<int(Args)>;

namespace {

struct Command {
    std::string name;
    std::string description;
    CommandFn execute;
};

class CommandRegistry {
public:
    auto add(Command command) -> void {
        commands_.push_back(std::move(command));
    }

    [[nodiscard]]
    auto execute(std::string_view const name, Args const args) const -> int {
        auto const cmd_it = std::ranges::find(commands_, name, &Command::name);

        if (cmd_it == commands_.end()) {
            std::println(stderr, "{}Unknown command: {}{}", azin::ux::color::red, name,
                         azin::ux::color::reset);
            return 1;
        }

        return cmd_it->execute(args);
    }

    [[nodiscard]]
    auto commands() const -> std::vector<Command> const & {
        return commands_;
    }

private:
    std::vector<Command> commands_;
};

auto help_command(CommandRegistry const &registry, Args /* unused */) -> int {
    auto const &commands = registry.commands();

    std::println("{}Usage: azin <command> [args...]{}\n", azin::ux::color::green,
                 azin::ux::color::reset);

    std::println("{}Available commands:{}", azin::ux::color::cyan, azin::ux::color::reset);

    std::size_t longest = 0;
    for (auto const &command : commands) {
        longest = std::max(longest, command.name.size());
    }

    for (auto const &command : commands) {
        std::string padded = command.name;
        padded.append(longest - padded.size(), ' ');
        std::println("  {}  {}", padded, command.description);
    }

    return 0;
}

auto test_command(Args /* unused */) -> int {
    std::println("Command!");
    return 0;
}

auto build_command(Args /* unused */) -> int {
    std::println("Building... (STUB)");
    return 0;
}

auto register_commands(CommandRegistry &registry) -> void {
    registry.add(Command{
        .name = "build",
        .description = "Build the project",
        .execute = build_command,
    });

    registry.add(Command{
        .name = "test",
        .description = "Testing command",
        .execute = test_command,
    });

    registry.add(Command{
        .name = "version",
        .description = "Display version information",
        .execute = version_command,
    });

    registry.add(Command{
        .name = "new",
        .description = "Create a new project <name>",
        .execute = new_command,
    });

    registry.add(Command{
        .name = "help",
        .description = "Display help information",
        .execute = [&registry](Args const args) -> int { return help_command(registry, args); },
    });
}

} // namespace

auto main(int const argc, char const *argv[]) -> int { // NOLINT(bugprone-exception-escape)
    try {
        CommandRegistry registry;

        register_commands(registry);

        std::span<char const *const> const argv_span{argv, static_cast<std::size_t>(argc)};

        if (argc < 2) {
            return help_command(registry, {});
        }

        std::vector<std::string_view> args;
        args.reserve(argv_span.size() > 2 ? argv_span.size() - 2 : 0);

        for (auto const *arg : argv_span.subspan(2)) {
            args.emplace_back(arg);
        }

        auto const *const command_name = argv_span.at(1);
        return registry.execute(command_name, args);
    }
    catch (std::exception const &exception) {
        std::println(stderr, "{}{}{}", azin::ux::color::red, exception.what(),
                     azin::ux::color::reset);
        return 1;
    }
    catch (...) {
        std::println(stderr, "{}An unknown error occurred.{}", azin::ux::color::red,
                     azin::ux::color::reset);
        return 1;
    }
}
