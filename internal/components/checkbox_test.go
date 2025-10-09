package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCheckbox(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name        string
		cfg         config.ComponentConfig
		expectError bool
		validate    func(*testing.T, *Checkbox)
	}{
		{
			name: "valid checkbox creation",
			cfg: config.ComponentConfig{
				Name:     "agree",
				Type:     config.TypeCheckbox,
				Label:    "I agree to the terms",
				Required: true,
				Help:     "Please read the terms carefully",
				Default:  false,
			},
			expectError: false,
			validate: func(t *testing.T, c *Checkbox) {
				assert.Equal(t, "agree", c.name)
				assert.Equal(t, "I agree to the terms", c.label)
				assert.True(t, c.required)
				assert.Equal(t, "Please read the terms carefully", c.help)
				assert.False(t, c.checked)
				assert.False(t, c.initialValue)
			},
		},
		{
			name: "checkbox with default true",
			cfg: config.ComponentConfig{
				Name:    "newsletter",
				Type:    config.TypeCheckbox,
				Label:   "Subscribe to newsletter",
				Default: true,
			},
			expectError: false,
			validate: func(t *testing.T, c *Checkbox) {
				assert.Equal(t, "newsletter", c.name)
				assert.True(t, c.checked)
				assert.True(t, c.initialValue)
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
			name: "default is not a boolean",
			cfg: config.ComponentConfig{
				Name:    "test",
				Type:    config.TypeCheckbox,
				Default: "true", // string instead of bool
			},
			expectError: false,
			validate: func(t *testing.T, c *Checkbox) {
				assert.False(t, c.checked, "Default should be ignored if not a bool")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCheckbox(tt.cfg, theme)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, c)
			} else {
				require.NoError(t, err)
				require.NotNil(t, c)
				if tt.validate != nil {
					tt.validate(t, c)
				}
			}
		})
	}
}

func TestCheckbox_Init(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name: "test",
		Type: config.TypeCheckbox,
	}

	c, err := NewCheckbox(cfg, theme)
	require.NoError(t, err)

	cmd := c.Init()
	assert.Nil(t, cmd)
}

func TestCheckbox_Update(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name:  "test",
		Type:  config.TypeCheckbox,
		Label: "Test checkbox",
	}

	c, err := NewCheckbox(cfg, theme)
	require.NoError(t, err)

	spaceMsg := tea.KeyPressMsg{Code: tea.KeySpace}
	enterMsg := tea.KeyPressMsg{Code: tea.KeyEnter}
	otherKeyMsg := tea.KeyPressMsg{Code: 'a'}

	t.Run("should not update when not focused", func(t *testing.T) {
		err := c.SetValue(false)
		require.NoError(t, err)
		model, cmd := c.Update(spaceMsg)
		assert.Equal(t, c, model)
		assert.Nil(t, cmd)
		assert.False(t, c.checked, "State should not change when not focused")
	})

	t.Run("should toggle on space key when focused and clear error", func(t *testing.T) {
		c.SetFocus(true)
		err := c.SetValue(false)
		require.NoError(t, err)
		c.SetError("an error")

		c.Update(spaceMsg)
		assert.True(t, c.checked, "Should be checked after space")
		assert.Empty(t, c.errorMsg, "Error should be cleared on toggle")

		c.Update(spaceMsg)
		assert.False(t, c.checked, "Should be unchecked after second space")
	})

	t.Run("should toggle on enter key when focused", func(t *testing.T) {
		c.SetFocus(true)
		err := c.SetValue(false)
		require.NoError(t, err)
		c.Update(enterMsg)
		assert.True(t, c.checked, "Should be checked after enter")
	})

	t.Run("should not toggle on other keys when focused", func(t *testing.T) {
		c.SetFocus(true)
		err := c.SetValue(true) // Start with a known state
		require.NoError(t, err)
		c.Update(otherKeyMsg)
		assert.True(t, c.checked, "State should not change for other keys")
	})

	t.Run("should not toggle on non-key messages", func(t *testing.T) {
		c.SetFocus(true)
		err := c.SetValue(false)
		require.NoError(t, err)
		msg := struct{}{} // Not a key message
		model, cmd := c.Update(msg)
		assert.False(t, c.checked, "Should not toggle on non-key messages")
		assert.Equal(t, c, model)
		assert.Nil(t, cmd)
	})
}

func TestCheckbox_View(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name        string
		setup       func(*Checkbox)
		contains    []string
		notContains []string
	}{
		{
			name:        "unchecked",
			setup:       func(c *Checkbox) { _ = c.SetValue(false) }, // Error ignored for brevity in tests
			contains:    []string{"[ ] Test Label"},
			notContains: []string{"[✓]"},
		},
		{
			name:        "checked",
			setup:       func(c *Checkbox) { _ = c.SetValue(true) }, // Error ignored for brevity in tests
			contains:    []string{"[✓] Test Label"},
			notContains: []string{"[ ]"},
		},
		{
			name:        "with help text and no error",
			setup:       func(c *Checkbox) {},
			contains:    []string{"Test Label", "This is help text"},
			notContains: []string{"✗"},
		},
		{
			name: "with error message, hiding help text",
			setup: func(c *Checkbox) {
				c.SetError("This field is required")
			},
			contains:    []string{"Test Label", "✗ This field is required"},
			notContains: []string{"This is help text"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.ComponentConfig{
				Name:  "test",
				Type:  config.TypeCheckbox,
				Label: "Test Label",
				Help:  "This is help text",
			}
			c, err := NewCheckbox(cfg, theme)
			require.NoError(t, err)

			tt.setup(c)

			view := c.View()
			for _, expected := range tt.contains {
				assert.Contains(t, view, expected)
			}
			for _, notExpected := range tt.notContains {
				assert.NotContains(t, view, notExpected)
			}
		})
	}
}

func TestCheckbox_ComponentInterface(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{Name: "test-checkbox", Type: config.TypeCheckbox}
	c, err := NewCheckbox(cfg, theme)
	require.NoError(t, err)

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, "test-checkbox", c.Name())
	})

	t.Run("CanFocus", func(t *testing.T) {
		assert.True(t, c.CanFocus())
	})

	t.Run("SetFocus", func(t *testing.T) {
		c.SetFocus(false) // Reset
		assert.False(t, c.focused)
		c.SetFocus(true)
		assert.True(t, c.focused)
		c.SetFocus(false)
		assert.False(t, c.focused)
	})

	t.Run("GetError and SetError", func(t *testing.T) {
		c.SetError("") // Reset
		assert.Empty(t, c.GetError())
		c.SetError("Test error")
		assert.Equal(t, "Test error", c.GetError())
	})

	t.Run("Value and SetValue", func(t *testing.T) {
		err := c.SetValue(true)
		assert.NoError(t, err)
		assert.True(t, c.Value().(bool))
	})

	t.Run("SetValue with invalid type", func(t *testing.T) {
		err := c.SetValue("not a bool")
		assert.Error(t, err)
		assert.Equal(t, "valor inválido: esperado bool, recebido string", err.Error())
	})
}

func TestCheckbox_Validation(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name     string
		cfg      config.ComponentConfig
		checked  bool
		isValid  bool
		errorMsg string
	}{
		{
			name:     "required checkbox checked is valid",
			cfg:      config.ComponentConfig{Name: "test", Type: config.TypeCheckbox, Required: true},
			checked:  true,
			isValid:  true,
			errorMsg: "",
		},
		{
			name:     "required checkbox unchecked is invalid",
			cfg:      config.ComponentConfig{Name: "test", Type: config.TypeCheckbox, Required: true},
			checked:  false,
			isValid:  false,
			errorMsg: "Esta opção deve ser marcada",
		},
		{
			name:     "optional checkbox unchecked is valid",
			cfg:      config.ComponentConfig{Name: "test", Type: config.TypeCheckbox, Required: false},
			checked:  false,
			isValid:  true,
			errorMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCheckbox(tt.cfg, theme)
			require.NoError(t, err)

			err = c.SetValue(tt.checked)
			require.NoError(t, err)
			isValid := c.IsValid()
			assert.Equal(t, tt.isValid, isValid)

			if !tt.isValid {
				assert.Equal(t, tt.errorMsg, c.GetError())
			} else {
				assert.Empty(t, c.GetError())
			}
		})
	}
}

func TestCheckbox_Reset(t *testing.T) {
	theme := styles.DefaultTheme()

	t.Run("when initial value is true", func(t *testing.T) {
		cfg := config.ComponentConfig{Name: "test", Type: config.TypeCheckbox, Default: true}
		c, err := NewCheckbox(cfg, theme)
		require.NoError(t, err)

		c.SetFocus(true)
		err = c.SetValue(false)
		require.NoError(t, err)
		c.SetError("Test error")

		c.Reset()

		assert.True(t, c.checked, "Should reset to initial value (true)")
		assert.Empty(t, c.GetError())
		assert.False(t, c.focused)
	})

	t.Run("when initial value is false", func(t *testing.T) {
		cfg := config.ComponentConfig{Name: "test", Type: config.TypeCheckbox, Default: false}
		c, err := NewCheckbox(cfg, theme)
		require.NoError(t, err)

		c.SetFocus(true)
		err = c.SetValue(true)
		require.NoError(t, err)
		c.SetError("Test error")

		c.Reset()

		assert.False(t, c.checked, "Should reset to initial value (false)")
		assert.Empty(t, c.GetError())
		assert.False(t, c.focused)
	})
}
