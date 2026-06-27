add_executable(azc
        src/main.cpp
        src/cli.cpp
)

target_include_directories(azc PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_compile_definitions(azc
        PRIVATE
        AZIN_COMPILER_VERSION="0.1.0"
)

target_link_libraries(azc PRIVATE
        Azin::CompilerFlags
        fmt
        Azin::Support::Filesystem

        Azin::Frontend::Source
        Azin::Frontend::Lexer

        CLI11::CLI11
)
azin_enable_analyzers(azc)
