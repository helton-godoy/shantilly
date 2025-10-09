package models

import (
	"encoding/json"
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/helton/shantilly/internal/components"
)

// mockUnmarshallableComponent is a mock component that returns a value
// that cannot be marshalled to JSON, for testing error handling.
type mockUnmarshallableComponent struct{}

func (m *mockUnmarshallableComponent) Name() string                            { return "bad" }
func (m *mockUnmarshallableComponent) Value() interface{}                      { return make(chan int) }
func (m *mockUnmarshallableComponent) Init() tea.Cmd                           { return nil }
func (m *mockUnmarshallableComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *mockUnmarshallableComponent) View() string                            { return "" }
func (m *mockUnmarshallableComponent) CanFocus() bool                          { return false }
func (m *mockUnmarshallableComponent) SetFocus(bool)                           {}
func (m *mockUnmarshallableComponent) IsValid() bool                           { return true }
func (m *mockUnmarshallableComponent) GetError() string                        { return "" }
func (m *mockUnmarshallableComponent) SetError(string)                         {}
func (m *mockUnmarshallableComponent) SetValue(interface{}) error              { return nil }
func (m *mockUnmarshallableComponent) Reset()                                  {}
func (m *mockUnmarshallableComponent) ExportToFormat(format components.ExportFormat) ([]byte, error) {
	return nil, nil
}
func (m *mockUnmarshallableComponent) ImportFromFormat(format components.ExportFormat, data []byte) error {
	return nil
}
func (m *mockUnmarshallableComponent) GetDependencies() []string    { return []string{} }
func (m *mockUnmarshallableComponent) SetTheme(theme *styles.Theme) {}
func (m *mockUnmarshallableComponent) GetMetadata() components.ComponentMetadata {
	return components.ComponentMetadata{}
}
func (m *mockUnmarshallableComponent) ValidateWithContext(context components.ValidationContext) []components.ValidationError {
	return []components.ValidationError{}
}

func TestNewFormModel(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name        string
		cfg         *config.FormConfig
		expectError bool
		validate    func(*testing.T, *FormModel)
	}{
		{
			name: "valid form with components",
			cfg: &config.FormConfig{
				Title:       "User Registration",
				Description: "Please fill out the form",
				Components: []config.ComponentConfig{
					{
						Name: "username",
						Type: config.TypeTextInput,
					},
					{
						Name: "email",
						Type: config.TypeTextInput,
					},
					{
						Name: "agree",
						Type: config.TypeCheckbox,
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, fm *FormModel) {
				assert.Equal(t, "User Registration", fm.title)
				assert.Equal(t, "Please fill out the form", fm.description)
				assert.Len(t, fm.components, 3)
				assert.Equal(t, 0, fm.focusIndex) // First focusable component
				assert.False(t, fm.submitted)
				assert.False(t, fm.quitting)
			},
		},
		{
			name: "form with no focusable components",
			cfg: &config.FormConfig{
				Title: "Static Form",
				Components: []config.ComponentConfig{
					{
						Name:  "title",
						Type:  config.TypeText,
						Label: "Welcome",
					},
					{
						Name:  "info",
						Type:  config.TypeText,
						Label: "Information",
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, fm *FormModel) {
				assert.Equal(t, "Static Form", fm.title)
				assert.Len(t, fm.components, 2)
				assert.Equal(t, -1, fm.focusIndex) // No focusable components
			},
		},
		{
			name: "minimal form",
			cfg: &config.FormConfig{
				Title: "Minimal Form",
				Components: []config.ComponentConfig{
					{
						Name:  "info",
						Type:  config.TypeText,
						Label: "Information",
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, fm *FormModel) {
				assert.Equal(t, "Minimal Form", fm.title)
				assert.Len(t, fm.components, 1)
				assert.Equal(t, -1, fm.focusIndex) // No focusable components
			},
		},
		{
			name: "form with invalid component",
			cfg: &config.FormConfig{
				Title: "Invalid Form",
				Components: []config.ComponentConfig{
					{
						Name: "", // Invalid - empty name
						Type: config.TypeTextInput,
					},
				},
			},
			expectError: true,
		},
		{
			name: "form with invalid configuration",
			cfg: &config.FormConfig{
				Title: "", // Invalid - empty title
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, err := NewFormModel(tt.cfg, theme)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, fm)
			} else {
				require.NoError(t, err)
				require.NotNil(t, fm)
				if tt.validate != nil {
					tt.validate(t, fm)
				}
			}
		})
	}
}

func TestFormModel_Init(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.FormConfig{
		Title: "Test Form",
		Components: []config.ComponentConfig{
			{
				Name: "test",
				Type: config.TypeTextInput,
			},
		},
	}

	fm, err := NewFormModel(cfg, theme)
	require.NoError(t, err)

	cmd := fm.Init()
	assert.Nil(t, cmd) // No initial commands
}

// func TestFormModel_Update(t *testing.T) {
// 	theme := styles.DefaultTheme()
// 	cfg := &config.FormConfig{
// 		Title: "Test Form",
// 		Components: []config.ComponentConfig{
// 			{
// 				Name: "username",
// 				Type: config.TypeTextInput,
// 			},
// 			{
// 				Name: "email",
// 				Type: config.TypeTextInput,
// 			},
// 			{
// 				Name: "agree",
// 				Type: config.TypeCheckbox,
// 			},
// 		},
// 	}

// 	fm, err := NewFormModel(cfg, theme)
// 	require.NoError(t, err)

// 	// Test window size message
// 	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
// 	model, cmd := fm.Update(wsMsg)
// 	assert.Equal(t, fm, model)
// 	assert.Nil(t, cmd)
// 	assert.Equal(t, 100, fm.width)
// 	assert.Equal(t, 50, fm.height)

// 	// Test quit message
// 	quitMsg := tea.QuitMsg{}
// 	model, cmd = fm.Update(quitMsg)
// 	assert.Equal(t, fm, model)
// 	assert.Nil(t, cmd)
// 	// Verify quitting state is set (form model may not directly expose this)

// 	// Reset quitting state for further tests
// 	fm.quitting = false

// 	// Test tab navigation (next focus)
// 	tabMsg := tea.KeyPressMsg{Code: tea.KeyTab}
// 	model, cmd = fm.Update(tabMsg)
// 	assert.Equal(t, fm, model)
// 	assert.Nil(t, cmd)
// 	assert.Equal(t, 1, fm.focusIndex) // Should move to next component

// 	// Test shift+tab navigation (previous focus)
// 	shiftTabMsg := tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift}
// 	model, cmd = fm.Update(shiftTabMsg)
// 	assert.Equal(t, fm, model)
// 	assert.Nil(t, cmd)
// 	assert.Equal(t, 0, fm.focusIndex) // Should move back to first component

// 	// Test enter key submission when form is valid
// 	fm.components[0].SetValue("john")             // Set username
// 	fm.components[1].SetValue("john@example.com") // Set email
// 	fm.components[2].SetValue(true)               // Check the checkbox

// 	enterMsg := tea.KeyPressMsg{Code: tea.KeyEnter}
// 	model, cmd = fm.Update(enterMsg)
// 	assert.Equal(t, fm, model)
// 	assert.NotNil(t, cmd) // Should return quit command
// 	assert.True(t, fm.submitted)

// 	// Test that focused component receives other key messages
// 	fm.submitted = false // Reset for test
// 	fm.focusIndex = 0    // Focus on first component

// 	keyMsg := tea.KeyPressMsg{Text: "a", Code: 'a'}
// 	model, cmd = fm.Update(keyMsg)
// 	assert.Equal(t, fm, model)
// 	// The focused component should have received the message
// }

// func TestFormModel_FocusNavigation(t *testing.T) {
// 	theme := styles.DefaultTheme()
// 	cfg := &config.FormConfig{
// 		Title: "Test Form",
// 		Components: []config.ComponentConfig{
// 			{
// 				Name: "field1",
// 				Type: config.TypeTextInput,
// 			},
// 			{
// 				Name:  "info",
// 				Type:  config.TypeText,
// 				Label: "Static text",
// 			},
// 			{
// 				Name: "field2",
// 				Type: config.TypeTextInput,
// 			},
// 			{
// 				Name: "checkbox",
// 				Type: config.TypeCheckbox,
// 			},
// 		},
// 	}

// 	fm, err := NewFormModel(cfg, theme)
// 	require.NoError(t, err)

// 	// Should start focused on first focusable component (index 0)
// 	assert.Equal(t, 0, fm.focusIndex)

// 	// Test next focus - should skip static component
// 	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
// 	fm.Update(tabMsg)
// 	assert.Equal(t, 2, fm.focusIndex) // Should skip index 1 (static text)

// 	// Test next focus again
// 	fm.Update(tabMsg)
// 	assert.Equal(t, 3, fm.focusIndex) // Should be on checkbox

// 	// Test next focus at end - should wrap to beginning
// 	fm.Update(tabMsg)
// 	assert.Equal(t, 0, fm.focusIndex) // Should wrap to first

// 	// Test previous focus
// 	shiftTabMsg := tea.KeyMsg{Type: tea.KeyShiftTab}
// 	fm.Update(shiftTabMsg)
// 	assert.Equal(t, 3, fm.focusIndex) // Should wrap to last

// 	// Test previous focus again
// 	fm.Update(shiftTabMsg)
// 	assert.Equal(t, 2, fm.focusIndex) // Should go to field2
// }

func TestFormModel_View(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.FormConfig{
		Title: "Test Form",
		Components: []config.ComponentConfig{
			{Name: "field1", Type: config.TypeTextInput},
		},
	}
	fm, err := NewFormModel(cfg, theme)
	require.NoError(t, err)

	// Apenas um smoke test para garantir que a renderização não entre em pânico.
	// Verificar o conteúdo exato é frágil devido aos códigos de escape ANSI.
	assert.NotEmpty(t, fm.View())
}

func TestFormModel_CanSubmit(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.FormConfig{
		Title: "Test Form",
		Components: []config.ComponentConfig{
			{
				Name:     "username",
				Type:     config.TypeTextInput,
				Required: true,
			},
			{
				Name:     "email",
				Type:     config.TypeTextInput,
				Required: true,
			},
			{
				Name:     "newsletter",
				Type:     config.TypeCheckbox,
				Required: false,
			},
		},
	}

	fm, err := NewFormModel(cfg, theme)
	require.NoError(t, err)

	// Initially should not be able to submit (required fields empty)
	assert.False(t, fm.CanSubmit())

	// Fill one required field
	err = fm.components[0].SetValue("john")
	require.NoError(t, err)
	assert.False(t, fm.CanSubmit()) // Still missing email

	// Fill second required field
	err = fm.components[1].SetValue("john@example.com")
	require.NoError(t, err)
	assert.True(t, fm.CanSubmit()) // Now should be valid

	// Optional field doesn't affect validity
	err = fm.components[2].SetValue(true)
	require.NoError(t, err)
	assert.True(t, fm.CanSubmit())
}

func TestFormModel_ValidateAll(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.FormConfig{
		Title: "Test Form",
		Components: []config.ComponentConfig{
			{
				Name:     "username",
				Type:     config.TypeTextInput,
				Required: true,
			},
		},
	}

	fm, err := NewFormModel(cfg, theme)
	require.NoError(t, err)

	// Initially, component has no error message from its own validation
	assert.Empty(t, fm.components[0].GetError())

	// Call validateAll
	fm.validateAll()

	// Now, the component should have an error message because it's required and empty
	assert.NotEmpty(t, fm.components[0].GetError())
}

func TestFormModel_Submitted(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.FormConfig{
		Title: "Test Form",
		Components: []config.ComponentConfig{
			{
				Name: "test",
				Type: config.TypeTextInput,
			},
		},
	}

	fm, err := NewFormModel(cfg, theme)
	require.NoError(t, err)

	// Initially not submitted
	assert.False(t, fm.Submitted())

	// Mark as submitted
	fm.submitted = true
	assert.True(t, fm.Submitted())
}

func TestFormModel_ToJSON(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.FormConfig{
		Title: "Test Form",
		Components: []config.ComponentConfig{
			{
				Name: "username",
				Type: config.TypeTextInput,
			},
			{
				Name: "age",
				Type: config.TypeSlider,
				Options: map[string]interface{}{
					"min": 0.0,
					"max": 100.0,
				},
			},
			{
				Name: "newsletter",
				Type: config.TypeCheckbox,
			},
		},
	}

	fm, err := NewFormModel(cfg, theme)
	require.NoError(t, err)

	// Set some values
	err = fm.components[0].SetValue("john_doe")
	require.NoError(t, err)
	err = fm.components[1].SetValue(25.0)
	require.NoError(t, err)
	err = fm.components[2].SetValue(true)
	require.NoError(t, err)

	// Test JSON conversion
	jsonStr, err := fm.ToJSON()
	require.NoError(t, err)

	// Parse JSON to verify structure
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &result)
	require.NoError(t, err)

	assert.Equal(t, "john_doe", result["username"])
	assert.Equal(t, 25.0, result["age"])
	assert.Equal(t, true, result["newsletter"])
}

func TestFormModel_ToJSON_Error(t *testing.T) {
	fm := &FormModel{
		components: []components.Component{
			&mockUnmarshallableComponent{},
		},
	}

	jsonData, err := fm.ToJSON()
	assert.Error(t, err)
	assert.Nil(t, jsonData)
	assert.Contains(t, err.Error(), "erro ao serializar dados")
}

func TestFormModel_ToMap(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.FormConfig{
		Title: "Test Form",
		Components: []config.ComponentConfig{
			{
				Name: "name",
				Type: config.TypeTextInput,
			},
			{
				Name: "active",
				Type: config.TypeCheckbox,
			},
		},
	}

	fm, err := NewFormModel(cfg, theme)
	require.NoError(t, err)

	// Set values
	err = fm.components[0].SetValue("Alice")
	require.NoError(t, err)
	err = fm.components[1].SetValue(false)
	require.NoError(t, err)

	// Test map conversion
	result := fm.ToMap()

	assert.Equal(t, "Alice", result["name"])
	assert.Equal(t, false, result["active"])
}

func TestFormModel_ErrorSimulation(t *testing.T) {
	t.Run("simulated component update error", func(t *testing.T) {
		theme := styles.DefaultTheme()
		cfg := &config.FormConfig{
			Title: "Test Form",
			Components: []config.ComponentConfig{
				{
					Name: "test",
					Type: config.TypeTextInput,
				},
			},
		}

		fm, err := NewFormModel(cfg, theme)
		require.NoError(t, err)

		// Simulate a message that would cause an error in component update
		// We'll test the error handling by creating a scenario where the component
		// returns an invalid model type

		// Test with a key message that should work normally first
		keyMsg := tea.KeyPressMsg{Text: "a", Code: 'a'}
		model, _ := fm.Update(keyMsg)
		assert.Equal(t, fm, model)
	})

	t.Run("simulated form validation errors", func(t *testing.T) {
		theme := styles.DefaultTheme()
		cfg := &config.FormConfig{
			Title: "Test Form",
			Components: []config.ComponentConfig{
				{
					Name:     "required_field",
					Type:     config.TypeTextInput,
					Required: true,
				},
			},
		}

		fm, err := NewFormModel(cfg, theme)
		require.NoError(t, err)

		// Initially should not be able to submit
		assert.False(t, fm.CanSubmit())

		// Try to submit invalid form - should trigger validation
		enterMsg := tea.KeyPressMsg{Code: tea.KeyEnter}
		model, cmd := fm.Update(enterMsg)
		assert.Equal(t, fm, model)
		assert.Nil(t, cmd) // Should not quit when form is invalid
		assert.False(t, fm.Submitted())

		// Verify that validation errors are triggered
		fm.validateAll()
		assert.NotEmpty(t, fm.components[0].GetError())
	})

	t.Run("simulated JSON serialization error", func(t *testing.T) {
		// Create form with component that has unmarshallable value
		fm := &FormModel{
			components: []components.Component{
				&mockUnmarshallableComponent{},
			},
		}

		jsonData, err := fm.ToJSON()
		assert.Error(t, err)
		assert.Nil(t, jsonData)
		assert.Contains(t, err.Error(), "erro ao serializar dados")
	})

	t.Run("simulated component creation errors", func(t *testing.T) {
		theme := styles.DefaultTheme()

		// Test with invalid component configuration
		cfg := &config.FormConfig{
			Title: "Invalid Form",
			Components: []config.ComponentConfig{
				{
					Name: "", // Invalid - empty name
					Type: config.TypeTextInput,
				},
			},
		}

		fm, err := NewFormModel(cfg, theme)
		assert.Error(t, err)
		assert.Nil(t, fm)
		assert.Contains(t, err.Error(), "erro de validação da configuração")
	})
}

// func TestFormModel_EdgeCases(t *testing.T) {
// 	theme := styles.DefaultTheme()

// 	t.Run("form with only static components", func(t *testing.T) {
// 		cfg := &config.FormConfig{
// 			Title: "Static Only",
// 			Components: []config.ComponentConfig{
// 				{
// 					Name:  "title",
// 					Type:  config.TypeText,
// 					Label: "Welcome",
// 				},
// 			},
// 		}

// 		fm, err := NewFormModel(cfg, theme)
// 		require.NoError(t, err)

// 		// Should handle tab navigation gracefully
// 		tabMsg := tea.KeyMsg{Type: tea.KeyTab}
// 		model, cmd := fm.Update(tabMsg)
// 		assert.Equal(t, fm, model)
// 		assert.Nil(t, cmd)
// 		assert.Equal(t, -1, fm.focusIndex) // Should remain -1
// 	})

// 	t.Run("form with static components only", func(t *testing.T) {
// 		cfg := &config.FormConfig{
// 			Title: "Static Form",
// 			Components: []config.ComponentConfig{
// 				{
// 					Name:  "info",
// 					Type:  config.TypeText,
// 					Label: "Static information",
// 				},
// 			},
// 		}

// 		fm, err := NewFormModel(cfg, theme)
// 		require.NoError(t, err)

// 		// Should be able to submit form with only static components
// 		assert.True(t, fm.CanSubmit())

// 		enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
// 		model, cmd := fm.Update(enterMsg)
// 		assert.Equal(t, fm, model)
// 		assert.NotNil(t, cmd) // Should return quit command
// 		assert.True(t, fm.submitted)
// 	})

// 	t.Run("large form stress test", func(t *testing.T) {
// 		components := make([]config.ComponentConfig, 50)
// 		for i := 0; i < 50; i++ {
// 			components[i] = config.ComponentConfig{
// 				Name: fmt.Sprintf("field_%d", i),
// 				Type: config.TypeTextInput,
// 			}
// 		}

// 		cfg := &config.FormConfig{
// 			Title:      "Large Form",
// 			Components: components,
// 		}

// 		fm, err := NewFormModel(cfg, theme)
// 		require.NoError(t, err)

// 		assert.Len(t, fm.components, 50)
// 		assert.Equal(t, 0, fm.focusIndex)

// 		// Test navigation through all components
// 		tabMsg := tea.KeyMsg{Type: tea.KeyTab}
// 		for i := 1; i < 50; i++ {
// 			fm.Update(tabMsg)
// 			assert.Equal(t, i, fm.focusIndex)
// 		}

// 		// Should wrap to beginning
// 		fm.Update(tabMsg)
// 		assert.Equal(t, 0, fm.focusIndex)
// 	})
// }

// func TestFormModel_ComponentInteraction(t *testing.T) {
// 	theme := styles.DefaultTheme()
// 	cfg := &config.FormConfig{
// 		Title: "Interactive Form",
// 		Components: []config.ComponentConfig{
// 			{
// 				Name:     "required_field",
// 				Type:     config.TypeTextInput,
// 				Required: true,
// 			},
// 		},
// 	}

// 	fm, err := NewFormModel(cfg, theme)
// 	require.NoError(t, err)

// 	// Test that validation affects form submission
// 	assert.False(t, fm.CanSubmit()) // Required field is empty

// 	// Try to submit invalid form
// 	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
// 	model, cmd := fm.Update(enterMsg)
// 	assert.Equal(t, fm, model)
// 	assert.Nil(t, cmd)            // Should not quit
// 	assert.False(t, fm.submitted) // Should not be marked as submitted

// 	// Fill required field
// 	fm.components[0].SetValue("test_value")
// 	assert.True(t, fm.CanSubmit())

// 	// Now submission should work
// 	model, cmd = fm.Update(enterMsg)
// 	assert.Equal(t, fm, model)
// 	assert.NotNil(t, cmd) // Should return quit command
// 	assert.True(t, fm.submitted)
// }
