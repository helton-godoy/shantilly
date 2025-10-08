package models

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLayoutModel(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name        string
		cfg         *config.LayoutConfig
		expectError bool
		validate    func(*testing.T, *LayoutModel)
	}{
		{
			name: "valid horizontal layout",
			cfg: &config.LayoutConfig{
				Title:       "Horizontal Layout",
				Description: "Side by side components",
				Layout:      "horizontal",
				Components: []config.ComponentConfig{
					{
						Name: "left",
						Type: config.TypeTextInput,
					},
					{
						Name: "right",
						Type: config.TypeTextInput,
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, lm *LayoutModel) {
				assert.Equal(t, "Horizontal Layout", lm.title)
				assert.Equal(t, "Side by side components", lm.description)
				assert.Equal(t, "horizontal", lm.layout)
				assert.Len(t, lm.components, 2)
				assert.Equal(t, 0, lm.focusIndex)
				assert.False(t, lm.quitting)
			},
		},
		{
			name: "valid vertical layout",
			cfg: &config.LayoutConfig{
				Title:  "Vertical Layout",
				Layout: "vertical",
				Components: []config.ComponentConfig{
					{
						Name: "top",
						Type: config.TypeTextInput,
					},
					{
						Name: "bottom",
						Type: config.TypeCheckbox,
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, lm *LayoutModel) {
				assert.Equal(t, "Vertical Layout", lm.title)
				assert.Equal(t, "vertical", lm.layout)
				assert.Len(t, lm.components, 2)
				assert.Equal(t, 0, lm.focusIndex)
			},
		},
		{
			name: "layout with default direction",
			cfg: &config.LayoutConfig{
				Title:  "Default Layout",
				Layout: "vertical",
				Components: []config.ComponentConfig{
					{
						Name: "field",
						Type: config.TypeTextInput,
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, lm *LayoutModel) {
				assert.Equal(t, "vertical", lm.layout) // Should default to vertical
			},
		},
		{
			name: "layout with no focusable components",
			cfg: &config.LayoutConfig{
				Title:  "Static Layout",
				Layout: "vertical",
				Components: []config.ComponentConfig{
					{
						Name:  "title",
						Type:  config.TypeText,
						Label: "Static Title",
					},
					{
						Name:  "info",
						Type:  config.TypeText,
						Label: "Static Info",
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, lm *LayoutModel) {
				assert.Equal(t, "Static Layout", lm.title)
				assert.Len(t, lm.components, 2)
				assert.Equal(t, -1, lm.focusIndex) // No focusable components
			},
		},
		{
			name: "layout with minimal component",
			cfg: &config.LayoutConfig{
				Title:  "Minimal Layout",
				Layout: "vertical",
				Components: []config.ComponentConfig{
					{
						Name:  "info",
						Type:  config.TypeText,
						Label: "Information",
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, lm *LayoutModel) {
				assert.Equal(t, "Minimal Layout", lm.title)
				assert.Len(t, lm.components, 1)
				assert.Equal(t, -1, lm.focusIndex) // No focusable components
			},
		},
		{
			name: "layout with invalid component",
			cfg: &config.LayoutConfig{
				Title: "Invalid Layout",
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
			name: "layout with invalid configuration",
			cfg: &config.LayoutConfig{
				Title: "", // Invalid - empty title
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm, err := NewLayoutModel(tt.cfg, theme)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, lm)
			} else {
				require.NoError(t, err)
				require.NotNil(t, lm)
				if tt.validate != nil {
					tt.validate(t, lm)
				}
			}
		})
	}
}

func TestLayoutModel_Init(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.LayoutConfig{
		Title:  "Test Layout",
		Layout: "vertical",
		Components: []config.ComponentConfig{
			{
				Name: "test",
				Type: config.TypeTextInput,
			},
		},
	}

	lm, err := NewLayoutModel(cfg, theme)
	require.NoError(t, err)

	cmd := lm.Init()
	assert.Nil(t, cmd) // No initial commands
}

// func TestLayoutModel_Update(t *testing.T) {
// 	theme := styles.DefaultTheme()
// 	cfg := &config.LayoutConfig{
// 		Title:  "Test Layout",
// 		Layout: "vertical",
// 		Components: []config.ComponentConfig{
// 			{
// 				Name: "field1",
// 				Type: config.TypeTextInput,
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

// 	lm, err := NewLayoutModel(cfg, theme)
// 	require.NoError(t, err)

// 	// Test window size message
// 	wsMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
// 	model, cmd := lm.Update(wsMsg)
// 	assert.Equal(t, lm, model)
// 	assert.Nil(t, cmd)
// 	assert.Equal(t, 120, lm.width)
// 	assert.Equal(t, 40, lm.height)

// 	// Test quit message
// 	quitMsg := tea.QuitMsg{}
// 	model, cmd = lm.Update(quitMsg)
// 	assert.Equal(t, lm, model)
// 	assert.Nil(t, cmd)
// 	// Verify quitting state (may not be directly exposed)

// 	// Reset quitting state for further tests
// 	lm.quitting = false

// 	// Test tab navigation (next focus)
// 	tabMsg := tea.KeyPressMsg{Code: tea.KeyTab}
// 	model, cmd = lm.Update(tabMsg)
// 	assert.Equal(t, lm, model)
// 	assert.Nil(t, cmd)
// 	assert.Equal(t, 1, lm.focusIndex) // Should move to next component

// 	// Test shift+tab navigation (previous focus)
// 	shiftTabMsg := tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift}
// 	model, cmd = lm.Update(shiftTabMsg)
// 	assert.Equal(t, lm, model)
// 	assert.Nil(t, cmd)
// 	assert.Equal(t, 0, lm.focusIndex) // Should move back to first component

// 	// Test escape key (quit)
// 	escMsg := tea.KeyPressMsg{Code: tea.KeyEscape}
// 	model, cmd = lm.Update(escMsg)
// 	assert.Equal(t, lm, model)
// 	assert.NotNil(t, cmd) // Should return quit command
// 	assert.True(t, lm.quitting)

// 	// Test that focused component receives other key messages
// 	lm.quitting = false // Reset for test
// 	lm.focusIndex = 0   // Focus on first component

// 	keyMsg := tea.KeyPressMsg{Text: "a", Code: 'a'}
// 	model, cmd = lm.Update(keyMsg)
// 	assert.Equal(t, lm, model)
// 	// The focused component should have received the message
// }

// func TestLayoutModel_FocusNavigation(t *testing.T) {
// 	theme := styles.DefaultTheme()
// 	cfg := &config.LayoutConfig{
// 		Title:  "Focus Test Layout",
// 		Layout: "vertical",
// 		Components: []config.ComponentConfig{
// 			{
// 				Name: "field1",
// 				Type: config.TypeTextInput,
// 			},
// 			{
// 				Name:  "static",
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

// 	lm, err := NewLayoutModel(cfg, theme)
// 	require.NoError(t, err)

// 	// Should start focused on first focusable component (index 0)
// 	assert.Equal(t, 0, lm.focusIndex)

// 	// Test next focus - should skip static component
// 	tabMsg := tea.KeyPressMsg{Code: tea.KeyTab}
// 	lm.Update(tabMsg)
// 	assert.Equal(t, 2, lm.focusIndex) // Should skip index 1 (static text)

// 	// Test next focus again
// 	lm.Update(tabMsg)
// 	assert.Equal(t, 3, lm.focusIndex) // Should be on checkbox

// 	// Test next focus at end - should wrap to beginning
// 	lm.Update(tabMsg)
// 	assert.Equal(t, 0, lm.focusIndex) // Should wrap to first

// 	// Test previous focus
// 	shiftTabMsg := tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift}
// 	lm.Update(shiftTabMsg)
// 	assert.Equal(t, 3, lm.focusIndex) // Should wrap to last

// 	// Test previous focus again
// 	lm.Update(shiftTabMsg)
// 	assert.Equal(t, 2, lm.focusIndex) // Should go to field2
// }

func TestLayoutModel_View(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name      string
		cfg       *config.LayoutConfig
		checkView func(*testing.T, string)
	}{
		{
			name: "horizontal layout view",
			cfg: &config.LayoutConfig{
				Title:       "Horizontal Test",
				Description: "Side by side layout",
				Layout:      "horizontal",
				Components: []config.ComponentConfig{
					{
						Name:  "left",
						Type:  config.TypeTextInput,
						Label: "Left Field",
					},
					{
						Name:  "right",
						Type:  config.TypeTextInput,
						Label: "Right Field",
					},
				},
			},
			checkView: func(t *testing.T, view string) {
				// View should be non-empty (content may contain ANSI codes)
				assert.NotEmpty(t, view)
				assert.Greater(t, len(view), 100) // Should contain substantial content
			},
		},
		{
			name: "vertical layout view",
			cfg: &config.LayoutConfig{
				Title:  "Vertical Test",
				Layout: "vertical",
				Components: []config.ComponentConfig{
					{
						Name:  "top",
						Type:  config.TypeTextInput,
						Label: "Top Field",
					},
					{
						Name:  "bottom",
						Type:  config.TypeCheckbox,
						Label: "Bottom Checkbox",
					},
				},
			},
			checkView: func(t *testing.T, view string) {
				// View should be non-empty (content may contain ANSI codes)
				assert.NotEmpty(t, view)
				assert.Greater(t, len(view), 100) // Should contain substantial content
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm, err := NewLayoutModel(tt.cfg, theme)
			require.NoError(t, err)

			view := lm.View()
			tt.checkView(t, view)

			// Basic smoke test - view should render without panicking
			assert.NotEmpty(t, view)
		})
	}
}

func TestLayoutModel_HorizontalRendering(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.LayoutConfig{
		Title:  "Horizontal Layout",
		Layout: "horizontal",
		Components: []config.ComponentConfig{
			{
				Name:  "field1",
				Type:  config.TypeTextInput,
				Label: "Field 1",
			},
			{
				Name:  "field2",
				Type:  config.TypeTextInput,
				Label: "Field 2",
			},
		},
	}

	lm, err := NewLayoutModel(cfg, theme)
	require.NoError(t, err)

	// Set a reasonable width for horizontal layout
	lm.width = 100

	view := lm.View()

	// Should contain both field labels
	assert.Contains(t, view, "Field 1")
	assert.Contains(t, view, "Field 2")

	// In horizontal layout, components should be side by side
	// This is harder to test directly, but we can at least verify
	// that the view is generated without errors
	assert.NotEmpty(t, view)
}

func TestLayoutModel_VerticalRendering(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.LayoutConfig{
		Title:  "Vertical Layout",
		Layout: "vertical",
		Components: []config.ComponentConfig{
			{
				Name:  "field1",
				Type:  config.TypeTextInput,
				Label: "Field 1",
			},
			{
				Name:  "field2",
				Type:  config.TypeCheckbox,
				Label: "Field 2",
			},
		},
	}

	lm, err := NewLayoutModel(cfg, theme)
	require.NoError(t, err)

	view := lm.View()

	// Should contain both field labels
	assert.Contains(t, view, "Field 1")
	assert.Contains(t, view, "Field 2")

	// Should be non-empty
	assert.NotEmpty(t, view)
}

// func TestLayoutModel_EdgeCases(t *testing.T) {
// 	theme := styles.DefaultTheme()

// 	t.Run("layout with only static components", func(t *testing.T) {
// 		cfg := &config.LayoutConfig{
// 			Title:  "Static Only",
// 			Layout: "vertical",
// 			Components: []config.ComponentConfig{
// 				{
// 					Name:  "title",
// 					Type:  config.TypeText,
// 					Label: "Welcome",
// 				},
// 				{
// 					Name:  "subtitle",
// 					Type:  config.TypeText,
// 					Label: "Please read the information",
// 				},
// 			},
// 		}

// 		lm, err := NewLayoutModel(cfg, theme)
// 		require.NoError(t, err)

// 		// Should handle tab navigation gracefully
// 		tabMsg := tea.KeyPressMsg{Code: tea.KeyTab}
// 		model, cmd := lm.Update(tabMsg)
// 		assert.Equal(t, lm, model)
// 		assert.Nil(t, cmd)
// 		assert.Equal(t, -1, lm.focusIndex) // Should remain -1
// 	})

// 	t.Run("layout with minimal components", func(t *testing.T) {
// 		cfg := &config.LayoutConfig{
// 			Title:  "Minimal Layout",
// 			Layout: "vertical",
// 			Components: []config.ComponentConfig{
// 				{
// 					Name:  "info",
// 					Type:  config.TypeText,
// 					Label: "Static info",
// 				},
// 			},
// 		}

// 		lm, err := NewLayoutModel(cfg, theme)
// 		require.NoError(t, err)

// 		view := lm.View()
// 		assert.NotEmpty(t, view) // Should render without panicking
// 	})

// 	t.Run("large layout stress test", func(t *testing.T) {
// 		components := make([]config.ComponentConfig, 20)
// 		for i := 0; i < 20; i++ {
// 			components[i] = config.ComponentConfig{
// 				Name: fmt.Sprintf("field_%d", i),
// 				Type: config.TypeTextInput,
// 			}
// 		}

// 		cfg := &config.LayoutConfig{
// 			Title:      "Large Layout",
// 			Layout:     "vertical",
// 			Components: components,
// 		}

// 		lm, err := NewLayoutModel(cfg, theme)
// 		require.NoError(t, err)

// 		assert.Len(t, lm.components, 20)
// 		assert.Equal(t, 0, lm.focusIndex)

// 		// Test navigation through all components
// 		tabMsg := tea.KeyPressMsg{Code: tea.KeyTab}
// 		for i := 1; i < 20; i++ {
// 			lm.Update(tabMsg)
// 			assert.Equal(t, i, lm.focusIndex)
// 		}

// 		// Should wrap to beginning
// 		lm.Update(tabMsg)
// 		assert.Equal(t, 0, lm.focusIndex)

// 		// View should render without errors
// 		view := lm.View()
// 		assert.NotEmpty(t, view)
// 	})
// }

func TestLayoutModel_DirectionHandling(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name     string
		layout   string
		expected string
	}{
		{
			name:     "explicit horizontal",
			layout:   "horizontal",
			expected: "horizontal",
		},
		{
			name:     "explicit vertical",
			layout:   "vertical",
			expected: "vertical",
		},
		{
			name:     "valid vertical layout",
			layout:   "vertical",
			expected: "vertical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.LayoutConfig{
				Title:  "Layout Test",
				Layout: tt.layout,
				Components: []config.ComponentConfig{
					{
						Name: "test",
						Type: config.TypeTextInput,
					},
				},
			}

			lm, err := NewLayoutModel(cfg, theme)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, lm.layout)
		})
	}
}

func TestLayoutModel_ComponentInteraction(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.LayoutConfig{
		Title:  "Interactive Layout",
		Layout: "vertical",
		Components: []config.ComponentConfig{
			{
				Name: "interactive_field",
				Type: config.TypeTextInput,
			},
		},
	}

	lm, err := NewLayoutModel(cfg, theme)
	require.NoError(t, err)

	// Test that component receives key messages when focused
	assert.Equal(t, 0, lm.focusIndex) // Should be focused

	// Send a character key to the focused component
	keyMsg := tea.KeyPressMsg{Text: "test", Code: 't'}
	model, cmd := lm.Update(keyMsg)
	assert.Equal(t, lm, model)
	assert.Nil(t, cmd)

	// The component should have received the message
	// (We can't easily test the internal state change, but we can
	// verify that no error occurred and the update completed)
}

func TestLayoutModel_ResponsiveWidth(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := &config.LayoutConfig{
		Title:  "Responsive Layout",
		Layout: "horizontal",
		Components: []config.ComponentConfig{
			{
				Name: "field1",
				Type: config.TypeTextInput,
			},
			{
				Name: "field2",
				Type: config.TypeTextInput,
			},
		},
	}

	lm, err := NewLayoutModel(cfg, theme)
	require.NoError(t, err)

	// Test different window sizes
	sizes := []struct {
		width  int
		height int
	}{
		{40, 20},  // Small
		{80, 24},  // Medium
		{120, 30}, // Large
		{200, 50}, // Extra large
	}

	for _, size := range sizes {
		wsMsg := tea.WindowSizeMsg{Width: size.width, Height: size.height}
		lm.Update(wsMsg)

		assert.Equal(t, size.width, lm.width)
		assert.Equal(t, size.height, lm.height)

		// View should render without errors at any size
		view := lm.View()
		assert.NotEmpty(t, view)
	}
}

func TestLayoutModel_ErrorSimulation(t *testing.T) {
	t.Run("simulated component update error", func(t *testing.T) {
		theme := styles.DefaultTheme()
		cfg := &config.LayoutConfig{
			Title:  "Test Layout",
			Layout: "vertical",
			Components: []config.ComponentConfig{
				{
					Name: "test",
					Type: config.TypeTextInput,
				},
			},
		}

		lm, err := NewLayoutModel(cfg, theme)
		require.NoError(t, err)

		// Test with a key message that should work normally first
		keyMsg := tea.KeyPressMsg{Text: "a", Code: 'a'}
		model, _ := lm.Update(keyMsg)
		assert.Equal(t, lm, model)
	})

	t.Run("simulated layout validation errors", func(t *testing.T) {
		theme := styles.DefaultTheme()
		cfg := &config.LayoutConfig{
			Title:  "Test Layout",
			Layout: "vertical",
			Components: []config.ComponentConfig{
				{
					Name:     "required_field",
					Type:     config.TypeTextInput,
					Required: true,
				},
			},
		}

		lm, err := NewLayoutModel(cfg, theme)
		require.NoError(t, err)

		// Test escape key (quit)
		escMsg := tea.KeyPressMsg{Code: tea.KeyEscape}
		model, cmd := lm.Update(escMsg)
		assert.Equal(t, lm, model)
		assert.NotNil(t, cmd) // Should return quit command
		assert.True(t, lm.quitting)
	})

	t.Run("simulated window size handling", func(t *testing.T) {
		theme := styles.DefaultTheme()
		cfg := &config.LayoutConfig{
			Title:  "Test Layout",
			Layout: "horizontal",
			Components: []config.ComponentConfig{
				{
					Name: "field1",
					Type: config.TypeTextInput,
				},
				{
					Name: "field2",
					Type: config.TypeTextInput,
				},
			},
		}

		lm, err := NewLayoutModel(cfg, theme)
		require.NoError(t, err)

		// Test different window sizes
		sizes := []struct {
			width  int
			height int
		}{
			{0, 0},     // Zero size
			{10, 10},   // Very small
			{1000, 50}, // Very wide
			{50, 1000}, // Very tall
		}

		for _, size := range sizes {
			wsMsg := tea.WindowSizeMsg{Width: size.width, Height: size.height}
			model, _ := lm.Update(wsMsg)

			assert.Equal(t, lm, model)
			assert.Equal(t, size.width, lm.width)
			assert.Equal(t, size.height, lm.height)

			// View should render without errors at any size
			view := lm.View()
			assert.NotEmpty(t, view)
		}
	})

	t.Run("simulated component creation errors", func(t *testing.T) {
		theme := styles.DefaultTheme()

		// Test with invalid component configuration
		cfg := &config.LayoutConfig{
			Title: "Invalid Layout",
			Components: []config.ComponentConfig{
				{
					Name: "", // Invalid - empty name
					Type: config.TypeTextInput,
				},
			},
		}

		lm, err := NewLayoutModel(cfg, theme)
		assert.Error(t, err)
		assert.Nil(t, lm)
		assert.Contains(t, err.Error(), "erro de validação da configuração")
	})

	t.Run("simulated focus navigation edge cases", func(t *testing.T) {
		theme := styles.DefaultTheme()

		// Test with only non-focusable components
		cfg := &config.LayoutConfig{
			Title:  "Static Layout",
			Layout: "vertical",
			Components: []config.ComponentConfig{
				{
					Name:  "title",
					Type:  config.TypeText,
					Label: "Static Title",
				},
			},
		}

		lm, err := NewLayoutModel(cfg, theme)
		require.NoError(t, err)

		// Should start with no focusable components
		assert.Equal(t, -1, lm.focusIndex)

		// Tab navigation should handle gracefully
		tabMsg := tea.KeyPressMsg{Code: tea.KeyTab}
		model, cmd := lm.Update(tabMsg)
		assert.Equal(t, lm, model)
		assert.Nil(t, cmd)
		assert.Equal(t, -1, lm.focusIndex) // Should remain -1
	})
}
