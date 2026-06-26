#pragma once

#include <cstddef>
#include <filesystem>
#include <string>
#include <string_view>

namespace source {

    class Manager {
    public:
        explicit Manager(std::filesystem::path path);

        [[nodiscard]]
        auto load() -> bool;

        [[nodiscard]]
        auto text() const noexcept -> std::string_view;

        [[nodiscard]]
        auto current() const noexcept -> char;

        [[nodiscard]]
        auto peek(std::size_t offset = 0) const noexcept -> char;

        auto advance() noexcept -> void;

        [[nodiscard]]
        auto eof() const noexcept -> bool;

        [[nodiscard]]
        auto remaining() const -> std::string_view;

        [[nodiscard]]
        auto position() const noexcept -> std::size_t;

        auto reset() noexcept -> void;

    private:
        std::filesystem::path m_path;
        std::string m_buffer;
        std::size_t m_position{0};
    };

} // namespace source