package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSlider(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name        string
		cfg         config.ComponentConfig
		expectError bool
		validate    func(*testing.T, *Slider)
	}{
		{
			name: "valid slider creation with default",
			cfg: config.ComponentConfig{
				Name:    "volume",
				Type:    config.TypeSlider,
				Label:   "Volume",
				Default: 50,
			},
			expectError: false,
			validate: func(t *testing.T, s *Slider) {
				assert.Equal(t, "volume", s.name)
				assert.Equal(t, 50.0, s.value)
				assert.Equal(t, 50.0, s.initialValue)
				assert.Equal(t, 0.0, s.min)
				assert.Equal(t, 100.0, s.max)
			},
		},
		{
			name: "slider with custom options as float",
			cfg: config.ComponentConfig{
				Name: "brightness",
				Type: config.TypeSlider,
				Options: map[string]interface{}{
					"min":   10.0,
					"max":   255.0,
					"step":  5.0,
					"width": 40,
				},
			},
			expectError: false,
			validate: func(t *testing.T, s *Slider) {
				assert.Equal(t, 10.0, s.min)
				assert.Equal(t, 255.0, s.max)
				assert.Equal(t, 5.0, s.step)
				assert.Equal(t, 40, s.width)
			},
		},
		{
			name: "slider with custom options as int",
			cfg: config.ComponentConfig{
				Name: "contrast",
				Type: config.TypeSlider,
				Options: map[string]interface{}{
					"min":  -100,
					"max":  100,
					"step": 10,
				},
			},
			expectError: false,
			validate: func(t *testing.T, s *Slider) {
				assert.Equal(t, -100.0, s.min)
				assert.Equal(t, 100.0, s.max)
				assert.Equal(t, 10.0, s.step)
			},
		},
		{
			name: "default value outside range is ignored",
			cfg: config.ComponentConfig{
				Name:    "test",
				Type:    config.TypeSlider,
				Options: map[string]interface{}{"min": 10.0, "max": 20.0},
				Default: 5.0,
			},
			expectError: false,
			validate: func(t *testing.T, s *Slider) {
				assert.Equal(t, 10.0, s.value, "Default value below min should be ignored, defaulting to min")
			},
		},
		{
			name:        "invalid component type",
			cfg:         config.ComponentConfig{Name: "invalid", Type: config.TypeTextInput},
			expectError: true,
		},
		{
			name: "invalid range min >= max",
			cfg: config.ComponentConfig{
				Name:    "invalid-range",
				Type:    config.TypeSlider,
				Options: map[string]interface{}{"min": 100, "max": 50},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewSlider(tt.cfg, theme)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, s)
			} else {
				require.NoError(t, err)
				require.NotNil(t, s)
				if tt.validate != nil {
					tt.validate(t, s)
				}
			}
		})
	}
}

func TestSlider_Init(t *testing.T) {
	cfg := config.ComponentConfig{Name: "test", Type: config.TypeSlider}
	s, err := NewSlider(cfg, styles.DefaultTheme())
	require.NoError(t, err)
	assert.Nil(t, s.Init())
}

func TestSlider_Update(t *testing.T) {
	cfg := config.ComponentConfig{
		Name:    "test",
		Type:    config.TypeSlider,
		Options: map[string]interface{}{"min": 0, "max": 10, "step": 1},
		Default: 5,
	}
	s, err := NewSlider(cfg, styles.DefaultTheme())
	require.NoError(t, err)

	leftMsg := tea.KeyPressMsg{Code: tea.KeyLeft}
	rightMsg := tea.KeyPressMsg{Code: tea.KeyRight}
	homeMsg := tea.KeyPressMsg{Code: tea.KeyHome}
	endMsg := tea.KeyPressMsg{Code: tea.KeyEnd}
	hMsg := tea.KeyPressMsg{Code: 'h'}
	lMsg := tea.KeyPressMsg{Code: 'l'}

	t.Run("should not update when not focused", func(t *testing.T) {
		err := s.SetValue(5)
		require.NoError(t, err)
		s.SetFocus(false)
		s.Update(rightMsg)
		assert.Equal(t, 5.0, s.value)
	})

	t.Run("should update when focused", func(t *testing.T) {
		err := s.SetValue(5)
		require.NoError(t, err)
		s.SetFocus(true)
		s.SetError("an error")

		s.Update(rightMsg)
		assert.Equal(t, 6.0, s.value)
		assert.Empty(t, s.errorMsg, "Error should be cleared on update")

		s.Update(leftMsg)
		assert.Equal(t, 5.0, s.value)
	})

	t.Run("should respect boundaries", func(t *testing.T) {
		s.SetFocus(true)
		err := s.SetValue(0)
		require.NoError(t, err)
		s.Update(leftMsg)
		assert.Equal(t, 0.0, s.value, "Should not go below min")

		err = s.SetValue(10)
		require.NoError(t, err)
		s.Update(rightMsg)
		assert.Equal(t, 10.0, s.value, "Should not go above max")
	})

	t.Run("should handle vim keys h and l", func(t *testing.T) {
		s.SetFocus(true)
		err := s.SetValue(5)
		require.NoError(t, err)
		s.Update(hMsg)
		assert.Equal(t, 4.0, s.value)

		s.Update(lMsg)
		assert.Equal(t, 5.0, s.value)
	})

	t.Run("should handle home and end keys", func(t *testing.T) {
		s.SetFocus(true)
		err := s.SetValue(5)
		require.NoError(t, err)
		s.Update(homeMsg)
		assert.Equal(t, 0.0, s.value, "Home key should go to min")

		s.Update(endMsg)
		assert.Equal(t, 10.0, s.value, "End key should go to max")
	})
}

func TestSlider_View(t *testing.T) {
	cfg := config.ComponentConfig{
		Name:    "test",
		Type:    config.TypeSlider,
		Label:   "Test Slider",
		Help:    "Helpful text",
		Options: map[string]interface{}{"min": 0, "max": 100, "width": 10},
	}
	s, err := NewSlider(cfg, styles.DefaultTheme())
	require.NoError(t, err)

	t.Run("default view", func(t *testing.T) {
		err := s.SetValue(50)
		require.NoError(t, err)
		view := s.View()
		assert.Contains(t, view, "Test Slider")
		assert.Contains(t, view, "50.0")
		assert.Contains(t, view, "Min: 0.0 | Max: 100.0")
		assert.Contains(t, view, "Helpful text")
	})

	t.Run("error view hides help text", func(t *testing.T) {
		s.SetError("an error")
		view := s.View()
		assert.Contains(t, view, "✗ an error")
		assert.NotContains(t, view, "Helpful text")
	})
}

func TestSlider_ComponentInterface(t *testing.T) {
	cfg := config.ComponentConfig{Name: "test-slider", Type: config.TypeSlider}
	s, err := NewSlider(cfg, styles.DefaultTheme())
	require.NoError(t, err)

	assert.Equal(t, "test-slider", s.Name())
	assert.True(t, s.CanFocus())
	assert.True(t, s.IsValid())

	s.SetFocus(true)
	assert.True(t, s.focused)

	s.SetError("error")
	assert.Equal(t, "error", s.GetError())

	t.Run("SetValue", func(t *testing.T) {
		assert.NoError(t, s.SetValue(50.0))
		assert.Equal(t, 50.0, s.Value())
		assert.NoError(t, s.SetValue(25))
		assert.Equal(t, 25.0, s.Value())

		assert.Error(t, s.SetValue("not a number"))
		assert.Error(t, s.SetValue(200.0), "Value above max should error")
		assert.Error(t, s.SetValue(-10.0), "Value below min should error")
	})
}

func TestSlider_Reset(t *testing.T) {
	t.Run("resets to default value", func(t *testing.T) {
		cfg := config.ComponentConfig{Name: "test", Type: config.TypeSlider, Default: 75}
		s, err := NewSlider(cfg, styles.DefaultTheme())
		require.NoError(t, err)

		err = s.SetValue(25)
		require.NoError(t, err)
		s.SetFocus(true)
		s.Reset()

		assert.Equal(t, 75.0, s.value)
		assert.False(t, s.focused)
	})

	t.Run("resets to min value when no default", func(t *testing.T) {
		cfg := config.ComponentConfig{Name: "test", Type: config.TypeSlider, Options: map[string]interface{}{"min": 10}}
		s, err := NewSlider(cfg, styles.DefaultTheme())
		require.NoError(t, err)

		err = s.SetValue(50)
		require.NoError(t, err)
		s.SetFocus(true)
		s.Reset()

		assert.Equal(t, 10.0, s.value)
		assert.False(t, s.focused)
	})
}
