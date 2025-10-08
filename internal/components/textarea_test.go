package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTextArea(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name        string
		cfg         config.ComponentConfig
		expectError bool
		validate    func(*testing.T, *TextArea)
	}{
		{
			name: "valid textarea creation",
			cfg: config.ComponentConfig{
				Name:        "description",
				Type:        config.TypeTextArea,
				Label:       "Description",
				Required:    true,
				Placeholder: "Enter description",
				Help:        "Help text",
				Default:     "Initial value",
			},
			expectError: false,
			validate: func(t *testing.T, ta *TextArea) {
				assert.Equal(t, "description", ta.name)
				assert.Equal(t, "Description", ta.label)
				assert.True(t, ta.required)
				assert.Equal(t, "Help text", ta.help)
				assert.Equal(t, "Initial value", ta.Value())
			},
		},
		{
			name: "with validation options",
			cfg: config.ComponentConfig{
				Name: "bio",
				Type: config.TypeTextArea,
				Options: map[string]interface{}{
					"min_length": 10,
					"max_length": 500,
					"height":     8,
					"width":      60,
				},
			},
			expectError: false,
			validate: func(t *testing.T, ta *TextArea) {
				assert.Equal(t, 10, ta.minLength)
				assert.Equal(t, 500, ta.maxLength)
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
		{
			name: "default is not a string",
			cfg: config.ComponentConfig{
				Name:    "test",
				Type:    config.TypeTextArea,
				Default: 123,
			},
			expectError: false,
			validate: func(t *testing.T, ta *TextArea) {
				assert.Equal(t, "", ta.Value(), "Default should be ignored and value should be empty")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta, err := NewTextArea(tt.cfg, theme)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, ta)
			} else {
				require.NoError(t, err)
				require.NotNil(t, ta)
				if tt.validate != nil {
					tt.validate(t, ta)
				}
			}
		})
	}
}

func TestTextArea_Init(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name: "test",
		Type: config.TypeTextArea,
	}

	ta, err := NewTextArea(cfg, theme)
	require.NoError(t, err)

	cmd := ta.Init()
	assert.Nil(t, cmd)
}

func TestTextArea_Update(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name: "test",
		Type: config.TypeTextArea,
	}

	ta, err := NewTextArea(cfg, theme)
	require.NoError(t, err)

	t.Run("window size message adjusts width", func(t *testing.T) {
		wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
		model, cmd := ta.Update(wsMsg)
		assert.Equal(t, ta, model)
		assert.Nil(t, cmd)
	})

	t.Run("focused update processes key press and clears error", func(t *testing.T) {
		ta.SetFocus(true)
		ta.SetError("An error")
		keyMsg := tea.KeyPressMsg{}
		model, cmd := ta.Update(keyMsg)
		assert.Equal(t, ta, model)
		assert.Nil(t, cmd)
		assert.Empty(t, ta.errorMsg, "Error should be cleared on typing")
	})

	t.Run("unfocused update does not process key press", func(t *testing.T) {
		initialValue := ta.Value()
		ta.SetFocus(false)
		ta.SetError("Test error")
		keyMsg := tea.KeyPressMsg{}
		model, cmd := ta.Update(keyMsg)
		assert.Equal(t, ta, model)
		assert.Nil(t, cmd)
		assert.Equal(t, "Test error", ta.errorMsg, "Error should remain when not focused")
		assert.Equal(t, initialValue, ta.Value(), "Value should not change when not focused")
	})

	t.Run("non-key message does not clear error when focused", func(t *testing.T) {
		ta.SetFocus(true)
		ta.SetError("Persistent error")
		otherMsg := struct{}{}
		model, cmd := ta.Update(otherMsg)
		assert.Equal(t, ta, model)
		assert.Nil(t, cmd)
		assert.Equal(t, "Persistent error", ta.GetError(), "Error should not be cleared by non-key messages")
	})
}

func TestTextArea_View(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name        string
		cfg         config.ComponentConfig
		setup       func(*TextArea)
		contains    []string
		notContains []string
	}{
		{
			name: "with label and help and no error",
			cfg: config.ComponentConfig{
				Name:  "test",
				Type:  config.TypeTextArea,
				Label: "Test Label",
				Help:  "Helpful text",
			},
			contains:    []string{"Test Label", "Helpful text"},
			notContains: []string{"✗"},
		},
		{
			name: "with error message, hiding help text",
			cfg: config.ComponentConfig{
				Name:  "test",
				Type:  config.TypeTextArea,
				Label: "Test Label",
				Help:  "This help should not be visible",
			},
			setup: func(ta *TextArea) {
				ta.SetError("Test error")
			},
			contains:    []string{"Test Label", "✗ Test error"},
			notContains: []string{"This help should not be visible"},
		},
		{
			name: "no label, no help, no error",
			cfg: config.ComponentConfig{
				Name: "test",
				Type: config.TypeTextArea,
			},
			setup: func(ta *TextArea) {
				ta.SetValue("Just the input")
			},
			contains:    []string{"Just the input"},
			notContains: []string{"Test Label", "Helpful text", "✗"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta, err := NewTextArea(tt.cfg, theme)
			require.NoError(t, err)

			if tt.setup != nil {
				tt.setup(ta)
			}

			view := ta.View()
			for _, expected := range tt.contains {
				assert.Contains(t, view, expected)
			}
			for _, notExpected := range tt.notContains {
				assert.NotContains(t, view, notExpected)
			}
		})
	}
}

func TestTextArea_ComponentInterface(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name: "test-textarea",
		Type: config.TypeTextArea,
	}

	ta, err := NewTextArea(cfg, theme)
	require.NoError(t, err)

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, "test-textarea", ta.Name())
	})

	t.Run("CanFocus", func(t *testing.T) {
		assert.True(t, ta.CanFocus())
	})

	t.Run("SetFocus", func(t *testing.T) {
		assert.False(t, ta.focused)
		ta.SetFocus(true)
		assert.True(t, ta.focused)
		ta.SetFocus(false)
		assert.False(t, ta.focused)
	})

	t.Run("GetError and SetError", func(t *testing.T) {
		assert.Empty(t, ta.GetError())
		ta.SetError("Test error")
		assert.Equal(t, "Test error", ta.GetError())
	})

	t.Run("Value and SetValue", func(t *testing.T) {
		err = ta.SetValue("Test value")
		assert.NoError(t, err)
		assert.Equal(t, "Test value", ta.Value())
	})

	t.Run("SetValue with invalid type", func(t *testing.T) {
		err = ta.SetValue(123)
		assert.Error(t, err)
		assert.Equal(t, "valor inválido: esperado string, recebido int", err.Error())
	})
}

func TestTextArea_Validation(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name     string
		cfg      config.ComponentConfig
		value    string
		isValid  bool
		errorMsg string
	}{
		{
			name: "required field with value",
			cfg: config.ComponentConfig{
				Name:     "test",
				Type:     config.TypeTextArea,
				Required: true,
			},
			value:   "Some text",
			isValid: true,
		},
		{
			name: "required field without value",
			cfg: config.ComponentConfig{
				Name:     "test",
				Type:     config.TypeTextArea,
				Required: true,
			},
			value:    "",
			isValid:  false,
			errorMsg: "Este campo é obrigatório",
		},
		{
			name: "required field with only whitespace",
			cfg: config.ComponentConfig{
				Name:     "test",
				Type:     config.TypeTextArea,
				Required: true,
			},
			value:    "   ",
			isValid:  false,
			errorMsg: "Este campo é obrigatório",
		},
		{
			name: "optional field without value is valid",
			cfg: config.ComponentConfig{
				Name:     "test",
				Type:     config.TypeTextArea,
				Required: false,
			},
			value:   "",
			isValid: true,
		},
		{
			name: "min length validation - valid",
			cfg: config.ComponentConfig{
				Name: "test",
				Type: config.TypeTextArea,
				Options: map[string]interface{}{
					"min_length": 5,
				},
			},
			value:   "Hello",
			isValid: true,
		},
		{
			name: "min length validation - invalid",
			cfg: config.ComponentConfig{
				Name: "test",
				Type: config.TypeTextArea,
				Options: map[string]interface{}{
					"min_length": 10,
				},
			},
			value:    "Short",
			isValid:  false,
			errorMsg: "Mínimo de 10 caracteres",
		},
		{
			name: "max length validation - valid",
			cfg: config.ComponentConfig{
				Name: "test",
				Type: config.TypeTextArea,
				Options: map[string]interface{}{
					"max_length": 10,
				},
			},
			value:   "Short text",
			isValid: true,
		},
		{
			name: "max length validation - invalid",
			cfg: config.ComponentConfig{
				Name: "test",
				Type: config.TypeTextArea,
				Options: map[string]interface{}{
					"max_length": 5,
				},
			},
			value:    "This is too long",
			isValid:  false,
			errorMsg: "Máximo de 5 caracteres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta, err := NewTextArea(tt.cfg, theme)
			require.NoError(t, err)

			err = ta.SetValue(tt.value)
			require.NoError(t, err)

			isValid := ta.IsValid()
			assert.Equal(t, tt.isValid, isValid)

			if !tt.isValid {
				assert.Equal(t, tt.errorMsg, ta.GetError())
			} else {
				assert.Empty(t, ta.GetError())
			}
		})
	}
}

func TestTextArea_Reset(t *testing.T) {
	t.Run("with initial value", func(t *testing.T) {
		theme := styles.DefaultTheme()
		cfg := config.ComponentConfig{
			Name:    "test",
			Type:    config.TypeTextArea,
			Default: "Initial value",
		}
		ta, err := NewTextArea(cfg, theme)
		require.NoError(t, err)

		ta.SetFocus(true)
		ta.SetValue("Changed value")
		ta.SetError("Test error")

		ta.Reset()

		assert.Equal(t, "Initial value", ta.Value())
		assert.Empty(t, ta.GetError())
		assert.False(t, ta.focused)
	})

	t.Run("without initial value", func(t *testing.T) {
		theme := styles.DefaultTheme()
		cfg := config.ComponentConfig{
			Name: "test",
			Type: config.TypeTextArea,
		}
		ta, err := NewTextArea(cfg, theme)
		require.NoError(t, err)

		ta.SetFocus(true)
		ta.SetValue("Changed value")
		ta.SetError("Test error")

		ta.Reset()

		assert.Equal(t, "", ta.Value())
		assert.Empty(t, ta.GetError())
		assert.False(t, ta.focused)
	})
}
