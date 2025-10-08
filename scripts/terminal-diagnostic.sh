#!/bin/bash

# Terminal Diagnostic Tool for Shantilly TUI
# Checks environment variables and terminal capabilities that might affect rendering

echo "=== SHANTILLY TERMINAL DIAGNOSTIC ==="
echo "Generated at: $(date)"
echo

echo "=== BASIC TERMINAL INFO ==="
echo "TERM: ${TERM:-Not set}"
echo "TERMINAL: ${TERMINAL:-Not set}"
echo "TERM_PROGRAM: ${TERM_PROGRAM:-Not set}"
echo "TERM_PROGRAM_VERSION: ${TERM_PROGRAM_VERSION:-Not set}"
echo "COLORTERM: ${COLORTERM:-Not set}"
echo

echo "=== LOCALE AND CHARACTER ENCODING ==="
echo "LANG: ${LANG:-Not set}"
echo "LC_ALL: ${LC_ALL:-Not set}"
echo "LC_CTYPE: ${LC_CTYPE:-Not set}"
echo

echo "=== TERMINAL DIMENSIONS ==="
echo "COLUMNS: ${COLUMNS:-$(tput cols 2>/dev/null || echo 'Unknown')}"
echo "LINES: ${LINES:-$(tput lines 2>/dev/null || echo 'Unknown')}"
echo

echo "=== UNICODE SUPPORT TEST ==="
echo "Testing Unicode characters used by lipgloss.RoundedBorder():"
echo "╭─╮"
echo "│ │"
echo "╰─╯"
echo
echo "Testing ASCII characters used by lipgloss.NormalBorder():"
echo "+--+"
echo "|  |"
echo "+--+"
echo

echo "=== SHELL INFORMATION ==="
echo "Current Shell: ${SHELL:-Unknown}"
echo "Shell Version: $($SHELL --version 2>/dev/null | head -1 || echo 'Unknown')"
echo "Running under: ${0}"
echo

echo "=== OH-MY-ZSH/OH-MY-BASH CHECK ==="
if [ -n "$ZSH" ]; then
    echo "Oh My Zsh detected: $ZSH"
    echo "Oh My Zsh Theme: ${ZSH_THEME:-Not set}"
fi

if [ -n "$OSH" ]; then
    echo "Oh My Bash detected: $OSH"
    echo "Oh My Bash Theme: ${OSH_THEME:-Not set}"
fi

echo

echo "=== PROMPT CUSTOMIZATION CHECK ==="
echo "PS1 length: ${#PS1}"
echo "PROMPT set: ${PROMPT:+Yes}"
echo "Starship: $(command -v starship >/dev/null && echo 'Installed' || echo 'Not found')"
echo "Powerlevel10k: $([ -f ~/.p10k.zsh ] && echo 'Config found' || echo 'Not found')"
echo

echo "=== TERMINAL CAPABILITIES ==="
echo "Colors (tput colors): $(tput colors 2>/dev/null || echo 'Unknown')"
echo "Can change colors: $(tput ccc 2>/dev/null && echo 'Yes' || echo 'No')"
echo "Has status line: $(tput hs 2>/dev/null && echo 'Yes' || echo 'No')"
echo

echo "=== ENVIRONMENT VARIABLES THAT MIGHT AFFECT RENDERING ==="
env | grep -E "(TERM|DISPLAY|COLORTERM|LANG|LC_|COLUMNS|LINES)" | sort
echo

echo "=== FONT AND RENDERING ==="
echo "Testing font rendering with problematic characters:"
printf "Box drawing: ┌─┐│└─┘\n"
printf "Block elements: █▓▒░\n"
printf "Arrows: ←↑→↓\n"
echo

echo "=== RECOMMENDATION ==="
echo "If you see garbled characters above, your terminal may have issues with"
echo "Unicode rendering. The shantilly app now uses ASCII borders for maximum"
echo "compatibility, but some visual elements might still be affected."
echo

echo "=== DEBUGGING TIPS ==="
echo "1. Try different TERM values: xterm, xterm-256color, screen-256color"
echo "2. Ensure UTF-8 locale: export LANG=en_US.UTF-8"
echo "3. Check terminal font supports Unicode box drawing"
echo "4. Test in minimal shell: sh -c './bin/shantilly ...'"
echo

echo "=== END DIAGNOSTIC ==="
