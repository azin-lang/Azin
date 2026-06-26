#pragma once
#include <cstdint>
#include <string_view>

namespace azin::frontend {

struct Loc {
    uint32_t line   = 1;
    uint32_t col    = 1;
    uint32_t offset = 0;
};

struct Span {
    Loc begin;
    Loc end;
};

struct SourceFile {
    std::string_view name;
    std::string_view text;
};

} // namespace azin::frontend