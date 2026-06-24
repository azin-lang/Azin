get_property(AZIN_IS_MULTI_CONFIG GLOBAL PROPERTY GENERATOR_IS_MULTI_CONFIG)

if (NOT AZIN_IS_MULTI_CONFIG AND NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE "Debug" CACHE STRING "Choose the type of build." FORCE)
    set_property(CACHE CMAKE_BUILD_TYPE PROPERTY STRINGS
            "Debug"
            "Release"
            "RelWithDebInfo"
            "MinSizeRel"
    )
endif ()


# Generates compile_commands.json for language servers/LSPs
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Build all libraries as position independent code (fPIC) by default
set(CMAKE_POSITION_INDEPENDENT_CODE ON)

# Explicitly force the standard MultiThreaded DLL runtime for MSVC
if (MSVC)
    set(CMAKE_MSVC_RUNTIME_LIBRARY "MultiThreaded$<$<CONFIG:Debug>:Debug>DLL")
endif ()

set_property(GLOBAL PROPERTY USE_FOLDERS ON)
set_property(GLOBAL PROPERTY PREDEFINED_TARGETS_FOLDER "CMake")

include(CheckIPOSupported)
check_ipo_supported(RESULT _ipo_supported OUTPUT _ipo_error)

if (_ipo_supported)
    set(CMAKE_INTERPROCEDURAL_OPTIMIZATION_RELEASE ON)
    set(CMAKE_INTERPROCEDURAL_OPTIMIZATION_RELWITHDEBINFO ON)

    if (AZIN_IS_MULTI_CONFIG)
        message(STATUS "LTO configured for Release targets.")
    elseif (CMAKE_BUILD_TYPE STREQUAL "Release" OR CMAKE_BUILD_TYPE STREQUAL "RelWithDebInfo")
        message(STATUS "Link-Time Optimization (LTO) enabled.")
    endif ()
else ()
    message(STATUS "LTO not supported by compiler: ${_ipo_error}")
endif ()