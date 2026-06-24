add_executable(azin
    src/main.cpp
    src/new.cpp
    src/version.cpp
)

target_include_directories(azin PRIVATE include)
target_link_libraries(azin PRIVATE Azin::CompilerFlags)
