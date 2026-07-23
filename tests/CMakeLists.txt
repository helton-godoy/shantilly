cmake_minimum_required(VERSION 3.16)

# Enable testing (already called in top-level CMakeLists) Add GoogleTest via
# vcpkg (gtest is listed in vcpkg.json)
find_package(GTest CONFIG REQUIRED)

add_executable(unit_tests test_widget_configs.cpp)

target_link_libraries(unit_tests PRIVATE GTest::gtest GTest::gtest_main
                                         Qt6::Widgets shantilly-ui)

enable_testing()
add_test(NAME unit_tests COMMAND unit_tests)

set(LEGACY_RUNTIME_DIR ${PROJECT_SOURCE_DIR}/src/code/shantilly)
add_executable(
  dialog_interaction_tests
  qt_test_main.cpp
  test_dialog_interactions.cpp
  ${LEGACY_RUNTIME_DIR}/SHantilly.cc
  ${LEGACY_RUNTIME_DIR}/SHantilly.h
  ${LEGACY_RUNTIME_DIR}/legacy/v1_monolith/dialog_parser.cc
  ${LEGACY_RUNTIME_DIR}/legacy/v1_monolith/dialog_private.cc
  ${LEGACY_RUNTIME_DIR}/legacy/v1_monolith/dialog_set_options.cc
  ${LEGACY_RUNTIME_DIR}/legacy/v1_monolith/dialog_slots.cc)
set_target_properties(dialog_interaction_tests PROPERTIES AUTOMOC ON)
target_include_directories(
  dialog_interaction_tests
  PRIVATE ${LEGACY_RUNTIME_DIR}
          ${LEGACY_RUNTIME_DIR}/legacy/v1_monolith
          ${PROJECT_SOURCE_DIR}/libs/shantilly-ui/include)
target_link_libraries(
  dialog_interaction_tests
  PRIVATE GTest::gtest Qt6::Test Qt6::Widgets Qt6::Charts shantilly-ui)

add_test(NAME dialog_interaction_tests COMMAND dialog_interaction_tests)
set_tests_properties(
  dialog_interaction_tests
  PROPERTIES LABELS "compatibility;interaction"
             ENVIRONMENT "QT_QPA_PLATFORM=offscreen"
             TIMEOUT 10)

add_test(
  NAME compatibility_fifo_background
  COMMAND
    ${CMAKE_CURRENT_SOURCE_DIR}/compatibility/assert_fifo_input.sh
    $<TARGET_FILE:shantilly>
    ${CMAKE_CURRENT_SOURCE_DIR}/compatibility/fixtures/fifo_background.txt)
set_tests_properties(
  compatibility_fifo_background
  PROPERTIES LABELS "compatibility;fifo"
             TIMEOUT 10)

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

foreach(alias IN ITEMS -h -v)
  if(alias STREQUAL "-h")
    set(fixture help.txt)
    set(test_name compatibility_cli_help_short)
  else()
    set(fixture version.txt)
    set(test_name compatibility_cli_version_short)
  endif()
  add_test(
    NAME ${test_name}
    COMMAND
      ${CMAKE_COMMAND} -DEXECUTABLE=$<TARGET_FILE:shantilly>
      -DARGUMENT=${alias}
      -DEXPECTED_FILE=${COMPATIBILITY_DIR}/fixtures/${fixture} -P
      ${COMPATIBILITY_DIR}/assert_cli_output.cmake)
endforeach()

add_test(
  NAME compatibility_cli_invalid_option
  COMMAND
    ${CMAKE_COMMAND} -DEXECUTABLE=$<TARGET_FILE:shantilly>
    -DARGUMENT=--invalid-option -DEXPECTED_RESULT=1
    "-DEXPECTED_stderr=Error: Unrecognized option --invalid-option" -P
    ${COMPATIBILITY_DIR}/assert_cli_result.cmake)

foreach(option IN ITEMS --resizable -r)
  string(REPLACE "--" "" suffix ${option})
  string(REPLACE "-" "short_" suffix ${suffix})
  add_test(
    NAME compatibility_cli_${suffix}
    COMMAND
      ${CMAKE_COMMAND} -DEXECUTABLE=$<TARGET_FILE:shantilly>
      -DARGUMENT=${option} -DTERMINATING_ARGUMENT=--version
      -DEXPECTED_RESULT=0 -P
      ${COMPATIBILITY_DIR}/assert_cli_result.cmake)
endforeach()

add_test(
  NAME compatibility_stdin_query
  COMMAND
    ${CMAKE_COMMAND} -DEXECUTABLE=$<TARGET_FILE:shantilly>
    -DINPUT_FILE=${COMPATIBILITY_DIR}/fixtures/minimal_query.input
    -DEXPECTED_OUTPUT=compat_checked=1 -P
    ${COMPATIBILITY_DIR}/assert_stdin_query.cmake)

add_test(
  NAME compatibility_protocol_mutations
  COMMAND
    ${CMAKE_COMMAND} -DEXECUTABLE=$<TARGET_FILE:shantilly>
    -DINPUT_FILE=${COMPATIBILITY_DIR}/fixtures/protocol_mutations.input
    -DEXPECTED_FILE=${COMPATIBILITY_DIR}/fixtures/protocol_mutations.txt -P
    ${COMPATIBILITY_DIR}/assert_stdin_output.cmake)

foreach(scenario
        IN ITEMS
           clear_semantics
           layout_widgets
           listbox_progressbar
           nested_tabs
           passive_widgets)
  add_test(
    NAME compatibility_${scenario}
    COMMAND
      ${CMAKE_COMMAND} -DEXECUTABLE=$<TARGET_FILE:shantilly>
      -DINPUT_FILE=${COMPATIBILITY_DIR}/fixtures/${scenario}.input
      -DEXPECTED_FILE=${COMPATIBILITY_DIR}/fixtures/${scenario}.txt -P
      ${COMPATIBILITY_DIR}/assert_stdin_output.cmake)
endforeach()

set_tests_properties(
  compatibility_cli_help compatibility_cli_help_short compatibility_cli_version
  compatibility_cli_version_short compatibility_cli_invalid_option
  compatibility_cli_resizable compatibility_cli_short_r compatibility_stdin_query
  compatibility_protocol_mutations compatibility_clear_semantics
  compatibility_layout_widgets compatibility_listbox_progressbar
  compatibility_nested_tabs compatibility_passive_widgets
  compatibility_fifo_background
  PROPERTIES LABELS "compatibility" TIMEOUT 10)

add_test(
  NAME naming_convention
  COMMAND
    ${CMAKE_COMMAND} -DSOURCE_DIR=${CMAKE_SOURCE_DIR} -P
    ${CMAKE_CURRENT_SOURCE_DIR}/quality/check_naming.cmake)
set_tests_properties(naming_convention PROPERTIES LABELS "quality")
