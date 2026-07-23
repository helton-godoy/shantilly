if(NOT DEFINED EXECUTABLE OR NOT DEFINED ARGUMENT OR NOT DEFINED EXPECTED_FILE)
  message(FATAL_ERROR "EXECUTABLE, ARGUMENT, and EXPECTED_FILE are required")
endif()

file(READ "${EXPECTED_FILE}" expected_output)

execute_process(
  COMMAND "${CMAKE_COMMAND}" -E env QT_QPA_PLATFORM=offscreen "${EXECUTABLE}"
          "${ARGUMENT}"
  RESULT_VARIABLE result
  OUTPUT_VARIABLE actual_output
  ERROR_VARIABLE error_output
  TIMEOUT 5)

string(REGEX REPLACE "[\r\n]+$" "" expected_output "${expected_output}")
string(REGEX REPLACE "[\r\n]+$" "" actual_output "${actual_output}")

if(NOT result EQUAL 0)
  message(FATAL_ERROR
          "${ARGUMENT} exited with ${result}. stderr: ${error_output}")
endif()

if(NOT actual_output STREQUAL expected_output)
  message(FATAL_ERROR
          "Unexpected ${ARGUMENT} output.\nExpected:\n${expected_output}\nActual:\n${actual_output}")
endif()
