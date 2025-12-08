cmake_minimum_required(VERSION 3.15)
project(MicroAuthContract)

# Set C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Compiler flags
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wextra -O3")

# Source files
set(SOURCES
    turboauth.cpp
)

set(HEADERS
    turboauth.hpp
)

# Create library
add_library(turboauth STATIC ${SOURCES} ${HEADERS})

# Include directories
target_include_directories(turboauth PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})

# TODO: Add Qubic SDK dependencies when available
# find_package(QubicSDK REQUIRED)
# target_link_libraries(turboauth QubicSDK::Core)

# Installation
install(TARGETS turboauth DESTINATION lib)
install(FILES ${HEADERS} DESTINATION include/turboauth)
