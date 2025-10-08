package styles

import (
	"github.com/charmbracelet/lipgloss/v2"
)

// Theme contains all Lipgloss styles for the Shantilly TUI.
// It uses simple colors that work well in most terminals.
type Theme struct {
	// Input styles
	Input        lipgloss.Style
	InputFocused lipgloss.Style
	InputError   lipgloss.Style

	// Label styles
	Label      lipgloss.Style
	LabelError lipgloss.Style

	// Button/Action styles
	Button        lipgloss.Style
	ButtonFocused lipgloss.Style

	// Container styles
	Title       lipgloss.Style
	Description lipgloss.Style
	Help        lipgloss.Style
	Error       lipgloss.Style

	// Layout styles
	Border       lipgloss.Style
	BorderActive lipgloss.Style

	// Component-specific styles
	CheckboxChecked   lipgloss.Style
	CheckboxUnchecked lipgloss.Style
	RadioSelected     lipgloss.Style
	RadioUnselected   lipgloss.Style
	SliderBar         lipgloss.Style
	SliderFilled      lipgloss.Style

	// Tab styles
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style
}

// Color palette
var (
	// Primary colors - Charm purple
	primaryColor = lipgloss.Color("#7D56F4")
	primaryDark  = lipgloss.Color("#5A3FBF")

	// Accent colors
	accentGreen = lipgloss.Color("#04B575")
	accentRed   = lipgloss.Color("#FF0000")
	accentBlue  = lipgloss.Color("#0087D7")

	// Neutral colors
	textPrimary   = lipgloss.Color("#E0E0E0")
	textSecondary = lipgloss.Color("#888888")
	textMuted     = lipgloss.Color("#666666")

	// Background colors
	bgNormal  = lipgloss.Color("#1A1A1A")
	bgFocused = lipgloss.Color("#2D2640")
	bgError   = lipgloss.Color("#3D2020")

	// Border colors
	borderNormal = lipgloss.Color("#404040")
	borderFocus  = lipgloss.Color("#7D56F4")
	borderError  = lipgloss.Color("#FF0000")
)

// DefaultTheme creates a theme using default Lipgloss styles.
func DefaultTheme() *Theme {
	t := &Theme{}

	// Input styles (without borders - borders are applied by layout/form models)
	t.Input = lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(bgNormal).
		Padding(0, 1)

	t.InputFocused = t.Input.
		Background(bgFocused)

	t.InputError = t.Input.
		Background(bgError)

	// Label styles
	t.Label = lipgloss.NewStyle().
		Foreground(textPrimary).
		Bold(true).
		MarginBottom(1)

	t.LabelError = t.Label.
		Foreground(accentRed)

	// Button styles
	t.Button = lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(bgNormal).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderNormal).
		Padding(0, 2).
		MarginRight(2)

	t.ButtonFocused = t.Button.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(primaryColor).
		BorderForeground(primaryColor).
		Bold(true)

	// Container styles
	t.Title = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		MarginBottom(1)

	t.Description = lipgloss.NewStyle().
		Foreground(textSecondary).
		Italic(true).
		MarginBottom(2)

	t.Help = lipgloss.NewStyle().
		Foreground(textMuted).
		Italic(true).
		MarginTop(1)

	t.Error = lipgloss.NewStyle().
		Foreground(accentRed).
		Bold(true).
		MarginTop(1)

	// Layout styles
	t.Border = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderNormal).
		Padding(1, 2)

	t.BorderActive = t.Border.
		BorderForeground(borderFocus)

	// Checkbox styles
	t.CheckboxChecked = lipgloss.NewStyle().
		Foreground(accentGreen).
		Bold(true)

	t.CheckboxUnchecked = lipgloss.NewStyle().
		Foreground(textMuted)

	// Radio styles
	t.RadioSelected = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)

	t.RadioUnselected = lipgloss.NewStyle().
		Foreground(textPrimary)

	// Slider styles
	t.SliderBar = lipgloss.NewStyle().
		Foreground(borderNormal)

	t.SliderFilled = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)

	// Tab styles
	t.TabActive = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(primaryColor).
		Padding(0, 2).
		Bold(true)

	t.TabInactive = lipgloss.NewStyle().
		Foreground(textSecondary).
		Background(bgNormal).
		Padding(0, 2)

	return t
}
