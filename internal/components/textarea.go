package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/v2/textarea"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
)

// TextArea wraps bubbles/textarea and implements the Component interface.
// It provides a multi-line text input with validation support.
type TextArea struct {
	name         string
	label        string
	required     bool
	help         string
	model        textarea.Model
	theme        *styles.Theme
	errorMsg     string
	focused      bool
	initialValue string

	// Validation options
	minLength int
	maxLength int
}

// NewTextArea creates a new TextArea component from configuration.
func NewTextArea(cfg config.ComponentConfig, theme *styles.Theme) (*TextArea, error) {
	if cfg.Type != config.TypeTextArea {
		return nil, fmt.Errorf("tipo de componente inválido: esperado textarea, recebido %s", cfg.Type)
	}

	// Initialize bubbles textarea model
	ta := textarea.New()
	ta.Placeholder = cfg.Placeholder
	ta.ShowLineNumbers = false
	ta.CharLimit = 0 // No default limit

	// Set default value if provided
	if cfg.Default != nil {
		if defaultStr, ok := cfg.Default.(string); ok {
			ta.SetValue(defaultStr)
		}
	}

	t := &TextArea{
		name:         cfg.Name,
		label:        cfg.Label,
		required:     cfg.Required,
		help:         cfg.Help,
		model:        ta,
		theme:        theme,
		initialValue: ta.Value(),
	}

	// Parse validation options
	if cfg.Options != nil {
		if minLen, ok := cfg.Options["min_length"].(int); ok {
			t.minLength = minLen
		}
		if maxLen, ok := cfg.Options["max_length"].(int); ok {
			t.maxLength = maxLen
			ta.CharLimit = maxLen
		}
		if height, ok := cfg.Options["height"].(int); ok {
			ta.SetHeight(height)
		} else {
			ta.SetHeight(5) // Default height
		}
		if width, ok := cfg.Options["width"].(int); ok {
			ta.SetWidth(width)
		} else {
			ta.SetWidth(50) // Default width
		}
	}

	return t, nil
}

// Init implements tea.Model.
func (t *TextArea) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (t *TextArea) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle window size messages for responsive layout
	if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
		// Adjust width but preserve minimum
		newWidth := wsMsg.Width - 10
		if newWidth < 30 {
			newWidth = 30
		}
		t.model.SetWidth(newWidth)
	}

	// Only process input messages if focused
	if t.focused {
		t.model, cmd = t.model.Update(msg)
		// Clear error when user types
		if _, ok := msg.(tea.KeyPressMsg); ok {
			t.errorMsg = ""
		}
	}

	return t, cmd
}

// View implements tea.Model.
func (t *TextArea) View() string {
	var b strings.Builder

	// Render label
	if t.label != "" {
		labelStyle := t.theme.Label
		if t.errorMsg != "" {
			labelStyle = t.theme.LabelError
		}
		b.WriteString(labelStyle.Render(t.label))
		b.WriteString("\n")
	}

	// Render textarea (without border - border is applied by layout)
	b.WriteString(t.model.View())

	// Render error message if present
	if t.errorMsg != "" {
		b.WriteString("\n")
		b.WriteString(t.theme.Error.Render("✗ " + t.errorMsg))
	}

	// Render help text if present and no error
	if t.help != "" && t.errorMsg == "" {
		b.WriteString("\n")
		b.WriteString(t.theme.Help.Render(t.help))
	}

	return b.String()
}

// Name implements Component.
func (t *TextArea) Name() string {
	return t.name
}

// CanFocus implements Component.
func (t *TextArea) CanFocus() bool {
	return true
}

// SetFocus implements Component.
func (t *TextArea) SetFocus(focused bool) {
	t.focused = focused
	if focused {
		t.model.Focus()
	} else {
		t.model.Blur()
	}
}

// IsValid implements Component.
func (t *TextArea) IsValid() bool {
	value := t.model.Value()

	// Required validation
	if t.required && strings.TrimSpace(value) == "" {
		t.errorMsg = "Este campo é obrigatório"
		return false
	}

	// Skip other validations if empty and not required
	if value == "" {
		t.errorMsg = ""
		return true
	}

	// Min length validation
	if t.minLength > 0 && len(value) < t.minLength {
		t.errorMsg = fmt.Sprintf("Mínimo de %d caracteres", t.minLength)
		return false
	}

	// Max length validation
	if t.maxLength > 0 && len(value) > t.maxLength {
		t.errorMsg = fmt.Sprintf("Máximo de %d caracteres", t.maxLength)
		return false
	}

	t.errorMsg = ""
	return true
}

// GetError implements Component.
func (t *TextArea) GetError() string {
	return t.errorMsg
}

// SetError implements Component.
func (t *TextArea) SetError(msg string) {
	t.errorMsg = msg
}

// Value implements Component.
func (t *TextArea) Value() interface{} {
	return t.model.Value()
}

// SetValue implements Component.
func (t *TextArea) SetValue(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("valor inválido: esperado string, recebido %T", value)
	}

	t.model.SetValue(strValue)
	return nil
}

// Reset implements Component.
func (t *TextArea) Reset() {
	t.model.SetValue(t.initialValue)
	t.errorMsg = ""
	t.model.Blur()
	t.focused = false
}
