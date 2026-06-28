/**
 * @file source.hpp
 * @brief Declares the source file manager.
 */

#pragma once

#include <cstddef>
#include <filesystem>
#include <string>
#include <string_view>

namespace source {

/**
 * @brief Manages loading and sequential access to a source file.
 *
 * The Manager owns the contents of a source file and provides
 * utilities for reading characters, peeking ahead, tracking the
 * current position, and retrieving metadata about the source.
 */
class Manager {
public:
    /**
     * @brief Constructs a source manager.
     *
     * @param path Path to the source file.
     */
    explicit Manager(std::filesystem::path path);

    /**
     * @brief Loads the source file into memory.
     *
     * @return true if the file was loaded successfully.
     * @return false if the file could not be opened.
     */
    [[nodiscard]]
    auto load() -> bool;

    /**
     * @brief Returns the entire source text.
     *
     * @return View of the loaded source buffer.
     */
    [[nodiscard]]
    auto text() const noexcept -> std::string_view;

    /**
     * @brief Returns the current character.
     *
     * Equivalent to calling peek(0).
     *
     * @return Current character or '\0' at end of file.
     */
    [[nodiscard]]
    auto current() const noexcept -> char;

    /**
     * @brief Returns a character relative to the current position.
     *
     * @param offset Number of characters ahead to inspect.
     * @return Character at the requested position, or '\0' if the
     *         position is outside the source buffer.
     */
    [[nodiscard]]
    auto peek(std::size_t offset = 0) const noexcept -> char;

    /**
     * @brief Advances to the next character.
     *
     * Does nothing if the end of the source has been reached.
     */
    auto advance() noexcept -> void;

    /**
     * @brief Returns whether the end of the source has been reached.
     *
     * @return true if no more characters remain.
     */
    [[nodiscard]]
    auto eof() const noexcept -> bool;

    /**
     * @brief Returns the unread portion of the source.
     *
     * @return View beginning at the current position.
     */
    [[nodiscard]]
    auto remaining() const -> std::string_view;

    /**
     * @brief Returns the current character position.
     *
     * @return Zero-based offset into the source buffer.
     */
    [[nodiscard]]
    auto position() const noexcept -> std::size_t;

    /**
     * @brief Returns the path to the source file.
     *
     * @return Constant reference to the source path.
     */
    [[nodiscard]]
    auto path() const noexcept -> const std::filesystem::path&;

    /**
     * @brief Returns the source file name.
     *
     * @return File name without directory components.
     */
    [[nodiscard]]
    auto file_name() const noexcept -> std::string;

    /**
     * @brief Resets the current position to the beginning of the source.
     */
    auto reset() noexcept -> void;

private:
    /// Path to the source file.
    std::filesystem::path m_path;

    /// Contents of the loaded source file.
    std::string m_buffer;

    /// Current zero-based position within the source buffer.
    std::size_t m_position{0};
};

} // namespace source