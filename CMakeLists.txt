cmake_minimum_required(VERSION 3.22)
project(Azin VERSION 0.0.1 LANGUAGES CXX)

list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/cmake")
include(Language)
include(ProjectOptions)
include(DetectCompiler)
include(CompilerOptions)
include(Analyzers)

add_subdirectory(cli)
