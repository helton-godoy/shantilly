package components

import (
	"fmt"
	"log"

	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/errors"
	"github.com/helton/shantilly/internal/styles"
)

// Global factory error manager for enhanced error handling
var factoryErrorManager *errors.ErrorManager

// SetFactoryErrorManager sets the global error manager for the factory
func SetFactoryErrorManager(em *errors.ErrorManager) {
	factoryErrorManager = em
}

// NewComponent creates a component from configuration using the factory pattern.
// This ensures that all components are created consistently and with proper validation.
func NewComponent(cfg config.ComponentConfig, theme *styles.Theme) (Component, error) {
	// Enhanced configuration validation with ErrorManager integration
	if err := cfg.Validate(); err != nil {
		if factoryErrorManager != nil {
			log.Printf("Factory component validation error for %s: %v", cfg.Name, err)
		}
		return nil, fmt.Errorf("erro de validação da configuração: %w", err)
	}

	// Component creation with enhanced error handling
	var component Component
	var err error

	switch cfg.Type {
	case config.TypeTextInput:
		component, err = NewTextInput(cfg, theme)
	case config.TypeTextArea:
		component, err = NewTextArea(cfg, theme)
	case config.TypeCheckbox:
		component, err = NewCheckbox(cfg, theme)
	case config.TypeRadioGroup:
		component, err = NewRadioGroup(cfg, theme)
	case config.TypeSlider:
		component, err = NewSlider(cfg, theme)
	case config.TypeText:
		component, err = NewTextLabel(cfg, theme)
	case config.TypeFilePicker:
		component, err = NewFilePicker(cfg, theme)
	default:
		err = fmt.Errorf("tipo de componente não suportado: %s", cfg.Type)
	}

	// Enhanced error handling for component creation
	if err != nil {
		if factoryErrorManager != nil {
			log.Printf("Factory component creation error for %s (%s): %v", cfg.Name, cfg.Type, err)
		}
		return nil, fmt.Errorf("erro ao criar componente %s: %w", cfg.Name, err)
	}

	// Set ErrorManager on components that support it
	if factoryErrorManager != nil {
		if textInput, ok := component.(*TextInput); ok {
			textInput.SetErrorManager(factoryErrorManager)
		}
		if textArea, ok := component.(*TextArea); ok {
			textArea.SetErrorManager(factoryErrorManager)
		}
	}

	return component, nil
}

// NewComponents creates multiple components from a slice of configurations.
// Returns an error if any component fails to create, with enhanced error reporting.
func NewComponents(configs []config.ComponentConfig, theme *styles.Theme) ([]Component, error) {
	if len(configs) == 0 {
		return []Component{}, nil
	}

	components := make([]Component, 0, len(configs))

	for i, cfg := range configs {
		comp, err := NewComponent(cfg, theme)
		if err != nil {
			// Enhanced error reporting with context
			if factoryErrorManager != nil {
				log.Printf("Factory batch creation error at index %d: %v", i, err)
			}
			return nil, fmt.Errorf("erro ao criar componente %d (%s): %w", i, cfg.Name, err)
		}
		components = append(components, comp)
	}

	// Success logging
	if factoryErrorManager != nil {
		log.Printf("Factory successfully created %d components", len(components))
	}

	return components, nil
}
