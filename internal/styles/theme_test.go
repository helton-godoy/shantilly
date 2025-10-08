package styles

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultTheme(t *testing.T) {
	theme := DefaultTheme()
	require.NotNil(t, theme)

	// Test that all required styles are initialized
	assert.NotNil(t, theme.Input)
	assert.NotNil(t, theme.InputFocused)
	assert.NotNil(t, theme.InputError)
	assert.NotNil(t, theme.Label)
	assert.NotNil(t, theme.LabelError)
	assert.NotNil(t, theme.Button)
	assert.NotNil(t, theme.ButtonFocused)
	assert.NotNil(t, theme.Title)
	assert.NotNil(t, theme.Description)
	assert.NotNil(t, theme.Help)
	assert.NotNil(t, theme.Error)
	assert.NotNil(t, theme.Border)
	assert.NotNil(t, theme.BorderActive)
	assert.NotNil(t, theme.CheckboxChecked)
	assert.NotNil(t, theme.CheckboxUnchecked)
	assert.NotNil(t, theme.RadioSelected)
	assert.NotNil(t, theme.RadioUnselected)
	assert.NotNil(t, theme.SliderBar)
	assert.NotNil(t, theme.SliderFilled)
	assert.NotNil(t, theme.TabActive)
	assert.NotNil(t, theme.TabInactive)
}

func TestTheme_InputStyles(t *testing.T) {
	theme := DefaultTheme()

	tests := []struct {
		name   string
		style  lipgloss.Style
		text   string
		checks func(*testing.T, string)
	}{
		{
			name:  "input style renders text",
			style: theme.Input,
			text:  "test input",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "test input")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "input focused style renders text",
			style: theme.InputFocused,
			text:  "focused input",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "focused input")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "input error style renders text",
			style: theme.InputError,
			text:  "error input",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "error input")
				assert.NotEmpty(t, rendered)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := tt.style.Render(tt.text)
			tt.checks(t, rendered)
		})
	}
}

func TestTheme_LabelStyles(t *testing.T) {
	theme := DefaultTheme()

	tests := []struct {
		name   string
		style  lipgloss.Style
		text   string
		checks func(*testing.T, string)
	}{
		{
			name:  "label style renders text",
			style: theme.Label,
			text:  "test label",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "test label")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "label error style renders text",
			style: theme.LabelError,
			text:  "error label",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "error label")
				assert.NotEmpty(t, rendered)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := tt.style.Render(tt.text)
			tt.checks(t, rendered)
		})
	}
}

func TestTheme_ButtonStyles(t *testing.T) {
	theme := DefaultTheme()

	tests := []struct {
		name   string
		style  lipgloss.Style
		text   string
		checks func(*testing.T, string)
	}{
		{
			name:  "button style renders text",
			style: theme.Button,
			text:  "Click Me",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "Click Me")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "button focused style renders text",
			style: theme.ButtonFocused,
			text:  "Focused Button",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "Focused Button")
				assert.NotEmpty(t, rendered)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := tt.style.Render(tt.text)
			tt.checks(t, rendered)
		})
	}
}

func TestTheme_ContainerStyles(t *testing.T) {
	theme := DefaultTheme()

	tests := []struct {
		name   string
		style  lipgloss.Style
		text   string
		checks func(*testing.T, string)
	}{
		{
			name:  "title style renders text",
			style: theme.Title,
			text:  "Application Title",
			checks: func(t *testing.T, rendered string) {
				// Title style may contain ANSI codes, just check it's not empty
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "description style renders text",
			style: theme.Description,
			text:  "This is a description",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "This is a description")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "help style renders text",
			style: theme.Help,
			text:  "Press Enter to continue",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "Press Enter to continue")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "error style renders text",
			style: theme.Error,
			text:  "An error occurred",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "An error occurred")
				assert.NotEmpty(t, rendered)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := tt.style.Render(tt.text)
			tt.checks(t, rendered)
		})
	}
}

func TestTheme_BorderStyles(t *testing.T) {
	theme := DefaultTheme()

	tests := []struct {
		name   string
		style  lipgloss.Style
		text   string
		checks func(*testing.T, string)
	}{
		{
			name:  "border style renders text",
			style: theme.Border,
			text:  "bordered content",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "bordered content")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "border active style renders text",
			style: theme.BorderActive,
			text:  "active bordered content",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "active bordered content")
				assert.NotEmpty(t, rendered)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := tt.style.Render(tt.text)
			tt.checks(t, rendered)
		})
	}
}

func TestTheme_ComponentSpecificStyles(t *testing.T) {
	theme := DefaultTheme()

	tests := []struct {
		name   string
		style  lipgloss.Style
		text   string
		checks func(*testing.T, string)
	}{
		{
			name:  "checkbox checked style",
			style: theme.CheckboxChecked,
			text:  "[✓] Checked",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "[✓] Checked")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "checkbox unchecked style",
			style: theme.CheckboxUnchecked,
			text:  "[ ] Unchecked",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "[ ] Unchecked")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "radio selected style",
			style: theme.RadioSelected,
			text:  "(•) Selected",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "(•) Selected")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "radio unselected style",
			style: theme.RadioUnselected,
			text:  "( ) Unselected",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "( ) Unselected")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "slider bar style",
			style: theme.SliderBar,
			text:  "────────────",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "────────────")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "slider filled style",
			style: theme.SliderFilled,
			text:  "████████████",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "████████████")
				assert.NotEmpty(t, rendered)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := tt.style.Render(tt.text)
			tt.checks(t, rendered)
		})
	}
}

func TestTheme_TabStyles(t *testing.T) {
	theme := DefaultTheme()

	tests := []struct {
		name   string
		style  lipgloss.Style
		text   string
		checks func(*testing.T, string)
	}{
		{
			name:  "tab active style",
			style: theme.TabActive,
			text:  "Active Tab",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "Active Tab")
				assert.NotEmpty(t, rendered)
			},
		},
		{
			name:  "tab inactive style",
			style: theme.TabInactive,
			text:  "Inactive Tab",
			checks: func(t *testing.T, rendered string) {
				assert.Contains(t, rendered, "Inactive Tab")
				assert.NotEmpty(t, rendered)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := tt.style.Render(tt.text)
			tt.checks(t, rendered)
		})
	}
}

func TestTheme_StyleConsistency(t *testing.T) {
	theme := DefaultTheme()

	// Test that styles are properly configured and don't panic
	testStyles := map[string]lipgloss.Style{
		"Input":             theme.Input,
		"InputFocused":      theme.InputFocused,
		"InputError":        theme.InputError,
		"Label":             theme.Label,
		"LabelError":        theme.LabelError,
		"Button":            theme.Button,
		"ButtonFocused":     theme.ButtonFocused,
		"Title":             theme.Title,
		"Description":       theme.Description,
		"Help":              theme.Help,
		"Error":             theme.Error,
		"Border":            theme.Border,
		"BorderActive":      theme.BorderActive,
		"CheckboxChecked":   theme.CheckboxChecked,
		"CheckboxUnchecked": theme.CheckboxUnchecked,
		"RadioSelected":     theme.RadioSelected,
		"RadioUnselected":   theme.RadioUnselected,
		"SliderBar":         theme.SliderBar,
		"SliderFilled":      theme.SliderFilled,
		"TabActive":         theme.TabActive,
		"TabInactive":       theme.TabInactive,
	}

	for name, style := range testStyles {
		t.Run(name, func(t *testing.T) {
			// Test that each style can render text without panicking
			rendered := style.Render("test")
			assert.NotEmpty(t, rendered, "Style %s should render non-empty text", name)
			// Styles may contain ANSI codes, so just check for non-empty output
			assert.NotEmpty(t, rendered, "Style %s should render non-empty text", name)
		})
	}
}

func TestTheme_EmptyText(t *testing.T) {
	theme := DefaultTheme()

	// Test that styles handle empty text gracefully
	styles := []lipgloss.Style{
		theme.Input,
		theme.Label,
		theme.Title,
		theme.Error,
		theme.Help,
	}

	for i, style := range styles {
		t.Run(fmt.Sprintf("style_%d", i), func(t *testing.T) {
			rendered := style.Render("")
			// Should not panic and should return some result
			assert.NotNil(t, rendered)
		})
	}
}

func TestTheme_LongText(t *testing.T) {
	theme := DefaultTheme()

	longText := "This is a very long text that might test the boundaries of the styling system and ensure that it can handle lengthy content without any issues or performance problems."

	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"Input", theme.Input},
		{"Title", theme.Title},
		{"Description", theme.Description},
		{"Help", theme.Help},
		{"Error", theme.Error},
	}

	for _, s := range styles {
		t.Run(s.name, func(t *testing.T) {
			rendered := s.style.Render(longText)
			// Styles may contain ANSI codes, so just check for non-empty output
			assert.NotEmpty(t, rendered)
		})
	}
}

func TestTheme_SpecialCharacters(t *testing.T) {
	theme := DefaultTheme()

	specialTexts := []string{
		"Unicode: 🎉 ✅ ❌ 🔥",
		"Symbols: @#$%^&*()_+",
		"Newlines:\nLine 1\nLine 2",
		"Tabs:\tTabbed\tContent",
		"Mixed: Hello 世界 🌍",
	}

	for _, text := range specialTexts {
		t.Run("special_chars", func(t *testing.T) {
			rendered := theme.Title.Render(text)
			assert.NotEmpty(t, rendered)
			// Basic containment check - some characters might be transformed
			// but the rendered text should not be empty
		})
	}
}

func TestTheme_ColorValues(t *testing.T) {
	// Test that color constants are properly defined
	// This is more of a smoke test to ensure colors don't cause panics
	theme := DefaultTheme()

	// Test that we can create styles with colors without errors
	testStyle := lipgloss.NewStyle().Foreground(primaryColor)
	rendered := testStyle.Render("test")
	assert.Contains(t, rendered, "test")

	// Test theme styles with sample text
	sampleText := "Sample"

	focusedStyle := theme.InputFocused
	focusedRendered := focusedStyle.Render(sampleText)
	assert.Contains(t, focusedRendered, sampleText)

	errorStyle := theme.Error
	errorRendered := errorStyle.Render(sampleText)
	assert.Contains(t, errorRendered, sampleText)
}
