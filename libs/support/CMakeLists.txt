#-------------------------------------------------------------------------------

add_library(support_filesystem STATIC
        src/filesystem/filesystem.cpp
        include/azin/support/fs/filesystem.hpp
)

add_library(Azin::Support::Filesystem ALIAS support_filesystem)

target_include_directories(support_filesystem PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
)

target_link_libraries(support_filesystem PRIVATE Azin::CompilerFlags)
azin_enable_analyzers(support_filesystem)

#-------------------------------------------------------------------------------