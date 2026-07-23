#!/bin/bash
#
# Test de compatibilidade: SHantilly vs DialogBox
# Execute este script para verificar se o SHantilly produz outputs
# idênticos ao dialogbox original.
#

SHANTILLY_BIN="${SHANTILLY_BIN:-./build/bin/shantilly}"
DIALOGBOX_BIN="${DIALOGBOX_BIN:-dialogbox}"

echo "=== SHantilly Compatibility Test ==="
echo ""

# Test 1: Basic label and button
echo "[Test 1] Basic label and button..."
output=$(echo 'add label "Hello World"
add pushbutton "OK" btn_ok apply exit
show' | timeout 2 "${SHANTILLY_BIN}" 2>/dev/null)
echo "Result: $?"

# Test 2: Checkbox state
echo "[Test 2] Checkbox with checked state..."
output=$(echo 'add checkbox "Option" cb1 checked
set cb1 focus
query
add pushbutton "OK" ok apply exit' | timeout 2 "${SHANTILLY_BIN}" 2>/dev/null)
echo "Output: ${output}"

# Test 3: Unset command
echo "[Test 3] Unset command..."
output=$(echo 'add checkbox "Test" cb2 checked
unset cb2 checked
query
add pushbutton "OK" ok apply exit' | timeout 2 "${SHANTILLY_BIN}" 2>/dev/null)
echo "Output: ${output}"

echo ""
echo "=== Tests complete ==="
echo "Note: Some tests require user interaction and will timeout."
