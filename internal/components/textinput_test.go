package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTextInput(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name    string
		config  config.ComponentConfig
		wantErr bool
	}{
		{
			name: "valid textinput",
			config: config.ComponentConfig{
				Type:        config.TypeTextInput,
				Name:        "username",
				Label:       "Username",
				Placeholder: "Enter username",
				Required:    true,
			},
			wantErr: false,
		},
		{
			name: "wrong type",
			config: config.ComponentConfig{
				Type: config.TypeCheckbox,
				Name: "test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti, err := NewTextInput(tt.config, theme)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, ti)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, ti)
				assert.Equal(t, tt.config.Name, ti.Name())
				assert.True(t, ti.CanFocus())
			}
		})
	}
}

func TestTextInput_Validation(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name       string
		config     config.ComponentConfig
		setValue   string
		wantValid  bool
		wantErrMsg string
	}{
		{
			name: "required field empty",
			config: config.ComponentConfig{
				Type:     config.TypeTextInput,
				Name:     "required",
				Required: true,
			},
			setValue:   "",
			wantValid:  false,
			wantErrMsg: "obrigatório",
		},
		{
			name: "required field filled",
			config: config.ComponentConfig{
				Type:     config.TypeTextInput,
				Name:     "required",
				Required: true,
			},
			setValue:  "test",
			wantValid: true,
		},
		{
			name: "min length validation",
			config: config.ComponentConfig{
				Type: config.TypeTextInput,
				Name: "minlen",
				Options: map[string]interface{}{
					"min_length": 5,
				},
			},
			setValue:   "abc",
			wantValid:  false,
			wantErrMsg: "Mínimo",
		},
		{
			name: "max length validation",
			config: config.ComponentConfig{
				Type: config.TypeTextInput,
				Name: "maxlen",
				Options: map[string]interface{}{
					"max_length": 5,
				},
			},
			setValue:   "abcdefghij",
			wantValid:  false,
			wantErrMsg: "Máximo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti, err := NewTextInput(tt.config, theme)
			require.NoError(t, err)

			err = ti.SetValue(tt.setValue)
			require.NoError(t, err)

			valid := ti.IsValid()
			assert.Equal(t, tt.wantValid, valid)

			if !tt.wantValid {
				assert.Contains(t, ti.GetError(), tt.wantErrMsg)
			}
		})
	}
}

func TestTextInput_Value(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Type: config.TypeTextInput,
		Name: "test",
	}

	ti, err := NewTextInput(cfg, theme)
	require.NoError(t, err)

	// Test SetValue and Value
	err = ti.SetValue("hello")
	require.NoError(t, err)
	assert.Equal(t, "hello", ti.Value())

	// Test invalid value type
	err = ti.SetValue(123)
	assert.Error(t, err)
}

func TestTextInput_Reset(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Type:    config.TypeTextInput,
		Name:    "test",
		Default: "initial",
	}

	ti, err := NewTextInput(cfg, theme)
	require.NoError(t, err)

	// Change value
	err = ti.SetValue("changed")
	require.NoError(t, err)
	assert.Equal(t, "changed", ti.Value())

	// Reset
	ti.Reset()
	assert.Equal(t, "initial", ti.Value())
	assert.Empty(t, ti.GetError())
}

func TestTextInput_Init(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Type: config.TypeTextInput,
		Name: "test",
	}

	ti, err := NewTextInput(cfg, theme)
	require.NoError(t, err)

	cmd := ti.Init()
	assert.Nil(t, cmd)
}

func TestTextInput_Update(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Type: config.TypeTextInput,
		Name: "test",
	}

	ti, err := NewTextInput(cfg, theme)
	require.NoError(t, err)

	// Test window size message
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	model, cmd := ti.Update(wsMsg)
	assert.Equal(t, ti, model)
	assert.Nil(t, cmd)

	// Test key message when not focused
	keyMsg := tea.KeyPressMsg{Text: "a", Code: 'a'}
	model, cmd = ti.Update(keyMsg)
	assert.Equal(t, ti, model)
	assert.Nil(t, cmd)

	// Test key message when focused
	ti.SetFocus(true)
	model, cmd = ti.Update(keyMsg)
	assert.Equal(t, ti, model)
	// Error should be cleared when typing
	assert.Empty(t, ti.GetError())
}

func TestTextInput_View(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Type:  config.TypeTextInput,
		Name:  "test",
		Label: "Test Label",
	}

	ti, err := NewTextInput(cfg, theme)
	require.NoError(t, err)

	view := ti.View()
	assert.NotEmpty(t, view)
}

func TestTextInput_SetFocus(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Type: config.TypeTextInput,
		Name: "test",
	}

	ti, err := NewTextInput(cfg, theme)
	require.NoError(t, err)

	// Initially not focused
	ti.SetFocus(false)
	// Set focus to true
	ti.SetFocus(true)
	// Set focus to false again
	ti.SetFocus(false)
}

func TestTextInput_SetError(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Type: config.TypeTextInput,
		Name: "test",
	}

	ti, err := NewTextInput(cfg, theme)
	require.NoError(t, err)

	// Set error
	ti.SetError("Test error")
	assert.Equal(t, "Test error", ti.GetError())

	// Clear error
	ti.SetError("")
	assert.Equal(t, "", ti.GetError())
}

func TestTextInput_JoinVertical(t *testing.T) {
	// Test the joinVertical utility function
	result := joinVertical("line1", "line2", "line3")
	assert.Contains(t, result, "line1")
	assert.Contains(t, result, "line2")
	assert.Contains(t, result, "line3")
	assert.NotEmpty(t, result)
}

func TestTextInput_JoinHorizontal(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Type: config.TypeTextInput,
		Name: "test",
	}

	ti, err := NewTextInput(cfg, theme)
	require.NoError(t, err)

	// Test the joinHorizontal method by calling View which uses it internally
	view := ti.View()
	assert.NotEmpty(t, view)
}
