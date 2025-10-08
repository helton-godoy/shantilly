package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
)

// RadioItem represents a single radio button option.
type RadioItem struct {
	ID    string
	Label string
}

// RadioGroup implements a radio button group component with custom Lip Gloss rendering.
type RadioGroup struct {
	name         string
	label        string
	required     bool
	help         string
	items        []RadioItem
	cursor       int // Current cursor position
	selected     int // Selected item index (-1 = none)
	theme        *styles.Theme
	errorMsg     string
	focused      bool
	initialValue int
}

// NewRadioGroup creates a new RadioGroup component from configuration.
func NewRadioGroup(cfg config.ComponentConfig, theme *styles.Theme) (*RadioGroup, error) {
	if cfg.Type != config.TypeRadioGroup {
		return nil, fmt.Errorf("tipo de componente inválido: esperado radiogroup, recebido %s", cfg.Type)
	}

	// Parse items from options
	var items []RadioItem
	if cfg.Options != nil {
		if itemsData, ok := cfg.Options["items"].([]interface{}); ok {
			for _, item := range itemsData {
				if itemMap, ok := item.(map[string]interface{}); ok {
					id, okID := itemMap["id"].(string)
					if !okID {
						return nil, fmt.Errorf("campo 'id' deve ser string: %T", itemMap["id"])
					}
					label, okLabel := itemMap["label"].(string)
					if !okLabel {
						return nil, fmt.Errorf("campo 'label' deve ser string: %T", itemMap["label"])
					}
					if id != "" && label != "" {
						items = append(items, RadioItem{ID: id, Label: label})
					}
				}
			}
		}
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("radiogroup deve conter pelo menos um item")
	}

	rg := &RadioGroup{
		name:         cfg.Name,
		label:        cfg.Label,
		required:     cfg.Required,
		help:         cfg.Help,
		items:        items,
		cursor:       0,
		selected:     -1, // None selected by default
		theme:        theme,
		initialValue: -1,
	}

	// Set default value if provided
	if cfg.Default != nil {
		if defaultID, ok := cfg.Default.(string); ok {
			for i, item := range items {
				if item.ID == defaultID {
					rg.selected = i
					rg.initialValue = i
					break
				}
			}
		}
	}

	return rg, nil
}

// Init implements tea.Model.
func (rg *RadioGroup) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (rg *RadioGroup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !rg.focused {
		return rg, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyUp, 'k':
			if rg.cursor > 0 {
				rg.cursor--
			}
		case tea.KeyDown, 'j':
			if rg.cursor < len(rg.items)-1 {
				rg.cursor++
			}
		case tea.KeyEnter, tea.KeySpace:
			rg.selected = rg.cursor
			rg.errorMsg = ""
		}
	}

	return rg, nil
}

// View implements tea.Model.
func (rg *RadioGroup) View() string {
	var b strings.Builder

	// Render label
	if rg.label != "" {
		labelStyle := rg.theme.Label
		if rg.errorMsg != "" {
			labelStyle = rg.theme.LabelError
		}
		b.WriteString(labelStyle.Render(rg.label))
		b.WriteString("\n")
	}

	// Render radio items
	for i, item := range rg.items {
		var symbol string
		var line string

		// Determine symbol based on selection
		if i == rg.selected {
			symbol = "(•)"
		} else {
			symbol = "( )"
		}

		line = symbol + " " + item.Label

		// Apply style based on cursor position and focus
		if rg.focused && i == rg.cursor {
			line = rg.theme.RadioSelected.Render(line)
		} else if i == rg.selected {
			line = rg.theme.RadioSelected.Render(line)
		} else {
			line = rg.theme.RadioUnselected.Render(line)
		}

		b.WriteString(line)
		if i < len(rg.items)-1 {
			b.WriteString("\n")
		}
	}

	// Render error message if present
	if rg.errorMsg != "" {
		b.WriteString("\n")
		b.WriteString(rg.theme.Error.Render("✗ " + rg.errorMsg))
	}

	// Render help text if present and no error
	if rg.help != "" && rg.errorMsg == "" {
		b.WriteString("\n")
		b.WriteString(rg.theme.Help.Render(rg.help))
	}

	return b.String()
}

// Name implements Component.
func (rg *RadioGroup) Name() string {
	return rg.name
}

// CanFocus implements Component.
func (rg *RadioGroup) CanFocus() bool {
	return true
}

// SetFocus implements Component.
func (rg *RadioGroup) SetFocus(focused bool) {
	rg.focused = focused
}

// IsValid implements Component.
func (rg *RadioGroup) IsValid() bool {
	// Required validation: must have a selection
	if rg.required && rg.selected == -1 {
		rg.errorMsg = "Selecione uma opção"
		return false
	}

	rg.errorMsg = ""
	return true
}

// GetError implements Component.
func (rg *RadioGroup) GetError() string {
	return rg.errorMsg
}

// SetError implements Component.
func (rg *RadioGroup) SetError(msg string) {
	rg.errorMsg = msg
}

// Value implements Component.
// Returns the ID of the selected item, not the index.
func (rg *RadioGroup) Value() interface{} {
	if rg.selected == -1 {
		return ""
	}
	return rg.items[rg.selected].ID
}

// SetValue implements Component.
func (rg *RadioGroup) SetValue(value interface{}) error {
	idValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("valor inválido: esperado string (ID), recebido %T", value)
	}

	// Find item by ID
	for i, item := range rg.items {
		if item.ID == idValue {
			rg.selected = i
			rg.cursor = i
			return nil
		}
	}

	return fmt.Errorf("ID não encontrado: %s", idValue)
}

// Reset implements Component.
func (rg *RadioGroup) Reset() {
	rg.selected = rg.initialValue
	rg.cursor = 0
	if rg.initialValue != -1 {
		rg.cursor = rg.initialValue
	}
	rg.errorMsg = ""
	rg.focused = false
}
