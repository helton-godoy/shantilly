if(NOT DEFINED EXECUTABLE OR NOT DEFINED INPUT_FILE OR
   NOT DEFINED EXPECTED_FILE)
  message(FATAL_ERROR "EXECUTABLE, INPUT_FILE, and EXPECTED_FILE are required")
endif()

file(READ "${EXPECTED_FILE}" expected_output)

execute_process(
  COMMAND "${CMAKE_COMMAND}" -E env QT_QPA_PLATFORM=offscreen "${EXECUTABLE}"
          --hidden
  INPUT_FILE "${INPUT_FILE}"
  RESULT_VARIABLE result
  OUTPUT_VARIABLE actual_output
  ERROR_VARIABLE error_output
  TIMEOUT 2)

string(REGEX REPLACE "[\r\n]+$" "" expected_output "${expected_output}")
string(REGEX REPLACE "[\r\n]+$" "" actual_output "${actual_output}")

if(NOT actual_output STREQUAL expected_output)
  message(FATAL_ERROR
          "Unexpected protocol output.\nExpected:\n${expected_output}\nActual:\n${actual_output}\nResult: ${result}\nStderr: ${error_output}")
endif()
