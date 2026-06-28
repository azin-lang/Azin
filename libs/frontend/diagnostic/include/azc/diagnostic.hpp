#pragma once

#include <string>
#include <cstdint>

namespace azc::frontend {

enum class diagnostic_severity: std::uint8_t {
    note,
    warning,
    error
};

struct diagnostic {
    diagnostic_severity severity;
    std::string message;
};

}
