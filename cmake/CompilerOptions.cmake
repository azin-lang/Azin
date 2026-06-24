add_library(azin_compiler_flags INTERFACE)
add_library(Azin::CompilerFlags ALIAS azin_compiler_flags)

option(AZIN_WARNINGS_AS_ERRORS "Treat compiler warnings as errors" ON)

if (AZIN_COMPILER_IS_MSVC OR AZIN_COMPILER_IS_CLANGCL)
    target_compile_options(azin_compiler_flags INTERFACE
            /W4
            /permissive-
            /Zc:__cplusplus
            /Zc:inline
            /volatile:iso
            /guard:cf

            $<$<CXX_COMPILER_ID:MSVC>:
            /external:anglebrackets
            /external:W0
            >

            $<$<BOOL:${AZIN_WARNINGS_AS_ERRORS}>:/WX>
    )
    if (AZIN_COMPILER_IS_MSVC)
        target_compile_options(azin_compiler_flags INTERFACE
                /w14242
                /w14254
                /w14263
                /w14265
                /w14287
                /we4289
                /w14296
                /w14311
                /w14545
                /w14546
                /w14547
                /w14549
                /w14555
                /w14619
                /w14640
                /w14826
                /w14905
                /w14906
                /w14928
                /w15038
        )
    elseif (AZIN_COMPILER_IS_CLANGCL)
        target_compile_options(azin_compiler_flags INTERFACE
                -fcolor-diagnostics
                -Winconsistent-missing-override
        )
    endif ()
elseif (AZIN_COMPILER_IS_GCC OR AZIN_COMPILER_IS_CLANG OR AZIN_COMPILER_IS_APPLECLANG)

    target_compile_options(azin_compiler_flags INTERFACE
            -Wall
            -Wextra
            -Wpedantic

            -Wshadow
            -Wnon-virtual-dtor
            -Wold-style-cast
            -Wcast-align
            -Wcast-qual

            -Wunused
            -Woverloaded-virtual
            -Wconversion
            -Wsign-conversion

            -Wnull-dereference
            -Wdouble-promotion

            -Wformat=2
            -Wmisleading-indentation
            -Wimplicit-fallthrough

            -Wswitch-enum
            -Wredundant-move
            -Wundef

            -Wextra-semi
            -Wreturn-type
            -Wuninitialized
            -Wstrict-overflow
            -Wrange-loop-construct
            -Wmissing-declarations
            -Wtautological-compare

            -fstack-protector-strong

            $<$<BOOL:${AZIN_WARNINGS_AS_ERRORS}>:-Werror>
    )

    if (AZIN_COMPILER_IS_GCC)

        target_compile_options(azin_compiler_flags INTERFACE
                -Wduplicated-cond
                -Wduplicated-branches
                -Wlogical-op
                -Wuseless-cast
                -Wsuggest-override
                -Walloc-zero
                -fdiagnostics-color=always
                -Wtrampolines
                -Wmaybe-uninitialized
        )

    elseif (AZIN_COMPILER_IS_CLANG OR AZIN_COMPILER_IS_APPLECLANG)

        target_compile_options(azin_compiler_flags INTERFACE
                -fcolor-diagnostics
                -Wdocumentation
                -Winconsistent-missing-override
                -Wweak-vtables
                -Wuninitialized
        )

    endif ()

endif ()

target_compile_definitions(azin_compiler_flags INTERFACE

        $<$<CONFIG:Debug,RelWithDebInfo>:_GLIBCXX_ASSERTIONS>

        $<$<CONFIG:Debug>:
        _LIBCPP_HARDENING_MODE=_LIBCPP_HARDENING_MODE_DEBUG
        >

        $<$<CONFIG:RelWithDebInfo>:
        _LIBCPP_HARDENING_MODE=_LIBCPP_HARDENING_MODE_EXTENSIVE
        >

        $<$<CONFIG:Release,MinSizeRel>:
        _LIBCPP_HARDENING_MODE=_LIBCPP_HARDENING_MODE_FAST
        >
)

if (CMAKE_CXX_COMPILER_ID STREQUAL "Clang" AND UNIX AND NOT APPLE)
    target_compile_options(azin_compiler_flags INTERFACE -stdlib=libc++)
    target_link_options(azin_compiler_flags INTERFACE -stdlib=libc++)
endif ()

if (CMAKE_CXX_COMPILER_ID STREQUAL "GNU")
    target_link_libraries(azin_compiler_flags INTERFACE stdc++exp)
endif ()