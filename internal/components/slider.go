package components

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/errors"
	"github.com/helton/shantilly/internal/styles"
)

// Slider implements a slider component with custom Lip Gloss rendering.
type Slider struct {
	name         string
	label        string
	required     bool
	help         string
	value        float64
	min          float64
	max          float64
	step         float64
	width        int
	theme        *styles.Theme
	errorMsg     string
	focused      bool
	initialValue float64

	// Error management integration
	errorManager *errors.ErrorManager
}

// NewSlider creates a new Slider component from configuration.
func NewSlider(cfg config.ComponentConfig, theme *styles.Theme) (*Slider, error) {
	if cfg.Type != config.TypeSlider {
		return nil, fmt.Errorf("tipo de componente inválido: esperado slider, recebido %s", cfg.Type)
	}

	s := &Slider{
		name:     cfg.Name,
		label:    cfg.Label,
		required: cfg.Required,
		help:     cfg.Help,
		min:      0.0,
		max:      100.0,
		step:     1.0,
		width:    30,
		theme:    theme,
	}

	// Parse options
	if cfg.Options != nil {
		if min, ok := cfg.Options["min"].(float64); ok {
			s.min = min
		} else if minInt, ok := cfg.Options["min"].(int); ok {
			s.min = float64(minInt)
		}

		if max, ok := cfg.Options["max"].(float64); ok {
			s.max = max
		} else if maxInt, ok := cfg.Options["max"].(int); ok {
			s.max = float64(maxInt)
		}

		if step, ok := cfg.Options["step"].(float64); ok {
			s.step = step
		} else if stepInt, ok := cfg.Options["step"].(int); ok {
			s.step = float64(stepInt)
		}

		if width, ok := cfg.Options["width"].(int); ok {
			s.width = width
		}
	}

	// Validate min/max
	if s.min >= s.max {
		return nil, fmt.Errorf("min deve ser menor que max")
	}

	// Set default value
	s.value = s.min
	s.initialValue = s.min

	if cfg.Default != nil {
		if defaultFloat, ok := cfg.Default.(float64); ok {
			if defaultFloat >= s.min && defaultFloat <= s.max {
				s.value = defaultFloat
				s.initialValue = defaultFloat
			}
		} else if defaultInt, ok := cfg.Default.(int); ok {
			defaultFloat := float64(defaultInt)
			if defaultFloat >= s.min && defaultFloat <= s.max {
				s.value = defaultFloat
				s.initialValue = defaultFloat
			}
		}
	}

	return s, nil
}

// SetErrorManager configura o ErrorManager para o componente
func (s *Slider) SetErrorManager(em *errors.ErrorManager) {
	s.errorManager = em
}

// Init implements tea.Model.
func (s *Slider) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (s *Slider) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !s.focused {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyLeft, 'h':
			s.value -= s.step
			if s.value < s.min {
				s.value = s.min
			}
			s.errorMsg = ""
		case tea.KeyRight, 'l':
			s.value += s.step
			if s.value > s.max {
				s.value = s.max
			}
			s.errorMsg = ""
		case tea.KeyHome:
			s.value = s.min
			s.errorMsg = ""
		case tea.KeyEnd:
			s.value = s.max
			s.errorMsg = ""
		}
	}

	return s, nil
}

// View implements tea.Model.
func (s *Slider) View() string {
	var b strings.Builder

	// Render label
	if s.label != "" {
		labelStyle := s.theme.Label
		if s.errorMsg != "" {
			labelStyle = s.theme.LabelError
		}
		b.WriteString(labelStyle.Render(s.label))
		b.WriteString("\n")
	}

	// Calculate position
	percentage := (s.value - s.min) / (s.max - s.min)
	filledWidth := int(float64(s.width) * percentage)
	if filledWidth < 0 {
		filledWidth = 0
	}
	if filledWidth > s.width {
		filledWidth = s.width
	}

	// Build slider bar
	filled := strings.Repeat("━", filledWidth)
	empty := strings.Repeat("━", s.width-filledWidth)

	// Apply styles to bar segments
	filledBar := s.theme.SliderFilled.Render(filled)
	emptyBar := s.theme.SliderBar.Render(empty)

	// Build slider line (without container border - border is applied by layout)
	sliderLine := filledBar + emptyBar + fmt.Sprintf(" %.1f", s.value)
	b.WriteString(sliderLine)
	b.WriteString("\n")

	// Render min/max labels
	rangeLabel := s.theme.Help.Render(fmt.Sprintf("Min: %.1f | Max: %.1f", s.min, s.max))
	b.WriteString(rangeLabel)

	// Render error message if present
	if s.errorMsg != "" {
		b.WriteString("\n")
		b.WriteString(s.theme.Error.Render("✗ " + s.errorMsg))
	}

	// Render help text if present and no error
	if s.help != "" && s.errorMsg == "" {
		b.WriteString("\n")
		b.WriteString(s.theme.Help.Render(s.help))
	}

	return b.String()
}

// Name implements Component.
func (s *Slider) Name() string {
	return s.name
}

// CanFocus implements Component.
func (s *Slider) CanFocus() bool {
	return true
}

// SetFocus implements Component.
func (s *Slider) SetFocus(focused bool) {
	s.focused = focused
}

// IsValid implements Component.
func (s *Slider) IsValid() bool {
	// Slider is always valid since value is constrained by min/max
	s.errorMsg = ""
	return true
}

// GetError implements Component.
func (s *Slider) GetError() string {
	return s.errorMsg
}

// SetError implements Component.
func (s *Slider) SetError(msg string) {
	s.errorMsg = msg
}

// Value implements Component.
func (s *Slider) Value() interface{} {
	return s.value
}

// SetValue implements Component.
func (s *Slider) SetValue(value interface{}) error {
	var floatValue float64

	switch v := value.(type) {
	case float64:
		floatValue = v
	case int:
		floatValue = float64(v)
	case int64:
		floatValue = float64(v)
	default:
		err := fmt.Errorf("valor inválido: esperado número, recebido %T", value)

		if s.errorManager != nil {
			log.Printf("Slider type validation error in %s: tipo inválido", s.name)
		}

		return err
	}

	if floatValue < s.min || floatValue > s.max {
		err := fmt.Errorf("valor fora do intervalo [%.1f, %.1f]", s.min, s.max)

		if s.errorManager != nil {
			log.Printf("Slider range validation error in %s: valor fora do intervalo", s.name)
		}

		return err
	}

	s.value = floatValue

	// Clear any previous error when setting a valid value
	s.errorMsg = ""

	return nil
}

// Reset implements Component.
func (s *Slider) Reset() {
	s.value = s.initialValue
	s.errorMsg = ""
	s.focused = false
}

// GetMetadata implements Component.
func (s *Slider) GetMetadata() ComponentMetadata {
	return ComponentMetadata{
		Version:      "1.0.0",
		Author:       "Shantilly Team",
		Description:  "Slider component for numeric selections",
		Dependencies: []string{},
		Examples: []ComponentExample{
			{
				Name:        "Volume Control",
				Description: "Slider for volume control",
				Config: map[string]interface{}{
					"type":  "slider",
					"name":  "volume",
					"label": "Volume",
					"options": map[string]interface{}{
						"min": 0,
						"max": 100,
					},
				},
			},
		},
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"value": map[string]interface{}{
					"type":        "number",
					"description": "The slider value",
				},
			},
		},
	}
}

// ValidateWithContext implements Component.
func (s *Slider) ValidateWithContext(context ValidationContext) []ValidationError {
	var errors []ValidationError

	if !s.IsValid() {
		validationErr := ValidationError{
			Code:     "VALIDATION_FAILED",
			Message:  s.GetError(),
			Field:    s.name,
			Severity: "error",
			Context: map[string]interface{}{
				"component":          "Slider",
				"value":              s.Value(),
				"validation_context": context,
				"required":           s.required,
				"min":                s.min,
				"max":                s.max,
				"step":               s.step,
			},
		}
		errors = append(errors, validationErr)

		// Log to ErrorManager if available
		if s.errorManager != nil {
			log.Printf("Slider validation failed in %s: %s", s.name, s.GetError())
		}
	}

	return errors
}

// ExportToFormat implements Component.
func (s *Slider) ExportToFormat(format ExportFormat) ([]byte, error) {
	data := map[string]interface{}{
		"name":     s.Name(),
		"value":    s.Value(),
		"metadata": s.GetMetadata(),
	}

	switch format {
	case FormatJSON:
		return json.MarshalIndent(data, "", "  ")
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

// ImportFromFormat implements Component.
func (s *Slider) ImportFromFormat(format ExportFormat, data []byte) error {
	var imported map[string]interface{}

	switch format {
	case FormatJSON:
		if err := json.Unmarshal(data, &imported); err != nil {
			return fmt.Errorf("erro ao fazer parse do JSON: %w", err)
		}
	default:
		return fmt.Errorf("formato não suportado: %s", format)
	}

	if value, ok := imported["value"].(float64); ok {
		return s.SetValue(value)
	}

	return nil
}

// GetDependencies implements Component.
func (s *Slider) GetDependencies() []string {
	return []string{}
}

// SetTheme implements Component.
func (s *Slider) SetTheme(theme *styles.Theme) {
	s.theme = theme
}

// JoinHorizontal is a helper for lipgloss compatibility.
func joinHorizontal(parts ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}
