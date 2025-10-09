package components

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/errors"
	"github.com/helton/shantilly/internal/styles"
)

// TextInput wraps bubbles/textinput and implements the Component interface.
// It provides a single-line text input with validation support.
type TextInput struct {
	name         string
	label        string
	required     bool
	help         string
	model        textinput.Model
	theme        *styles.Theme
	errorMsg     string
	focused      bool
	initialValue string

	// Validation options
	minLength int
	maxLength int
	pattern   *regexp.Regexp

	// Error management integration
	errorManager *errors.ErrorManager
}

// NewTextInput creates a new TextInput component from configuration.
func NewTextInput(cfg config.ComponentConfig, theme *styles.Theme) (*TextInput, error) {
	if cfg.Type != config.TypeTextInput {
		return nil, fmt.Errorf("tipo de componente inválido: esperado textinput, recebido %s", cfg.Type)
	}

	// Initialize bubbles textinput model
	ti := textinput.New()
	ti.Placeholder = cfg.Placeholder
	ti.CharLimit = 0 // No default limit

	// Set default value if provided
	if cfg.Default != nil {
		if defaultStr, ok := cfg.Default.(string); ok {
			ti.SetValue(defaultStr)
		}
	}

	t := &TextInput{
		name:         cfg.Name,
		label:        cfg.Label,
		required:     cfg.Required,
		help:         cfg.Help,
		model:        ti,
		theme:        theme,
		initialValue: ti.Value(),
	}

	// Parse validation options
	if cfg.Options != nil {
		if minLen, ok := cfg.Options["min_length"].(int); ok {
			t.minLength = minLen
		}
		if maxLen, ok := cfg.Options["max_length"].(int); ok {
			t.maxLength = maxLen
			ti.CharLimit = maxLen
		}
		if patternStr, ok := cfg.Options["pattern"].(string); ok {
			pattern, err := regexp.Compile(patternStr)
			if err != nil {
				return nil, fmt.Errorf("erro ao compilar regex pattern: %w", err)
			}
			t.pattern = pattern
		}
	}

	return t, nil
}

// SetErrorManager configura o ErrorManager para o componente
func (t *TextInput) SetErrorManager(em *errors.ErrorManager) {
	t.errorManager = em
}

// Init implements tea.Model.
func (t *TextInput) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (t *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Only process messages if focused
	if t.focused {
		t.model, cmd = t.model.Update(msg)
		// Clear error when user types
		if _, ok := msg.(tea.KeyMsg); ok {
			t.errorMsg = ""
		}
	}

	return t, cmd
}

// View implements tea.Model.
func (t *TextInput) View() string {
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

	// Render input (without border - border is applied by layout)
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
func (t *TextInput) Name() string {
	return t.name
}

// CanFocus implements Component.
func (t *TextInput) CanFocus() bool {
	return true
}

// SetFocus implements Component.
func (t *TextInput) SetFocus(focused bool) {
	t.focused = focused
	if focused {
		t.model.Focus()
	} else {
		t.model.Blur()
	}
}

// IsValid implements Component.
func (t *TextInput) IsValid() bool {
	value := t.model.Value()

	// Required validation
	if t.required && strings.TrimSpace(value) == "" {
		t.errorMsg = "Este campo é obrigatório"

		if t.errorManager != nil {
			// Log validation error for debugging
			log.Printf("TextInput validation error: Campo obrigatório não preenchido - %s", t.name)
		}
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

		if t.errorManager != nil {
			log.Printf("TextInput min length validation error in %s: valor abaixo do mínimo", t.name)
		}
		return false
	}

	// Max length validation (already enforced by CharLimit, but check anyway)
	if t.maxLength > 0 && len(value) > t.maxLength {
		t.errorMsg = fmt.Sprintf("Máximo de %d caracteres", t.maxLength)

		if t.errorManager != nil {
			log.Printf("TextInput max length validation error in %s: valor excede o máximo", t.name)
		}
		return false
	}

	// Pattern validation
	if t.pattern != nil && !t.pattern.MatchString(value) {
		t.errorMsg = "Formato inválido"

		if t.errorManager != nil {
			log.Printf("TextInput pattern validation error in %s: formato inválido", t.name)
		}
		return false
	}

	// Clear error if validation passes
	t.errorMsg = ""
	return true
}

// GetError implements Component.
func (t *TextInput) GetError() string {
	return t.errorMsg
}

// SetError implements Component.
func (t *TextInput) SetError(msg string) {
	t.errorMsg = msg
}

// Value implements Component.
func (t *TextInput) Value() interface{} {
	return t.model.Value()
}

// SetValue implements Component.
func (t *TextInput) SetValue(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		err := fmt.Errorf("valor inválido: esperado string, recebido %T", value)

		if t.errorManager != nil {
			log.Printf("TextInput type validation error in %s: tipo inválido", t.name)
		}

		return err
	}

	t.model.SetValue(strValue)

	// Clear any previous error when setting a valid value
	t.errorMsg = ""

	return nil
}

// Reset implements Component.
func (t *TextInput) Reset() {
	t.model.SetValue(t.initialValue)
	t.errorMsg = ""
	t.model.Blur()
	t.focused = false
}

// GetMetadata implements Component.
func (t *TextInput) GetMetadata() ComponentMetadata {
	return ComponentMetadata{
		Version:      "1.0.0",
		Author:       "Shantilly Team",
		Description:  "Single-line text input component with validation support",
		Dependencies: []string{},
		Examples: []ComponentExample{
			{
				Name:        "Simple Text Input",
				Description: "Basic text input for user names",
				Config: map[string]interface{}{
					"type":        "textinput",
					"name":        "username",
					"label":       "Username",
					"placeholder": "Enter your username",
					"required":    true,
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
func (t *TextInput) ValidateWithContext(context ValidationContext) []ValidationError {
	var errors []ValidationError

	// Basic validation with ErrorManager integration
	if !t.IsValid() {
		validationErr := ValidationError{
			Code:     "VALIDATION_FAILED",
			Message:  t.GetError(),
			Field:    t.name,
			Severity: "error",
			Context: map[string]interface{}{
				"component":          "TextInput",
				"value":              t.Value(),
				"validation_context": context,
			},
		}
		errors = append(errors, validationErr)

		// Log to ErrorManager if available (simplified for compilation)
		if t.errorManager != nil {
			log.Printf("TextInput validation error in %s: %s", t.name, t.GetError())
		}
	}

	// Cross-field validation examples with ErrorManager
	if componentValues, ok := context.ComponentValues["password"]; ok {
		if password, ok := componentValues.(string); ok && t.name == "confirm_password" {
			if currentValue := t.Value().(string); currentValue != password {
				validationErr := ValidationError{
					Code:     "PASSWORD_MISMATCH",
					Message:  "Senhas não coincidem",
					Field:    t.name,
					Severity: "error",
					Context: map[string]interface{}{
						"component":       "TextInput",
						"related_field":   "password",
						"password_length": len(password),
						"current_length":  len(currentValue),
					},
				}
				errors = append(errors, validationErr)

				// Log cross-field validation error (simplified for compilation)
				if t.errorManager != nil {
					log.Printf("TextInput cross-field validation error in %s: senhas não coincidem", t.name)
				}
			}
		}
	}

	return errors
}

// ExportToFormat implements Component.
func (t *TextInput) ExportToFormat(format ExportFormat) ([]byte, error) {
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
func (t *TextInput) ImportFromFormat(format ExportFormat, data []byte) error {
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
func (t *TextInput) GetDependencies() []string {
	return []string{} // TextInput has no dependencies
}

// SetTheme implements Component.
func (t *TextInput) SetTheme(theme *styles.Theme) {
	t.theme = theme
}

// JoinVertical is a helper for lipgloss compatibility.
func joinVertical(parts ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
