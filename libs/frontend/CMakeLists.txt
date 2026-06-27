add_library(frontend_source STATIC
        source/src/source.cpp
        source/include/azc/source.hpp
)

add_library(Azin::Frontend::Source ALIAS frontend_source)

target_include_directories(frontend_source PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/source/include>
        $<INSTALL_INTERFACE:include>
)

target_link_libraries(frontend_source PRIVATE
        Azin::CompilerFlags
)

#-------------------------------------------------------------------------------

add_library(frontend_token INTERFACE)

add_library(Azin::Frontend::Token ALIAS frontend_token)

target_include_directories(frontend_token INTERFACE
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/token/include>
        $<INSTALL_INTERFACE:include>
)

target_link_libraries(frontend_token INTERFACE
        Azin::CompilerFlags
)

#-------------------------------------------------------------------------------

add_library(frontend_lexer STATIC
        lexer/src/lexer.cpp
        lexer/include/azc/lexer.hpp
)

add_library(Azin::Frontend::Lexer ALIAS frontend_lexer)

target_include_directories(frontend_lexer PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/lexer/include>
        $<INSTALL_INTERFACE:include>
)

target_link_libraries(frontend_lexer PUBLIC
        Azin::CompilerFlags
        Azin::Frontend::Token
        fmt::fmt
)

#-------------------------------------------------------------------------------

add_library(frontend_parser STATIC
        parser/src/parser.cpp
        parser/include/azc/parser.hpp
)

add_library(Azin::Frontend::Parser ALIAS frontend_parser)

target_include_directories(frontend_parser PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/parser/include>
        $<INSTALL_INTERFACE:include>
)

target_link_libraries(frontend_parser PRIVATE
        Azin::CompilerFlags
)