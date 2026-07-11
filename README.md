# Azin

A modern, high-performance systems programming language engineered for structural clarity, low-level execution control, and human readability.

[![CI](https://github.com/azin-lang/azin/actions/workflows/build.yml/badge.svg)](https://github.com/azin-lang/azin/actions)
[![Website](https://img.shields.io/website?url=https%3A%2F%2Fazin-lang.github.io%2Fazin%2F&label=website&color=7289da)](https://azin-lang.org/)
[![Documentation](https://img.shields.io/website?url=https%3A%2F%2Fazin-lang.github.io%2Fazin%2F&label=documentation&color=7289da)](https://docs.azin-lang.org/)

---

## Key Design Philosophies

Azin bridges the gap between low-level performance-critical systems development and modern, highly expressive syntax language mechanics.

* **Explicit Block Scoping:** Replaces traditional brace nesting (`{}`) with clean, `do` / `end` scope structures.
* **Statically Typed:** Highly rigid compiler-enforced types to eliminate runtime type safety bugs.
* **Systems Architecture Native:** Completely free from heavy runtime garbage collection constraints, built for direct compilation targets.

---

## Language Showcase

Here is a quick glance at writing code in Azin:

### Structs and Pure Functions
```az
struct Point is
    x: int;
    y: int;
end

fn distance(a: Point, b: Point): int do
    return abs(a.x - b.x) + abs(a.y - b.y);
end
```

# Documentation

The official documentation for Azin can be found [here](https://docs.azin-lang.org/).
