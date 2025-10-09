#!/bin/bash

# Clean Environment Test for Shantilly TUI
# Tests the application in a minimal environment without shell customizations

set -e

echo "=== SHANTILLY CLEAN ENVIRONMENT TEST ==="
echo "This test runs shantilly with minimal environment variables"
echo "to ensure no shell configurations interfere with TUI rendering."
echo

# Save current directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "Project directory: $PROJECT_DIR"
echo "Current environment: $(env | wc -l) variables"
echo

# Check if shantilly binary exists
if [ ! -f "$PROJECT_DIR/bin/shantilly" ]; then
    echo "❌ Shantilly binary not found. Building..."
    cd "$PROJECT_DIR"
    go build -o bin/shantilly ./cmd/shantilly
    echo "✅ Build completed."
fi

# Check if test file exists
TEST_FILE="$PROJECT_DIR/docs/examples/radiogroup-test.yaml"
if [ ! -f "$TEST_FILE" ]; then
    echo "❌ Test file not found: $TEST_FILE"
    echo "Available test files:"
    ls -la "$PROJECT_DIR/docs/examples/" || echo "Examples directory not found"
    exit 1
fi

echo "✅ All prerequisites met."
echo

# Function to run test in different environments
run_test() {
    local env_name="$1"
    local env_vars="$2"
    local description="$3"

    echo "--- TEST: $env_name ---"
    echo "Description: $description"
    echo "Environment: $env_vars"
    echo

    echo "Running for 3 seconds (use Tab/arrows to test navigation)..."
    echo "Press Ctrl+C or Esc to exit early."
    echo

    cd "$PROJECT_DIR"
    echo "🔧 Executando com variáveis de ambiente: $env_vars"
    echo "📊 Definindo SHANTILLY_TEST=1 para ativar modo de teste"
    timeout 3s env -i $env_vars SHANTILLY_TEST=1 ./bin/shantilly form "$TEST_FILE" || true

    echo
    echo "✅ Test '$env_name' completed."
    echo "================================================"
    echo
}

echo "Starting clean environment tests..."
echo "Each test runs for 3 seconds or until you press Esc/Ctrl+C"
echo

# Test 1: Absolute minimal environment
run_test "MINIMAL" \
    "HOME=$HOME PATH=/usr/local/bin:/usr/bin:/bin TERM=xterm-256color" \
    "Bare minimum environment variables"

# Test 2: Standard terminal environment
run_test "STANDARD" \
    "HOME=$HOME PATH=$PATH TERM=xterm-256color LANG=C.UTF-8 LC_ALL=C.UTF-8" \
    "Standard terminal with UTF-8 locale"

# Test 3: Current environment but clean shell
run_test "CLEAN_SHELL" \
    "$(env | grep -E '^(HOME|PATH|TERM|LANG|LC_|DISPLAY|COLORTERM)=' | tr '\n' ' ')" \
    "Current terminal settings without shell customizations"

# Test 4: Force different TERM values
echo "--- TESTING DIFFERENT TERM VALUES ---"
for term_type in "xterm" "xterm-256color" "screen" "linux"; do
    echo "Testing with TERM=$term_type..."
    timeout 2s env -i HOME="$HOME" PATH="$PATH" TERM="$term_type" LANG=C.UTF-8 \
        "$PROJECT_DIR/bin/shantilly" form "$TEST_FILE" 2>&1 | head -5 || true
    echo "✅ $term_type test completed"
    echo
done

echo "=== COMPARISON WITH CURRENT SHELL ==="
echo "Now running in your current shell environment for comparison:"
echo "Environment variables: $(env | wc -l)"
echo "Shell: $SHELL"
echo "Oh My Zsh: ${ZSH:+Enabled}"
echo "Powerlevel10k: $([ -f ~/.p10k.zsh ] && echo 'Enabled' || echo 'Not found')"
echo

timeout 3s "$PROJECT_DIR/bin/shantilly" form "$TEST_FILE" || true

echo
echo "=== TEST RESULTS SUMMARY ==="
echo "✅ All tests completed successfully"
echo
echo "WHAT TO LOOK FOR:"
echo "- Layout should be identical in all test environments"
echo "- Borders should be ASCII characters (+, -, |) in all cases"
echo "- Navigation with Tab/arrows should not break component layout"
echo "- No visual differences between minimal and full environments"
echo
echo "IF YOU SEE DIFFERENCES:"
echo "- Compare the minimal vs current environment tests"
echo "- Note any visual inconsistencies or layout breaks"
echo "- Check if specific TERM values cause issues"
echo
echo "CONCLUSION:"
echo "If all tests show consistent behavior, then your shell"
echo "configurations (Oh My Zsh, Powerlevel10k, etc.) are NOT"
echo "interfering with Shantilly TUI rendering."
echo
echo "=== END CLEAN ENVIRONMENT TEST ==="
