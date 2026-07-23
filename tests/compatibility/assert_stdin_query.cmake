if(NOT DEFINED EXECUTABLE OR NOT DEFINED INPUT_FILE OR NOT DEFINED EXPECTED_OUTPUT)
  message(FATAL_ERROR "EXECUTABLE, INPUT_FILE, and EXPECTED_OUTPUT are required")
endif()

execute_process(
  COMMAND "${CMAKE_COMMAND}" -E env QT_QPA_PLATFORM=offscreen "${EXECUTABLE}"
          --hidden
  INPUT_FILE "${INPUT_FILE}"
  RESULT_VARIABLE result
  OUTPUT_VARIABLE actual_output
  ERROR_VARIABLE error_output
  TIMEOUT 2)

string(FIND "${actual_output}" "${EXPECTED_OUTPUT}" expected_position)
if(expected_position EQUAL -1)
  message(FATAL_ERROR
          "Expected '${EXPECTED_OUTPUT}' on stdout. Result: ${result}. stderr: ${error_output}. stdout: ${actual_output}")
endif()
