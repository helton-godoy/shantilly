if(NOT DEFINED EXECUTABLE OR NOT DEFINED EXPECTED_RESULT)
  message(FATAL_ERROR "EXECUTABLE and EXPECTED_RESULT are required")
endif()

set(command "${EXECUTABLE}")
if(DEFINED ARGUMENT)
  list(APPEND command "${ARGUMENT}")
endif()
if(DEFINED TERMINATING_ARGUMENT)
  list(APPEND command "${TERMINATING_ARGUMENT}")
endif()

execute_process(
  COMMAND "${CMAKE_COMMAND}" -E env QT_QPA_PLATFORM=offscreen ${command}
  RESULT_VARIABLE result
  OUTPUT_VARIABLE actual_stdout
  ERROR_VARIABLE actual_stderr
  TIMEOUT 5)

if(NOT result EQUAL EXPECTED_RESULT)
  message(FATAL_ERROR
          "${ARGUMENT} exited with ${result}; expected ${EXPECTED_RESULT}. stderr: ${actual_stderr}")
endif()

foreach(channel stdout stderr)
  if(DEFINED EXPECTED_${channel})
    set(actual "${actual_${channel}}")
    set(expected "${EXPECTED_${channel}}")
    string(REGEX REPLACE "[\r\n]+$" "" actual "${actual}")
    string(REGEX REPLACE "[\r\n]+$" "" expected "${expected}")
    if(NOT actual STREQUAL expected)
      message(FATAL_ERROR
              "Unexpected ${channel} for ${ARGUMENT}.\nExpected: ${expected}\nActual: ${actual}")
    endif()
  endif()
endforeach()
