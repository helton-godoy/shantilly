package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupRadioGroup is a helper to create a default RadioGroup for tests.
func setupRadioGroup(t *testing.T) (*RadioGroup, *styles.Theme) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Name:  "test-radiogroup",
		Type:  config.TypeRadioGroup,
		Label: "Select an option",
		Help:  "This is help text",
		Options: map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"id": "item1", "label": "Item 1"},
				map[string]interface{}{"id": "item2", "label": "Item 2"},
				map[string]interface{}{"id": "item3", "label": "Item 3"},
			},
		},
	}
	rg, err := NewRadioGroup(cfg, theme)
	require.NoError(t, err)
	return rg, theme
}

func TestNewRadioGroup(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name        string
		cfg         config.ComponentConfig
		expectError bool
		validate    func(*testing.T, *RadioGroup)
	}{
		{
			name: "valid radiogroup with default",
			cfg: config.ComponentConfig{
				Name:     "options",
				Type:     config.TypeRadioGroup,
				Label:    "Options",
				Required: true,
				Help:     "Select one",
				Options: map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{"id": "opt1", "label": "Option 1"},
						map[string]interface{}{"id": "opt2", "label": "Option 2"},
					},
				},
				Default: "opt2",
			},
			expectError: false,
			validate: func(t *testing.T, rg *RadioGroup) {
				assert.Equal(t, "options", rg.name)
				assert.Len(t, rg.items, 2)
				assert.Equal(t, 1, rg.selected)
				assert.Equal(t, 1, rg.initialValue)
			},
		},
		{
			name: "valid radiogroup without default",
			cfg: config.ComponentConfig{
				Name: "options",
				Type: config.TypeRadioGroup,
				Options: map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{"id": "opt1", "label": "Option 1"},
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, rg *RadioGroup) {
				assert.Equal(t, -1, rg.selected)
				assert.Equal(t, -1, rg.initialValue)
			},
		},
		{
			name: "invalid default value type",
			cfg: config.ComponentConfig{
				Name:    "options",
				Type:    config.TypeRadioGroup,
				Options: map[string]interface{}{"items": []interface{}{map[string]interface{}{"id": "opt1", "label": "Option 1"}}},
				Default: 123, // int instead of string
			},
			expectError: false,
			validate: func(t *testing.T, rg *RadioGroup) {
				assert.Equal(t, -1, rg.selected, "Invalid default type should be ignored")
			},
		},
		{
			name:        "invalid component type",
			cfg:         config.ComponentConfig{Name: "invalid", Type: config.TypeTextInput},
			expectError: true,
		},
		{
			name:        "no items in config",
			cfg:         config.ComponentConfig{Name: "no-items", Type: config.TypeRadioGroup},
			expectError: true,
		},
		{
			name:        "empty items array",
			cfg:         config.ComponentConfig{Name: "empty-items", Type: config.TypeRadioGroup, Options: map[string]interface{}{"items": []interface{}{}}},
			expectError: true,
		},
		{
			name: "invalid item format in array",
			cfg: config.ComponentConfig{
				Name: "invalid-item-format",
				Type: config.TypeRadioGroup,
				Options: map[string]interface{}{
					"items": []interface{}{"not-a-map"},
				},
			},
			// This does not error because of the type assertion guard, but it results in 0 items which is an error.
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rg, err := NewRadioGroup(tt.cfg, theme)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, rg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, rg)
				if tt.validate != nil {
					tt.validate(t, rg)
				}
			}
		})
	}
}

func TestRadioGroup_ErrorSimulation(t *testing.T) {
	t.Run("simulated type assertion error for id field", func(t *testing.T) {
		theme := styles.DefaultTheme()
		cfg := config.ComponentConfig{
			Name: "test-error",
			Type: config.TypeRadioGroup,
			Options: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"id": 123, "label": "Invalid ID Type"}, // ID should be string
				},
			},
		}

		rg, err := NewRadioGroup(cfg, theme)
		assert.Error(t, err)
		assert.Nil(t, rg)
		assert.Contains(t, err.Error(), "campo 'id' deve ser string")
	})

	t.Run("simulated type assertion error for label field", func(t *testing.T) {
		theme := styles.DefaultTheme()
		cfg := config.ComponentConfig{
			Name: "test-error",
			Type: config.TypeRadioGroup,
			Options: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"id": "valid-id", "label": 456}, // Label should be string
				},
			},
		}

		rg, err := NewRadioGroup(cfg, theme)
		assert.Error(t, err)
		assert.Nil(t, rg)
		assert.Contains(t, err.Error(), "campo 'label' deve ser string")
	})

	t.Run("simulated mixed valid and invalid items", func(t *testing.T) {
		theme := styles.DefaultTheme()
		cfg := config.ComponentConfig{
			Name: "test-error",
			Type: config.TypeRadioGroup,
			Options: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"id": "valid1", "label": "Valid 1"},
					map[string]interface{}{"id": 123, "label": "Invalid"}, // Invalid ID
					map[string]interface{}{"id": "valid2", "label": "Valid 2"},
				},
			},
		}

		rg, err := NewRadioGroup(cfg, theme)
		assert.Error(t, err)
		assert.Nil(t, rg)
		assert.Contains(t, err.Error(), "campo 'id' deve ser string")
	})

	t.Run("simulated error propagation in SetValue", func(t *testing.T) {
		rg, _ := setupRadioGroup(t)

		// Test error when setting invalid value type
		err := rg.SetValue(123) // Should be string
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "valor inválido: esperado string")

		// Test error when setting non-existent ID
		err = rg.SetValue("non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID não encontrado")
	})

	t.Run("simulated validation error", func(t *testing.T) {
		rg, _ := setupRadioGroup(t)

		// Set as required but no selection
		rg.required = true
		rg.selected = -1

		isValid := rg.IsValid()
		assert.False(t, isValid)
		assert.Equal(t, "Selecione uma opção", rg.GetError())
	})
}

func TestRadioGroup_Init(t *testing.T) {
	rg, _ := setupRadioGroup(t)
	assert.Nil(t, rg.Init())
}

func TestRadioGroup_Update(t *testing.T) {
	rg, _ := setupRadioGroup(t)

	// Define key messages
	upMsg := tea.KeyPressMsg{Code: tea.KeyUp}
	downMsg := tea.KeyPressMsg{Code: tea.KeyDown}
	kMsg := tea.KeyPressMsg{Code: 'k'}
	jMsg := tea.KeyPressMsg{Code: 'j'}
	enterMsg := tea.KeyPressMsg{Code: tea.KeyEnter}
	spaceMsg := tea.KeyPressMsg{Code: tea.KeySpace}
	otherMsg := tea.KeyPressMsg{Code: 'a'}

	t.Run("should not update when not focused", func(t *testing.T) {
		rg.cursor = 0
		rg.Update(downMsg)
		assert.Equal(t, 0, rg.cursor, "Cursor should not move when not focused")
	})

	t.Run("navigation", func(t *testing.T) {
		rg.SetFocus(true)
		rg.cursor = 1

		rg.Update(upMsg)
		assert.Equal(t, 0, rg.cursor)
		rg.Update(upMsg)
		assert.Equal(t, 0, rg.cursor, "Cursor should not go above 0")

		rg.Update(downMsg)
		assert.Equal(t, 1, rg.cursor)
		rg.Update(downMsg)
		assert.Equal(t, 2, rg.cursor)
		rg.Update(downMsg)
		assert.Equal(t, 2, rg.cursor, "Cursor should not go below max items")

		rg.cursor = 1 // reset for j/k
		rg.Update(kMsg)
		assert.Equal(t, 0, rg.cursor, "k should move up")
		rg.Update(jMsg)
		assert.Equal(t, 1, rg.cursor, "j should move down")
	})

	t.Run("selection", func(t *testing.T) {
		rg.SetFocus(true)
		rg.cursor = 1
		rg.selected = -1
		rg.SetError("an error")

		rg.Update(enterMsg)
		assert.Equal(t, 1, rg.selected, "Enter should select item at cursor")
		assert.Empty(t, rg.errorMsg, "Error should be cleared on selection")

		rg.cursor = 2
		rg.Update(spaceMsg)
		assert.Equal(t, 2, rg.selected, "Space should select item at cursor")

		rg.Update(otherMsg)
		assert.Equal(t, 2, rg.selected, "Other keys should not change selection")
	})
}

func TestRadioGroup_View(t *testing.T) {
	rg, _ := setupRadioGroup(t)

	t.Run("unselected state", func(t *testing.T) {
		rg.selected = -1
		view := rg.View()
		assert.Contains(t, view, "Select an option")
		assert.Contains(t, view, "( ) Item 1")
		assert.Contains(t, view, "( ) Item 2")
		assert.Contains(t, view, "( ) Item 3")
		assert.Contains(t, view, "This is help text")
	})

	t.Run("selected state", func(t *testing.T) {
		rg.selected = 1 // Select "Item 2"
		view := rg.View()
		assert.Contains(t, view, "( ) Item 1")
		assert.Contains(t, view, "(•) Item 2")
		assert.Contains(t, view, "( ) Item 3")
	})

	t.Run("focused state with cursor", func(t *testing.T) {
		rg.SetFocus(true)
		rg.cursor = 0 // Cursor on "Item 1"
		view := rg.View()
		// The theme renders focused items, let's just check for the content
		assert.Contains(t, view, "Item 1")
	})

	t.Run("error state should hide help text", func(t *testing.T) {
		rg.SetError("A selection is required")
		view := rg.View()
		assert.Contains(t, view, "✗ A selection is required")
		assert.NotContains(t, view, "This is help text")
	})
}

func TestRadioGroup_ComponentInterface(t *testing.T) {
	rg, _ := setupRadioGroup(t)

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, "test-radiogroup", rg.Name())
	})

	t.Run("CanFocus", func(t *testing.T) {
		assert.True(t, rg.CanFocus())
	})

	t.Run("SetFocus", func(t *testing.T) {
		rg.SetFocus(true)
		assert.True(t, rg.focused)
		rg.SetFocus(false)
		assert.False(t, rg.focused)
	})

	t.Run("GetError and SetError", func(t *testing.T) {
		rg.SetError("my error")
		assert.Equal(t, "my error", rg.GetError())
	})

	t.Run("Value and SetValue", func(t *testing.T) {
		assert.Equal(t, "", rg.Value(), "Initial value should be empty string if nothing is selected")

		err := rg.SetValue("item2")
		require.NoError(t, err)
		assert.Equal(t, "item2", rg.Value())
		assert.Equal(t, 1, rg.selected)
		assert.Equal(t, 1, rg.cursor, "Cursor should move to selected item on SetValue")

		err = rg.SetValue("non-existent-id")
		assert.Error(t, err)

		err = rg.SetValue(123)
		assert.Error(t, err)
	})
}

func TestRadioGroup_Validation(t *testing.T) {
	rg, _ := setupRadioGroup(t)

	t.Run("required field", func(t *testing.T) {
		rg.required = true
		rg.selected = -1
		assert.False(t, rg.IsValid())
		assert.Equal(t, "Selecione uma opção", rg.GetError())

		rg.selected = 0
		assert.True(t, rg.IsValid())
		assert.Empty(t, rg.GetError())
	})

	t.Run("optional field", func(t *testing.T) {
		rg.required = false
		rg.selected = -1
		assert.True(t, rg.IsValid())
		assert.Empty(t, rg.GetError())
	})
}

func TestRadioGroup_Reset(t *testing.T) {
	t.Run("reset with initial value", func(t *testing.T) {
		rg, _ := setupRadioGroup(t)
		err := rg.SetValue("item3") // Set initial selection
		require.NoError(t, err)
		rg.initialValue = rg.selected

		// Modify state
		rg.SetFocus(true)
		rg.selected = 0
		rg.cursor = 0
		rg.SetError("an error")

		rg.Reset()

		assert.Equal(t, 2, rg.selected, "Should reset to initial value")
		assert.Equal(t, 2, rg.cursor, "Cursor should reset to initial value")
		assert.False(t, rg.focused)
		assert.Empty(t, rg.GetError())
	})

	t.Run("reset without initial value", func(t *testing.T) {
		rg, _ := setupRadioGroup(t)
		rg.initialValue = -1 // Ensure no default

		// Modify state
		rg.SetFocus(true)
		rg.selected = 1
		rg.cursor = 1
		rg.SetError("an error")

		rg.Reset()

		assert.Equal(t, -1, rg.selected, "Should reset to no selection")
		assert.Equal(t, 0, rg.cursor, "Cursor should reset to 0")
		assert.False(t, rg.focused)
		assert.Empty(t, rg.GetError())
	})
}
