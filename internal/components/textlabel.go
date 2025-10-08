package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
)

// TextLabel implements a static text label component.
// This is used for displaying text that doesn't require user interaction.
type TextLabel struct {
	name  string
	text  string
	theme *styles.Theme
}

// NewTextLabel creates a new TextLabel component from configuration.
func NewTextLabel(cfg config.ComponentConfig, theme *styles.Theme) (*TextLabel, error) {
	if cfg.Type != config.TypeText {
		return nil, fmt.Errorf("tipo de componente inválido: esperado text, recebido %s", cfg.Type)
	}

	text := cfg.Label
	if text == "" && cfg.Default != nil {
		if defaultStr, ok := cfg.Default.(string); ok {
			text = defaultStr
		}
	}

	return &TextLabel{
		name:  cfg.Name,
		text:  text,
		theme: theme,
	}, nil
}

// Init implements tea.Model.
func (t *TextLabel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (t *TextLabel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Static component, no updates needed
	return t, nil
}

// View implements tea.Model.
func (t *TextLabel) View() string {
	return t.theme.Label.Render(t.text) + "\n"
}

// Name implements Component.
func (t *TextLabel) Name() string {
	return t.name
}

// CanFocus implements Component.
func (t *TextLabel) CanFocus() bool {
	return false // Static component cannot receive focus
}

// SetFocus implements Component.
func (t *TextLabel) SetFocus(focused bool) {
	// No-op for static component
}

// IsValid implements Component.
func (t *TextLabel) IsValid() bool {
	return true // Static component is always valid
}

// GetError implements Component.
func (t *TextLabel) GetError() string {
	return "" // Static component never has errors
}

// SetError implements Component.
func (t *TextLabel) SetError(msg string) {
	// No-op for static component
}

// Value implements Component.
func (t *TextLabel) Value() interface{} {
	return t.text
}

// SetValue implements Component.
func (t *TextLabel) SetValue(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("valor inválido: esperado string, recebido %T", value)
	}

	t.text = strValue
	return nil
}

// Reset implements Component.
func (t *TextLabel) Reset() {
	// No-op for static component
}
