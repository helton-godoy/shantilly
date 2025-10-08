#!/usr/bin/env bash
# Coverage Report Validator - Enforces 85% minimum coverage threshold
# Usage: ./scripts/coverage-report.sh [threshold]

set -euo pipefail

THRESHOLD=${1:-85}
COVERAGE_FILE="coverage.out"

if [ ! -f "$COVERAGE_FILE" ]; then
    echo "❌ Error: Coverage file not found: $COVERAGE_FILE"
    echo "Run 'make coverage' first to generate the coverage report."
    exit 1
fi

# Calculate total coverage percentage
COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total: | awk '{print $3}' | sed 's/%//')

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Coverage Report"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Total Coverage: ${COVERAGE}%"
echo "Threshold:      ${THRESHOLD}%"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Compare coverage with threshold (handle decimal comparison)
if (( $(echo "$COVERAGE >= $THRESHOLD" | bc -l) )); then
    echo "✅ Coverage check passed (${COVERAGE}% >= ${THRESHOLD}%)"
    echo ""
    echo "Detailed breakdown:"
    go tool cover -func="$COVERAGE_FILE" | grep -v "total:" | tail -n 10
    exit 0
else
    echo "❌ Coverage check failed (${COVERAGE}% < ${THRESHOLD}%)"
    echo ""
    echo "Files with low coverage:"
    go tool cover -func="$COVERAGE_FILE" | grep -v "total:" | awk -v threshold="$THRESHOLD" '{if ($3+0 < threshold) print}'
    echo ""
    echo "Please add more tests to reach the ${THRESHOLD}% coverage threshold."
    exit 1
fi
