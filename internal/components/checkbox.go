package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
)

// Checkbox implements a checkbox component.
// Note: bubbles/v2 doesn't have a checkbox yet, so we implement it with custom logic.
type Checkbox struct {
	name         string
	label        string
	required     bool
	help         string
	checked      bool
	theme        *styles.Theme
	errorMsg     string
	focused      bool
	initialValue bool
}

// NewCheckbox creates a new Checkbox component from configuration.
func NewCheckbox(cfg config.ComponentConfig, theme *styles.Theme) (*Checkbox, error) {
	if cfg.Type != config.TypeCheckbox {
		return nil, fmt.Errorf("tipo de componente inválido: esperado checkbox, recebido %s", cfg.Type)
	}

	c := &Checkbox{
		name:     cfg.Name,
		label:    cfg.Label,
		required: cfg.Required,
		help:     cfg.Help,
		theme:    theme,
	}

	// Set default value if provided
	if cfg.Default != nil {
		if defaultBool, ok := cfg.Default.(bool); ok {
			c.checked = defaultBool
			c.initialValue = defaultBool
		}
	}

	return c, nil
}

// Init implements tea.Model.
func (c *Checkbox) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (c *Checkbox) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !c.focused {
		return c, nil
	}

	switch msg := msg.(type) {
	// tea.KeyMsg is used for special keys, like space and enter.
	case tea.KeyMsg:
		switch msg.String() {
		case " ", "space", "enter":
			c.checked = !c.checked
			c.errorMsg = ""
		}
	}

	return c, nil
}

// View implements tea.Model.
func (c *Checkbox) View() string {
	var b strings.Builder

	// Determine checkbox symbol
	var symbol string
	if c.checked {
		symbol = "[✓]"
	} else {
		symbol = "[ ]"
	}

	// Build the checkbox line (without border - border is applied by layout)
	checkboxLine := symbol + " " + c.label
	b.WriteString(checkboxLine)

	// Render error message if present
	if c.errorMsg != "" {
		b.WriteString("\n")
		b.WriteString(c.theme.Error.Render("✗ " + c.errorMsg))
	}

	// Render help text if present and no error
	if c.help != "" && c.errorMsg == "" {
		b.WriteString("\n")
		b.WriteString(c.theme.Help.Render(c.help))
	}

	return b.String()
}

// Name implements Component.
func (c *Checkbox) Name() string {
	return c.name
}

// CanFocus implements Component.
func (c *Checkbox) CanFocus() bool {
	return true
}

// SetFocus implements Component.
func (c *Checkbox) SetFocus(focused bool) {
	c.focused = focused
}

// IsValid implements Component.
func (c *Checkbox) IsValid() bool {
	// For checkboxes, required means it must be checked
	if c.required && !c.checked {
		c.errorMsg = "Esta opção deve ser marcada"
		return false
	}

	c.errorMsg = ""
	return true
}

// GetError implements Component.
func (c *Checkbox) GetError() string {
	return c.errorMsg
}

// SetError implements Component.
func (c *Checkbox) SetError(msg string) {
	c.errorMsg = msg
}

// Value implements Component.
func (c *Checkbox) Value() interface{} {
	return c.checked
}

// SetValue implements Component.
func (c *Checkbox) SetValue(value interface{}) error {
	boolValue, ok := value.(bool)
	if !ok {
		return fmt.Errorf("valor inválido: esperado bool, recebido %T", value)
	}

	c.checked = boolValue
	return nil
}

// Reset implements Component.
func (c *Checkbox) Reset() {
	c.checked = c.initialValue
	c.errorMsg = ""
	c.focused = false
}
