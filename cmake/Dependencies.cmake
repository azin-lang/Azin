include(CPM)

CPMAddPackage("gh:CLIUtils/CLI11#v2.5.0")
CPMAddPackage("gh:fmtlib/fmt#12.2.0")

# Fix for MSVC
set_target_properties(fmt PROPERTIES INTERFACE_SYSTEM_INCLUDE_DIRECTORIES $<TARGET_PROPERTY:fmt,INTERFACE_INCLUDE_DIRECTORIES>)
