#pragma once

#include <vector>
#include <span>

#include <azc/diagnostic.hpp>

namespace azc::frontend {

class diagnostic_engine {
public:
    void report(diagnostic diagnostic);

    [[nodiscard]]
    auto has_errors() const noexcept -> bool;

    [[nodiscard]]
    auto diagnostics() const noexcept -> std::span<const diagnostic>;

private:
    std::vector<diagnostic> m_diagnostics;
};

}