cmake_minimum_required(VERSION 3.16)

# Enable testing (already called in top-level CMakeLists) Add GoogleTest via
# vcpkg (gtest is listed in vcpkg.json)
find_package(GTest CONFIG REQUIRED)

add_executable(unit_tests test_widget_configs.cpp)

target_link_libraries(unit_tests PRIVATE GTest::gtest GTest::gtest_main
                                         Qt6::Widgets SHantilly-ui)

enable_testing()
add_test(NAME unit_tests COMMAND unit_tests)

set(COMPATIBILITY_DIR ${CMAKE_CURRENT_SOURCE_DIR}/compatibility)

add_test(
  NAME compatibility_cli_help
  COMMAND
    ${CMAKE_COMMAND} -DEXECUTABLE=$<TARGET_FILE:shantilly> -DARGUMENT=--help
    -DEXPECTED_FILE=${COMPATIBILITY_DIR}/fixtures/help.txt -P
    ${COMPATIBILITY_DIR}/assert_cli_output.cmake)

add_test(
  NAME compatibility_cli_version
  COMMAND
    ${CMAKE_COMMAND} -DEXECUTABLE=$<TARGET_FILE:shantilly>
    -DARGUMENT=--version
    -DEXPECTED_FILE=${COMPATIBILITY_DIR}/fixtures/version.txt -P
    ${COMPATIBILITY_DIR}/assert_cli_output.cmake)

add_test(
  NAME compatibility_stdin_query
  COMMAND
    ${CMAKE_COMMAND} -DEXECUTABLE=$<TARGET_FILE:shantilly>
    -DINPUT_FILE=${COMPATIBILITY_DIR}/fixtures/minimal_query.input
    -DEXPECTED_OUTPUT=compat_checked=1 -P
    ${COMPATIBILITY_DIR}/assert_stdin_query.cmake)

set_tests_properties(
  compatibility_cli_help compatibility_cli_version compatibility_stdin_query
  PROPERTIES LABELS "compatibility" TIMEOUT 10)
