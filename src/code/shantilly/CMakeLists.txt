cmake_minimum_required(VERSION 3.16)

project(SHantilly-cli)

# Configuração Automática do Qt (MOC, UIC, RCC)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTOUIC ON)
set(CMAKE_AUTORCC ON)

# Encontrar fontes
file(GLOB_RECURSE SRCS "*.cpp" "*.cc" "*.h")
list(FILTER SRCS EXCLUDE REGEX ".*/obj/.*")
# Keep the characterized CLI entry point active until Phase 2 extracts it from
# the legacy directory. The static demonstration is not the production CLI.
list(FILTER SRCS EXCLUDE REGEX ".*/main.cpp")

add_executable(shantilly ${SRCS})

# Dependências
target_link_libraries(
  shantilly PRIVATE Qt6::Widgets Qt6::Charts Qt6::PrintSupport Qt6::Network
                    SHantilly-ui)

# Includes
target_include_directories(
  shantilly
  PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}
          ${CMAKE_CURRENT_SOURCE_DIR}/legacy/v2_incomplete
          ${CMAKE_CURRENT_SOURCE_DIR}/legacy/v1_monolith)
