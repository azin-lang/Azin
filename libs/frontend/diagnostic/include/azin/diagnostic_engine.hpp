/**
 * @file diagnostic_engine.hpp
 * @brief Declares the diagnostic engine.
 */

#pragma once

#include <span>
#include <vector>

#include <azin/diagnostic.hpp>

namespace azc::frontend {

    /**
     * @brief Collects and manages compiler diagnostics.
     *
     * The diagnostic engine stores diagnostics reported by the frontend
     * and provides utilities for querying them after compilation.
     */
    class diagnostic_engine {
    public:
        /**
         * @brief Reports a diagnostic.
         *
         * The diagnostic is appended to the internal collection.
         *
         * @param diagnostic Diagnostic to report.
         */
        void report(diagnostic diagnostic);

        /**
         * @brief Returns whether any reported diagnostic is an error.
         *
         * @return true if at least one error diagnostic exists.
         * @return false otherwise.
         */
        [[nodiscard]]
        auto has_errors() const noexcept -> bool;

        /**
         * @brief Returns all reported diagnostics.
         *
         * @return Read-only view of the stored diagnostics.
         */
        [[nodiscard]]
        auto diagnostics() const noexcept -> std::span<const diagnostic>;

    private:
        /// Collection of reported diagnostics.
        std::vector<diagnostic> m_diagnostics;
    };

} // namespace azc::frontend