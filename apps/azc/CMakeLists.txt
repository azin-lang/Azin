add_executable(azc
        src/main.cpp
        src/cli.cpp
)

target_include_directories(azc PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/include
        ${CMAKE_SOURCE_DIR}/libs
)

target_compile_definitions(azc
        PRIVATE
        AZIN_COMPILER_VERSION="0.1.0"
)

target_link_libraries(azc PRIVATE
        Azin::CompilerFlags
        Azin::Support::ANSI
        Azin::Support::Filesystem
        Azin::Frontend::Source
        CLI11::CLI11
)
