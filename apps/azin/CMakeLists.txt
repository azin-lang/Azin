add_executable(azin
        src/main.cpp
        src/new.cpp
        include/azin/version.hpp
        src/version.cpp
)

target_include_directories(azin PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_link_libraries(azin PRIVATE
        Azin::CompilerFlags
        Azin::Support::ANSI
        Azin::Support::Filesystem
)
