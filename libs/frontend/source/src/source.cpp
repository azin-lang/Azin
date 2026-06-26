#include <azc/source.hpp>

#include <cstddef>
#include <filesystem>
#include <ios>
#include <string_view>
#include <utility>

#include <fstream>
#include <iterator>


namespace source {

    Manager::Manager(std::filesystem::path path)
        : m_path(std::move(path)) {}

    auto Manager::load() -> bool {
        std::ifstream file(m_path, std::ios::binary);

        if (!file) {
            return false;
        }

        m_buffer.assign(
            std::istreambuf_iterator<char>{file},
            std::istreambuf_iterator<char>{}
        );

        m_position = 0;
        return true;
    }

    auto Manager::text() const noexcept -> std::string_view {
        return m_buffer;
    }

    auto Manager::current() const noexcept -> char {
        return peek();
    }

    auto Manager::peek(std::size_t offset) const noexcept -> char {
        const auto pos = m_position + offset;

        if (pos >= m_buffer.size()) {
            return '\0';
        }

        return m_buffer[pos]; // NOLINT(cppcoreguidelines-pro-bounds-avoid-unchecked-container-access)
    }

    auto Manager::advance() noexcept -> void {
        if (!eof()) {
            ++m_position;
        }
    }

    auto Manager::eof() const noexcept -> bool {
        return m_position >= m_buffer.size();
    }

    auto Manager::remaining() const -> std::string_view {
        return std::string_view{m_buffer}.substr(m_position);
    }

    auto Manager::position() const noexcept -> std::size_t {
        return m_position;
    }

    auto Manager::reset() noexcept -> void {
        m_position = 0;
    }

} // namespace source