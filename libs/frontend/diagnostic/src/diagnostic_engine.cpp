#include <azin/diagnostic_engine.hpp>

#include <span>

namespace azc::frontend {

void diagnostic_engine::report(diagnostic diagnostic) {
    m_diagnostics.push_back(std::move(diagnostic));
}

auto diagnostic_engine::has_errors() const noexcept -> bool {
    for (auto const& diagnostic : m_diagnostics) {
        if (diagnostic.severity ==
        diagnostic_severity::error) {
            return true;
        }
    }

    return false;
}

auto diagnostic_engine::diagnostics() const noexcept -> std::span<const diagnostic> {

    return m_diagnostics;
}

}