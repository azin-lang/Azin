add_executable(azc
        src/main.cpp
)

target_include_directories(azc PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_link_libraries(azc PRIVATE
        Azin::CompilerFlags
        Azin::Support::ANSI
        Azin::Support::Filesystem
)
