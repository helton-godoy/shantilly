package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTextLabel(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name        string
		cfg         config.ComponentConfig
		expectError bool
		validate    func(*testing.T, *TextLabel)
	}{
		{
			name: "valid textlabel creation with label",
			cfg: config.ComponentConfig{
				Name:  "title",
				Type:  config.TypeText,
				Label: "Welcome to our application",
			},
			expectError: false,
			validate: func(t *testing.T, tl *TextLabel) {
				assert.Equal(t, "title", tl.name)
				assert.Equal(t, "Welcome to our application", tl.text)
			},
		},
		{
			name: "textlabel with default value",
			cfg: config.ComponentConfig{
				Name:    "message",
				Type:    config.TypeText,
				Default: "Default message text",
			},
			expectError: false,
			validate: func(t *testing.T, tl *TextLabel) {
				assert.Equal(t, "message", tl.name)
				assert.Equal(t, "Default message text", tl.text)
			},
		},
		{
			name: "textlabel with both label and default (label takes precedence)",
			cfg: config.ComponentConfig{
				Name:    "info",
				Type:    config.TypeText,
				Label:   "Label text",
				Default: "Default text",
			},
			expectError: false,
			validate: func(t *testing.T, tl *TextLabel) {
				assert.Equal(t, "info", tl.name)
				assert.Equal(t, "Label text", tl.text)
			},
		},
		{
			name: "textlabel with empty label and no default",
			cfg: config.ComponentConfig{
				Name:  "empty",
				Type:  config.TypeText,
				Label: "",
			},
			expectError: false,
			validate: func(t *testing.T, tl *TextLabel) {
				assert.Equal(t, "empty", tl.name)
				assert.Equal(t, "", tl.text)
			},
		},
		{
			name: "textlabel with invalid default type",
			cfg: config.ComponentConfig{
				Name:    "invalid-default",
				Type:    config.TypeText,
				Default: 123, // Not a string
			},
			expectError: false,
			validate: func(t *testing.T, tl *TextLabel) {
				assert.Equal(t, "invalid-default", tl.name)
				assert.Equal(t, "", tl.text) // Should default to empty string
			},
		},
		{
			name: "invalid component type",
			cfg: config.ComponentConfig{
				Name: "invalid",
				Type: config.TypeTextInput,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl, err := NewTextLabel(tt.cfg, theme)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, tl)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tl)
				if tt.validate != nil {
					tt.validate(t, tl)
				}
			}
		})
	}
}

func TestTextLabel_Init(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name:  "test",
		Type:  config.TypeText,
		Label: "Test Label",
	}

	tl, err := NewTextLabel(cfg, theme)
	require.NoError(t, err)

	cmd := tl.Init()
	assert.Nil(t, cmd)
}

func TestTextLabel_Update(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name:  "test",
		Type:  config.TypeText,
		Label: "Test Label",
	}

	tl, err := NewTextLabel(cfg, theme)
	require.NoError(t, err)

	// Test various message types - should always return unchanged
	tests := []tea.Msg{
		tea.KeyPressMsg{Text: "a", Code: 'a'},
		tea.KeyPressMsg{Code: tea.KeyEnter},
		tea.WindowSizeMsg{Width: 100, Height: 50},
	}

	for _, msg := range tests {
		model, cmd := tl.Update(msg)
		assert.Equal(t, tl, model)
		assert.Nil(t, cmd)
	}
}

func TestTextLabel_View(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name     string
		cfg      config.ComponentConfig
		expected string
	}{
		{
			name: "simple text label",
			cfg: config.ComponentConfig{
				Name:  "test",
				Type:  config.TypeText,
				Label: "Hello World",
			},
			expected: "Hello World",
		},
		{
			name: "empty text label",
			cfg: config.ComponentConfig{
				Name:  "empty",
				Type:  config.TypeText,
				Label: "",
			},
			expected: "",
		},
		{
			name: "text label with special characters",
			cfg: config.ComponentConfig{
				Name:  "special",
				Type:  config.TypeText,
				Label: "🎉 Special Characters! @#$%",
			},
			expected: "🎉 Special Characters! @#$%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl, err := NewTextLabel(tt.cfg, theme)
			require.NoError(t, err)

			expectedView := theme.Label.Render(tt.expected) + "\n"
			assert.Equal(t, expectedView, tl.View())
		})
	}
}

func TestTextLabel_ComponentInterface(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name:  "test-textlabel",
		Type:  config.TypeText,
		Label: "Test Label",
	}

	tl, err := NewTextLabel(cfg, theme)
	require.NoError(t, err)

	// Test Name
	assert.Equal(t, "test-textlabel", tl.Name())

	// Test CanFocus - should be false for static component
	assert.False(t, tl.CanFocus())

	// Test SetFocus - should be no-op
	tl.SetFocus(true)  // Should not panic
	tl.SetFocus(false) // Should not panic

	// Test IsValid - should always be true
	assert.True(t, tl.IsValid())

	// Test GetError - should always be empty
	assert.Empty(t, tl.GetError())

	// Test SetError - should be no-op
	tl.SetError("Some error")      // Should not panic
	assert.Empty(t, tl.GetError()) // Should still be empty

	// Test Value/SetValue
	assert.Equal(t, "Test Label", tl.Value())

	err = tl.SetValue("New label text")
	assert.NoError(t, err)
	assert.Equal(t, "New label text", tl.Value())

	// Test invalid value type
	err = tl.SetValue(123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "valor inválido: esperado string")

	// Test Reset - should be no-op
	originalText := tl.text
	tl.Reset()                             // Should not panic
	assert.Equal(t, originalText, tl.text) // Should remain unchanged
}

func TestTextLabel_StaticBehavior(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name:  "static",
		Type:  config.TypeText,
		Label: "Static Text",
	}

	tl, err := NewTextLabel(cfg, theme)
	require.NoError(t, err)

	// Test that the component maintains its static nature
	originalText := tl.text

	// Various operations that should not affect the component
	tl.SetFocus(true)
	tl.SetError("Error message")
	tl.Update(tea.KeyPressMsg{Text: "test", Code: 't'})
	tl.Reset()

	// Text should remain unchanged
	assert.Equal(t, originalText, tl.text)

	// Component should still be non-focusable and valid
	assert.False(t, tl.CanFocus())
	assert.True(t, tl.IsValid())
	assert.Empty(t, tl.GetError())
}

func TestTextLabel_MultilineText(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name:  "multiline",
		Type:  config.TypeText,
		Label: "Line 1\nLine 2\nLine 3",
	}

	tl, err := NewTextLabel(cfg, theme)
	require.NoError(t, err)

	view := tl.View()
	assert.Contains(t, view, "Line 1")
	assert.Contains(t, view, "Line 2")
	assert.Contains(t, view, "Line 3")

	// Test that the value correctly preserves newlines
	assert.Equal(t, "Line 1\nLine 2\nLine 3", tl.Value())
}

func TestTextLabel_UncoveredMethods(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name:  "test-coverage",
		Type:  config.TypeText,
		Label: "Coverage Test",
	}

	tl, err := NewTextLabel(cfg, theme)
	require.NoError(t, err)

	// Test SetFocus method (should be no-op)
	tl.SetFocus(true)
	tl.SetFocus(false)
	// No assertions needed - just testing that it doesn't panic

	// Test SetError method (should be no-op)
	tl.SetError("Some error message")
	assert.Empty(t, tl.GetError()) // Should still be empty

	// Test Reset method (should be no-op)
	originalText := tl.text
	tl.Reset()
	assert.Equal(t, originalText, tl.text) // Should remain unchanged
}
