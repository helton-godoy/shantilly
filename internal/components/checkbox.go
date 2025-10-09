package components

import (
	"encoding/json"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/errors"
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

	// Error management integration
	errorManager *errors.ErrorManager
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

// SetErrorManager configura o ErrorManager para o componente
func (c *Checkbox) SetErrorManager(em *errors.ErrorManager) {
	c.errorManager = em
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
	var parts []string

	// Determine checkbox symbol
	var symbol string
	if c.checked {
		symbol = "[✓]"
	} else {
		symbol = "[ ]"
	}

	// Build the checkbox line
	checkboxLine := symbol + " " + c.label
	parts = append(parts, checkboxLine)

	// Render error or help text
	if c.errorMsg != "" {
		parts = append(parts, c.theme.Error.Render("✗ "+c.errorMsg))
	} else if c.help != "" {
		parts = append(parts, c.theme.Help.Render(c.help))
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
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

		if c.errorManager != nil {
			log.Printf("Checkbox validation error in %s: opção obrigatória não marcada", c.name)
		}
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
		err := fmt.Errorf("valor inválido: esperado bool, recebido %T", value)

		if c.errorManager != nil {
			log.Printf("Checkbox type validation error in %s: tipo inválido", c.name)
		}

		return err
	}

	c.checked = boolValue

	// Clear any previous error when setting a valid value
	c.errorMsg = ""

	return nil
}

// Reset implements Component.
func (c *Checkbox) Reset() {
	c.checked = c.initialValue
	c.errorMsg = ""
	c.focused = false
}

// GetMetadata implements Component.
func (c *Checkbox) GetMetadata() ComponentMetadata {
	return ComponentMetadata{
		Version:      "1.0.0",
		Author:       "Shantilly Team",
		Description:  "Checkbox component for boolean selections",
		Dependencies: []string{},
		Examples: []ComponentExample{
			{
				Name:        "Simple Checkbox",
				Description: "Basic checkbox for accepting terms",
				Config: map[string]interface{}{
					"type":     "checkbox",
					"name":     "accept_terms",
					"label":    "I accept the terms and conditions",
					"required": true,
				},
			},
		},
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"value": map[string]interface{}{
					"type":        "boolean",
					"description": "The checkbox state",
				},
			},
		},
	}
}

// ValidateWithContext implements Component.
func (c *Checkbox) ValidateWithContext(context ValidationContext) []ValidationError {
	var errors []ValidationError

	if !c.IsValid() {
		validationErr := ValidationError{
			Code:     "VALIDATION_FAILED",
			Message:  c.GetError(),
			Field:    c.name,
			Severity: "error",
			Context: map[string]interface{}{
				"component":          "Checkbox",
				"value":              c.Value(),
				"validation_context": context,
				"required":           c.required,
				"checked":            c.checked,
			},
		}
		errors = append(errors, validationErr)

		// Log to ErrorManager if available
		if c.errorManager != nil {
			log.Printf("Checkbox validation failed in %s: %s", c.name, c.GetError())
		}
	}

	return errors
}

// ExportToFormat implements Component.
func (c *Checkbox) ExportToFormat(format ExportFormat) ([]byte, error) {
	data := map[string]interface{}{
		"name":     c.Name(),
		"value":    c.Value(),
		"metadata": c.GetMetadata(),
	}

	switch format {
	case FormatJSON:
		return json.MarshalIndent(data, "", "  ")
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

// ImportFromFormat implements Component.
func (c *Checkbox) ImportFromFormat(format ExportFormat, data []byte) error {
	var imported map[string]interface{}

	switch format {
	case FormatJSON:
		if err := json.Unmarshal(data, &imported); err != nil {
			return fmt.Errorf("erro ao fazer parse do JSON: %w", err)
		}
	default:
		return fmt.Errorf("formato não suportado: %s", format)
	}

	if value, ok := imported["value"].(bool); ok {
		return c.SetValue(value)
	}

	return nil
}

// GetDependencies implements Component.
func (c *Checkbox) GetDependencies() []string {
	return []string{}
}

// SetTheme implements Component.
func (c *Checkbox) SetTheme(theme *styles.Theme) {
	c.theme = theme
}
