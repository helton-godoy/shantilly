package components

import (
	"fmt"

	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
)

// NewComponent creates a component from configuration using the factory pattern.
// This ensures that all components are created consistently and with proper validation.
func NewComponent(cfg config.ComponentConfig, theme *styles.Theme) (Component, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("erro de validação da configuração: %w", err)
	}

	switch cfg.Type {
	case config.TypeTextInput:
		return NewTextInput(cfg, theme)
	case config.TypeTextArea:
		return NewTextArea(cfg, theme)
	case config.TypeCheckbox:
		return NewCheckbox(cfg, theme)
	case config.TypeRadioGroup:
		return NewRadioGroup(cfg, theme)
	case config.TypeSlider:
		return NewSlider(cfg, theme)
	case config.TypeText:
		return NewTextLabel(cfg, theme)
	default:
		return nil, fmt.Errorf("tipo de componente não suportado: %s", cfg.Type)
	}
}

// NewComponents creates multiple components from a slice of configurations.
// Returns an error if any component fails to create.
func NewComponents(configs []config.ComponentConfig, theme *styles.Theme) ([]Component, error) {
	components := make([]Component, 0, len(configs))

	for i, cfg := range configs {
		comp, err := NewComponent(cfg, theme)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar componente %d (%s): %w", i, cfg.Name, err)
		}
		components = append(components, comp)
	}

	return components, nil
}
