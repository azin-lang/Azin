add_executable(azinc
    src/helpers.cpp
    src/version.cpp
    src/main.cpp
)

target_include_directories(azinc PRIVATE include)
target_link_libraries(azinc PRIVATE Azin::CompilerFlags)
