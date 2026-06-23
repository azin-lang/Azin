option(AZIN_ENABLE_ANALYZERS "Enable static analysis tools" ON)

option(AZIN_ENABLE_CLANG_TIDY "Enable clang-tidy" ON)
option(AZIN_ENABLE_CPPCHECK "Enable cppcheck" OFF)
option(AZIN_ENABLE_IWYU "Enable include-what-you-use" OFF)

option(AZIN_ANALYZERS_STRICT "Treat analyzer warnings as errors" ON)

find_program(AZIN_CLANG_TIDY_EXE NAMES clang-tidy)
find_program(AZIN_CPPCHECK_EXE NAMES cppcheck)
find_program(AZIN_IWYU_EXE NAMES include-what-you-use)

set(AZIN_HAS_CLANG_TIDY OFF CACHE INTERNAL "")
set(AZIN_HAS_CPPCHECK OFF CACHE INTERNAL "")
set(AZIN_HAS_IWYU OFF CACHE INTERNAL "")

# MSVC and clang-tidy don't work well together from what I've tested; use clang-cl instead.
if (AZIN_ENABLE_ANALYZERS
        AND AZIN_ENABLE_CLANG_TIDY
        AND AZIN_CLANG_TIDY_EXE
        AND NOT MSVC)

    set(_clang_tidy_args
            ${AZIN_CLANG_TIDY_EXE}
            --extra-arg=-Wno-unknown-warning-option
            --extra-arg=-Wno-error=unknown-warning-option
    )

    if (CMAKE_CXX_COMPILER_ID STREQUAL "GNU")
        list(APPEND _clang_tidy_args
                --extra-arg-before=--driver-mode=g++)
    elseif (CMAKE_CXX_COMPILER_ID MATCHES "Clang")
        list(APPEND _clang_tidy_args
                --extra-arg-before=--driver-mode=clang++)
    endif ()

    if (AZIN_ANALYZERS_STRICT)
        list(APPEND _clang_tidy_args
                --warnings-as-errors=*)
    endif ()

    set(CMAKE_CXX_CLANG_TIDY
            "${_clang_tidy_args}"
            CACHE STRING "" FORCE)

    set(AZIN_HAS_CLANG_TIDY ON CACHE INTERNAL "" FORCE)

    message(STATUS "Analyzer: clang-tidy enabled")
endif ()

if (AZIN_ENABLE_ANALYZERS
        AND AZIN_ENABLE_CPPCHECK
        AND AZIN_CPPCHECK_EXE)

    set(_cppcheck_args
            ${AZIN_CPPCHECK_EXE}
            --enable=warning,performance,portability
            --inline-suppr
            --quiet
    )

    if (AZIN_ANALYZERS_STRICT)
        list(APPEND _cppcheck_args
                --error-exitcode=1)
    endif ()

    set(CMAKE_CXX_CPPCHECK
            "${_cppcheck_args}"
            CACHE STRING "" FORCE)

    set(AZIN_HAS_CPPCHECK ON CACHE INTERNAL "" FORCE)

    message(STATUS "Analyzer: cppcheck enabled")
endif ()

if (AZIN_ENABLE_ANALYZERS
        AND AZIN_ENABLE_IWYU
        AND AZIN_IWYU_EXE)

    set(CMAKE_CXX_INCLUDE_WHAT_YOU_USE
            ${AZIN_IWYU_EXE})

    set(AZIN_HAS_IWYU ON CACHE INTERNAL "" FORCE)

    message(STATUS "Analyzer: include-what-you-use enabled")
endif ()