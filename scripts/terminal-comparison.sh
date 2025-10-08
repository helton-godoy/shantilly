#!/bin/bash

# Terminal Comparison Test for Shantilly TUI
# Tests rendering differences between terminals and identifies potential issues

echo "=== SHANTILLY TERMINAL COMPARISON TEST ==="
echo "Test Date: $(date)"
echo "Running in: ${TERM_PROGRAM:-${TERMINAL:-Unknown Terminal}}"
echo

# Function to test character rendering
test_rendering() {
    local description="$1"
    local chars="$2"

    echo "--- $description ---"
    echo "Characters: $chars"
    echo "Rendered:"
    echo "$chars"
    echo "Length reported by shell: ${#chars}"
    printf "Width test: [%s]\n" "$chars"
    echo
}

echo "=== CHARACTER RENDERING TESTS ==="

# Test 1: Box Drawing Characters (Unicode)
test_rendering "Unicode Box Drawing (RoundedBorder)" "╭─────╮
│     │
╰─────╯"

# Test 2: ASCII Box Characters (Normal Border)
test_rendering "ASCII Box Drawing (NormalBorder)" "+-----+
|     |
+-----+"

# Test 3: Mixed Unicode Characters
test_rendering "Mixed Unicode Symbols" "█▓▒░ ←↑→↓ ✓✗ ★☆"

# Test 4: Lipgloss Border Styles Visual Test
echo "--- LIPGLOSS BORDER STYLES COMPARISON ---"
echo "1. NormalBorder (Current - ASCII):"
cat << 'EOF'
+------------------------+
| Component Content Here |
+------------------------+
EOF

echo
echo "2. RoundedBorder (Previous - Unicode):"
cat << 'EOF'
╭────────────────────────╮
│ Component Content Here │
╰────────────────────────╯
EOF

echo
echo "3. ThickBorder (Unicode):"
cat << 'EOF'
┏━━━━━━━━━━━━━━━━━━━━━━━━┓
┃ Component Content Here ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━┛
EOF

echo
echo "4. DoubleBorder (Unicode):"
cat << 'EOF'
╔════════════════════════╗
║ Component Content Here ║
╚════════════════════════╝
EOF

echo

# Test 5: Color rendering
echo "--- COLOR SUPPORT TEST ---"
echo "Testing 256-color support:"
for i in {0..15}; do
    printf "\033[48;5;${i}m  \033[0m"
done
echo
echo

echo "Testing truecolor support:"
for r in 255 200 150 100 50 0; do
    printf "\033[48;2;${r};100;150m  \033[0m"
done
echo
echo

# Test 6: Terminal-specific issues
echo "--- TERMINAL-SPECIFIC ISSUES ---"
echo "Checking for common terminal rendering problems:"

echo "Test A: Zero-width characters"
printf "Before[%s]After\n" "‌‍"  # Zero-width non-joiner and joiner

echo "Test B: Combining characters"
printf "Combined: a%s e%s\n" "̊" "́"  # Ring above, acute accent

echo "Test C: Double-width characters"
printf "Japanese: 日本語 [Width test]\n"
printf "Emoji: 🎯📋✅ [Width test]\n"

echo

# Test 7: Environment that might affect rendering
echo "--- RENDERING ENVIRONMENT ---"
echo "Terminal size: ${COLUMNS}x${LINES}"
echo "Font info (if available):"
# Try to get font info (works in some terminals)
printf '\033]11;?\033\\' 2>/dev/null || echo "Font query not supported"
sleep 0.1

echo
echo "Cursor position test:"
printf "Start[%s]End\n" "$(printf '\033[6n' 2>/dev/null || echo 'No cursor info')"

echo

# Test 8: Performance test
echo "--- RENDERING PERFORMANCE ---"
echo "Testing rendering speed for different border types..."

time_start=$(date +%s%N)
for i in {1..100}; do
    printf "+-----+\n|     |\n+-----+\n" > /dev/null
done
time_ascii=$(($(date +%s%N) - time_start))

time_start=$(date +%s%N)
for i in {1..100}; do
    printf "╭─────╮\n│     │\n╰─────╯\n" > /dev/null
done
time_unicode=$(($(date +%s%N) - time_start))

echo "ASCII borders:   ${time_ascii} nanoseconds (100 iterations)"
echo "Unicode borders: ${time_unicode} nanoseconds (100 iterations)"
echo "Performance ratio: $(echo "scale=2; $time_unicode / $time_ascii" | bc 2>/dev/null || echo "N/A")"

echo

# Test 9: Shantilly-specific test
echo "--- SHANTILLY APP TEST ---"
echo "Testing with actual shantilly binary..."

if [ -f "./bin/shantilly" ]; then
    echo "Binary found. Testing basic execution:"

    # Test with minimal example
    if [ -f "docs/examples/radiogroup-test.yaml" ]; then
        echo "Running quick test (will exit automatically after 2 seconds)..."
        timeout 2s ./bin/shantilly form docs/examples/radiogroup-test.yaml 2>&1 || true
        echo "Test completed."
    else
        echo "Test file not found: docs/examples/radiogroup-test.yaml"
    fi
else
    echo "Shantilly binary not found. Run 'make build' first."
fi

echo

# Test 10: Recommendations
echo "--- RECOMMENDATIONS ---"
echo "Based on this test, here are recommendations for terminal compatibility:"

if command -v tput >/dev/null; then
    colors=$(tput colors 2>/dev/null)
    if [ "$colors" -ge 256 ]; then
        echo "✅ Color support: Excellent ($colors colors)"
    else
        echo "⚠️  Color support: Limited ($colors colors)"
    fi
else
    echo "❌ tput not available - cannot determine color support"
fi

if [ "$COLORTERM" = "truecolor" ]; then
    echo "✅ Truecolor: Supported"
else
    echo "⚠️  Truecolor: Not detected"
fi

if [ "${LANG}" = *"UTF-8"* ]; then
    echo "✅ UTF-8: Properly configured"
else
    echo "⚠️  UTF-8: May not be configured (LANG=$LANG)"
fi

echo
echo "--- TROUBLESHOOTING TIPS ---"
echo "If you experience layout issues:"
echo "1. Ensure TERM=xterm-256color (current: $TERM)"
echo "2. Use UTF-8 locale: export LANG=en_US.UTF-8"
echo "3. Try in a clean shell: env -i sh -c 'export TERM=xterm-256color; ./bin/shantilly ...'"
echo "4. Check terminal font supports box-drawing characters"
echo "5. Consider using ASCII-only mode if Unicode issues persist"

echo
echo "=== END TERMINAL COMPARISON TEST ==="
