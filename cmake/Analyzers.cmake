option(AZIN_ENABLE_ANALYZERS "Enable static analysis tools" ON)

option(AZIN_ENABLE_CLANG_TIDY "Enable clang-tidy" ON)
option(AZIN_ENABLE_CPPCHECK "Enable cppcheck" OFF)
option(AZIN_ENABLE_IWYU "Enable include-what-you-use" OFF)

option(AZIN_ANALYZERS_STRICT "Treat analyzer warnings as errors" ON)

find_program(AZIN_CLANG_TIDY_EXE NAMES clang-tidy)
find_program(AZIN_CPPCHECK_EXE NAMES cppcheck)
find_program(AZIN_IWYU_EXE NAMES include-what-you-use)

set(AZIN_HAS_CLANG_TIDY OFF)
set(AZIN_HAS_CPPCHECK OFF)
set(AZIN_HAS_IWYU OFF)

if (AZIN_ENABLE_ANALYZERS
        AND AZIN_ENABLE_CLANG_TIDY
        AND AZIN_CLANG_TIDY_EXE
        AND NOT MSVC)

    set(AZIN_CLANG_TIDY_COMMAND
            ${AZIN_CLANG_TIDY_EXE}
            --extra-arg=-Wno-unknown-warning-option
            --extra-arg=-Wno-error=unknown-warning-option
            --exclude-header-filter=^${CMAKE_BINARY_DIR}
    )

    if (CMAKE_CXX_COMPILER_ID STREQUAL "GNU")
        list(APPEND AZIN_CLANG_TIDY_COMMAND
                --extra-arg-before=--driver-mode=g++)
    elseif (CMAKE_CXX_COMPILER_ID MATCHES "Clang")
        list(APPEND AZIN_CLANG_TIDY_COMMAND
                --extra-arg-before=--driver-mode=clang++)
    endif ()

    if (AZIN_ANALYZERS_STRICT)
        list(APPEND AZIN_CLANG_TIDY_COMMAND
                --warnings-as-errors=*)
    endif ()

    set(AZIN_HAS_CLANG_TIDY ON)

    message(STATUS "Analyzer: clang-tidy enabled")
endif ()

if (AZIN_ENABLE_ANALYZERS
        AND AZIN_ENABLE_CPPCHECK
        AND AZIN_CPPCHECK_EXE)

    set(AZIN_CPPCHECK_COMMAND
            ${AZIN_CPPCHECK_EXE}
            --enable=warning,performance,portability
            --inline-suppr
            --quiet
    )

    if (AZIN_ANALYZERS_STRICT)
        list(APPEND AZIN_CPPCHECK_COMMAND
                --error-exitcode=1)
    endif ()

    set(AZIN_HAS_CPPCHECK ON)

    message(STATUS "Analyzer: cppcheck enabled")
endif ()

if (AZIN_ENABLE_ANALYZERS
        AND AZIN_ENABLE_IWYU
        AND AZIN_IWYU_EXE)

    set(AZIN_IWYU_COMMAND
            ${AZIN_IWYU_EXE})

    set(AZIN_HAS_IWYU ON)

    message(STATUS "Analyzer: include-what-you-use enabled")
endif ()

function(azin_enable_analyzers target)
    if (AZIN_HAS_CLANG_TIDY)
        set_property(TARGET ${target} PROPERTY
                CXX_CLANG_TIDY "${AZIN_CLANG_TIDY_COMMAND}")
    endif ()

    if (AZIN_HAS_CPPCHECK)
        set_property(TARGET ${target} PROPERTY
                CXX_CPPCHECK "${AZIN_CPPCHECK_COMMAND}")
    endif ()

    if (AZIN_HAS_IWYU)
        set_property(TARGET ${target} PROPERTY
                CXX_INCLUDE_WHAT_YOU_USE "${AZIN_IWYU_COMMAND}")
    endif ()
endfunction()