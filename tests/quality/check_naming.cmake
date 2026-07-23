if(NOT DEFINED SOURCE_DIR)
  message(FATAL_ERROR "SOURCE_DIR is required")
endif()

set(scan_roots docs man packaging examples tests src libs)
set(patterns
    "SHantilly[ \t]+--"
    "SHantilly[ \t]*<<"
    "Exec=SHantilly"
    "command:[ \t]+SHantilly"
    "/bin/SHantilly"
    "src/code/SHantilly"
    "build/bin/SHantilly"
    "SHantilly-ui"
    "SHantilly\\.(parser|gui)")

set(violations)
foreach(root IN LISTS scan_roots)
  file(
    GLOB_RECURSE files
    LIST_DIRECTORIES false
    "${SOURCE_DIR}/${root}/*.cc"
    "${SOURCE_DIR}/${root}/*.cpp"
    "${SOURCE_DIR}/${root}/*.desktop"
    "${SOURCE_DIR}/${root}/*.h"
    "${SOURCE_DIR}/${root}/*.md"
    "${SOURCE_DIR}/${root}/*.sh"
    "${SOURCE_DIR}/${root}/*.spec"
    "${SOURCE_DIR}/${root}/*.xml"
    "${SOURCE_DIR}/${root}/*.yaml")

  foreach(file IN LISTS files)
    if(file STREQUAL "${CMAKE_CURRENT_LIST_FILE}")
      continue()
    endif()
    file(READ "${file}" content)
    foreach(pattern IN LISTS patterns)
      if(content MATCHES "${pattern}")
        file(RELATIVE_PATH relative_file "${SOURCE_DIR}" "${file}")
        list(APPEND violations "${relative_file}: ${pattern}")
      endif()
    endforeach()
  endforeach()
endforeach()

if(violations)
  list(JOIN violations "\n" violation_report)
  message(FATAL_ERROR "Invalid SHantilly naming contexts:\n${violation_report}")
endif()
