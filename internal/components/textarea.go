package components

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/v2/textarea"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/errors"
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

	// Error management integration
	errorManager *errors.ErrorManager
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

// SetErrorManager configura o ErrorManager para o componente
func (t *TextArea) SetErrorManager(em *errors.ErrorManager) {
	t.errorManager = em
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
	var parts []string

	// Render label
	if t.label != "" {
		labelStyle := t.theme.Label
		if t.errorMsg != "" {
			labelStyle = t.theme.LabelError
		}
		parts = append(parts, labelStyle.Render(t.label))
	}

	// Render textarea
	parts = append(parts, t.model.View())

	// Render error or help text
	if t.errorMsg != "" {
		parts = append(parts, t.theme.Error.Render("✗ "+t.errorMsg))
	} else if t.help != "" {
		parts = append(parts, t.theme.Help.Render(t.help))
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
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

	// Required validation with ErrorManager integration
	if t.required && strings.TrimSpace(value) == "" {
		t.errorMsg = "Este campo é obrigatório"

		if t.errorManager != nil {
			log.Printf("TextArea validation error in %s: campo obrigatório não preenchido", t.name)
		}
		return false
	}

	// Skip other validations if empty and not required
	if value == "" {
		t.errorMsg = ""
		return true
	}

	// Min length validation with ErrorManager integration
	if t.minLength > 0 && len(value) < t.minLength {
		t.errorMsg = fmt.Sprintf("Mínimo de %d caracteres", t.minLength)

		if t.errorManager != nil {
			log.Printf("TextArea min length validation error in %s: valor abaixo do mínimo", t.name)
		}
		return false
	}

	// Max length validation with ErrorManager integration
	if t.maxLength > 0 && len(value) > t.maxLength {
		t.errorMsg = fmt.Sprintf("Máximo de %d caracteres", t.maxLength)

		if t.errorManager != nil {
			log.Printf("TextArea max length validation error in %s: valor excede o máximo", t.name)
		}
		return false
	}

	// Clear error if validation passes
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
		err := fmt.Errorf("valor inválido: esperado string, recebido %T", value)

		if t.errorManager != nil {
			log.Printf("TextArea type validation error in %s: tipo inválido", t.name)
		}

		return err
	}

	t.model.SetValue(strValue)

	// Clear any previous error when setting a valid value
	t.errorMsg = ""

	return nil
}

// Reset implements Component.
func (t *TextArea) Reset() {
	t.model.SetValue(t.initialValue)
	t.errorMsg = ""
	t.model.Blur()
	t.focused = false
}

// GetMetadata implements Component.
func (t *TextArea) GetMetadata() ComponentMetadata {
	return ComponentMetadata{
		Version:      "1.0.0",
		Author:       "Shantilly Team",
		Description:  "Multi-line text area component with validation support",
		Dependencies: []string{},
		Examples: []ComponentExample{
			{
				Name:        "Comment Field",
				Description: "Multi-line text area for user comments",
				Config: map[string]interface{}{
					"type":        "textarea",
					"name":        "comments",
					"label":       "Comments",
					"placeholder": "Enter your comments here...",
				},
			},
		},
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"value": map[string]interface{}{
					"type":        "string",
					"description": "The text value",
				},
			},
		},
	}
}

// ValidateWithContext implements Component.
func (t *TextArea) ValidateWithContext(context ValidationContext) []ValidationError {
	var errors []ValidationError

	// Basic validation with ErrorManager integration
	if !t.IsValid() {
		validationErr := ValidationError{
			Code:     "VALIDATION_FAILED",
			Message:  t.GetError(),
			Field:    t.name,
			Severity: "error",
			Context: map[string]interface{}{
				"component":          "TextArea",
				"value":              t.Value(),
				"validation_context": context,
				"required":           t.required,
				"min_length":         t.minLength,
				"max_length":         t.maxLength,
			},
		}
		errors = append(errors, validationErr)

		// Log to ErrorManager if available
		if t.errorManager != nil {
			log.Printf("TextArea validation failed in %s: %s", t.name, t.GetError())
		}
	}

	return errors
}

// ExportToFormat implements Component.
func (t *TextArea) ExportToFormat(format ExportFormat) ([]byte, error) {
	data := map[string]interface{}{
		"name":     t.Name(),
		"value":    t.Value(),
		"metadata": t.GetMetadata(),
	}

	switch format {
	case FormatJSON:
		return json.MarshalIndent(data, "", "  ")
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

// ImportFromFormat implements Component.
func (t *TextArea) ImportFromFormat(format ExportFormat, data []byte) error {
	var imported map[string]interface{}

	switch format {
	case FormatJSON:
		if err := json.Unmarshal(data, &imported); err != nil {
			return fmt.Errorf("erro ao fazer parse do JSON: %w", err)
		}
	default:
		return fmt.Errorf("formato não suportado: %s", format)
	}

	if value, ok := imported["value"].(string); ok {
		return t.SetValue(value)
	}

	return nil
}

// GetDependencies implements Component.
func (t *TextArea) GetDependencies() []string {
	return []string{}
}

// SetTheme implements Component.
func (t *TextArea) SetTheme(theme *styles.Theme) {
	t.theme = theme
}
