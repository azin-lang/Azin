/**
 * @file diagnostic.hpp
 * @brief Defines diagnostic types used by the Azin frontend.
 */

#pragma once

#include <cstdint>
#include <string>

namespace azc::frontend {

    /**
     * @brief Severity level of a diagnostic.
     */
    enum class diagnostic_severity : std::uint8_t {
        /// Informational message.
        note,

        /// Warning that does not prevent compilation.
        warning,

        /// Error that prevents successful compilation.
        error
    };

    /**
     * @brief Represents a compiler diagnostic.
     *
     * A diagnostic consists of a severity level and a human-readable
     * message describing the issue.
     */
    struct diagnostic {
        /// Severity of the diagnostic.
        diagnostic_severity severity;

        /// Human-readable diagnostic message.
        std::string message;
    };

} // namespace azc::frontend