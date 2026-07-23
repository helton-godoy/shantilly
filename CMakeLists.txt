cmake_minimum_required(VERSION 3.16)

# Integração automática com vcpkg (Manifest Mode)
if(NOT DEFINED CMAKE_TOOLCHAIN_FILE
   AND EXISTS "$ENV{HOME}/vcpkg/scripts/buildsystems/vcpkg.cmake")
  set(CMAKE_TOOLCHAIN_FILE
      "$ENV{HOME}/vcpkg/scripts/buildsystems/vcpkg.cmake"
      CACHE STRING "")
endif()

project(shantilly-suite LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Configuração Global do Qt (via vcpkg)
find_package(Qt6 REQUIRED COMPONENTS Widgets Charts PrintSupport Network Svg Test)

# Diretórios de Saída
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)

# ============================================================================
# Cobertura de Código (Coverage)
# ============================================================================
option(ENABLE_COVERAGE "Enable coverage reporting" OFF)
if(ENABLE_COVERAGE)
  if(CMAKE_CXX_COMPILER_ID STREQUAL "GNU" OR CMAKE_CXX_COMPILER_ID STREQUAL
                                             "Clang")
    add_compile_options(--coverage -O0 -g)
    add_link_options(--coverage)
  endif()
endif()

# Subprojetos
add_subdirectory(libs/shantilly-ui)
add_subdirectory(src/code/shantilly)
# add_subdirectory(apps/studio) # Habilitaremos quando iniciarmos o Studio

enable_testing()
add_subdirectory(tests)

# poc_modern_cli removido devido à refatoração
# add_executable(
#   poc_modern_cli
#   src/poc_modern_cli.cpp
#   src/code/shantilly/core/parser/ParserMain.cpp
#   src/code/shantilly/core/parser/ParserMain.h
#   src/code/shantilly/core/parser/CommandParser.cpp
#   src/code/shantilly/core/parser/Tokenizer.cpp
#   src/code/shantilly/core/parser/WidgetParserUtils.cpp)

# target_link_libraries(poc_modern_cli PRIVATE Qt6::Widgets shantilly-ui)

# target_include_directories(
#   poc_modern_cli PRIVATE src/code/shantilly/core/parser
#                          src/code/shantilly/core/builder)

# set_property(TARGET poc_modern_cli PROPERTY AUTOMOC ON)

add_executable(visual_verifier src/visual_verifier.cpp)

target_link_libraries(visual_verifier PRIVATE Qt6::Widgets shantilly-ui)

# ============================================================================
# Instalação
# ============================================================================
include(GNUInstallDirs)

install(TARGETS shantilly visual_verifier
        RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR})

install(DIRECTORY libs/shantilly-ui/include/
        DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}/SHantilly)

install(
  TARGETS shantilly-ui
  LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
  ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR})

# ============================================================================
# CPack - Empacotamento DEB
# ============================================================================
set(CPACK_PACKAGE_NAME "shantilly")
set(CPACK_PACKAGE_VENDOR "helton-godoy")
set(CPACK_PACKAGE_DESCRIPTION_SUMMARY
    "SHantilly - Modern UI Builder for Shell Scripts")
set(CPACK_PACKAGE_VERSION "1.0.0")
set(CPACK_DEBIAN_PACKAGE_MAINTAINER "Helton Godoy <helton@example.com>")
set(CPACK_DEBIAN_PACKAGE_DEPENDS "libqt6widgets6, libqt6charts6, libqt6svg6")
set(CPACK_GENERATOR "DEB")

include(CPack)
