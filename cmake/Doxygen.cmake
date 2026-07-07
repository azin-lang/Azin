find_package(Doxygen)

if(NOT DOXYGEN_FOUND)
	message(STATUS "Doxygen not found; documentation target disabled.")
	return()
endif()

configure_file(
		${PROJECT_SOURCE_DIR}/Doxyfile.in
		${PROJECT_BINARY_DIR}/Doxyfile
		@ONLY
)

add_custom_target(docs
                  COMMAND ${DOXYGEN_EXECUTABLE} ${PROJECT_BINARY_DIR}/Doxyfile
                  WORKING_DIRECTORY ${PROJECT_BINARY_DIR}
                  COMMENT "Generating API documentation"
)