package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ComponentType defines the type of UI component.
type ComponentType string

const (
	TypeTextInput  ComponentType = "textinput"
	TypeTextArea   ComponentType = "textarea"
	TypeCheckbox   ComponentType = "checkbox"
	TypeRadioGroup ComponentType = "radiogroup"
	TypeSlider     ComponentType = "slider"
	TypeFilePicker ComponentType = "filepicker"
	TypeText       ComponentType = "text" // Static label
)

// ComponentConfig represents the declarative configuration for a single component.
// This structure is parsed from YAML and used to initialize components.
type ComponentConfig struct {
	Type        ComponentType          `yaml:"type"`
	Name        string                 `yaml:"name"`
	Label       string                 `yaml:"label,omitempty"`
	Placeholder string                 `yaml:"placeholder,omitempty"`
	Default     interface{}            `yaml:"default,omitempty"`
	Required    bool                   `yaml:"required,omitempty"`
	Help        string                 `yaml:"help,omitempty"`
	Options     map[string]interface{} `yaml:"options,omitempty"`
}

// Validate performs validation on the ComponentConfig.
// Returns an error if the configuration is invalid.
func (c *ComponentConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("nome do componente é obrigatório")
	}

	validTypes := []ComponentType{
		TypeTextInput, TypeTextArea, TypeCheckbox,
		TypeRadioGroup, TypeSlider, TypeFilePicker, TypeText,
	}

	valid := false
	for _, t := range validTypes {
		if c.Type == t {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("tipo de componente inválido: %s", c.Type)
	}

	return nil
}

// FormConfig represents the complete form configuration with multiple components.
type FormConfig struct {
	Title       string            `yaml:"title,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Components  []ComponentConfig `yaml:"components"`
}

// Validate performs validation on the FormConfig.
func (f *FormConfig) Validate() error {
	if len(f.Components) == 0 {
		return fmt.Errorf("a configuração deve conter pelo menos um componente")
	}

	// Validate each component
	for i, comp := range f.Components {
		if err := comp.Validate(); err != nil {
			return fmt.Errorf("erro no componente %d: %w", i, err)
		}
	}

	// Check for duplicate component names
	names := make(map[string]bool)
	for _, comp := range f.Components {
		if names[comp.Name] {
			return fmt.Errorf("nome de componente duplicado: %s", comp.Name)
		}
		names[comp.Name] = true
	}

	return nil
}

// LayoutConfig represents a layout configuration with positioned components.
type LayoutConfig struct {
	Title       string            `yaml:"title,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Layout      string            `yaml:"layout"` // "horizontal" or "vertical"
	Components  []ComponentConfig `yaml:"components"`
}

// Validate performs validation on the LayoutConfig.
func (l *LayoutConfig) Validate() error {
	if l.Layout != "horizontal" && l.Layout != "vertical" {
		return fmt.Errorf("layout deve ser 'horizontal' ou 'vertical', recebido: %s", l.Layout)
	}

	if len(l.Components) == 0 {
		return fmt.Errorf("a configuração deve conter pelo menos um componente")
	}

	for i, comp := range l.Components {
		if err := comp.Validate(); err != nil {
			return fmt.Errorf("erro no componente %d: %w", i, err)
		}
	}

	return nil
}

// MenuConfig represents a menu/list selection configuration.
type MenuConfig struct {
	Title       string   `yaml:"title,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Items       []string `yaml:"items"`
	MultiSelect bool     `yaml:"multi_select,omitempty"`
}

// Validate performs validation on the MenuConfig.
func (m *MenuConfig) Validate() error {
	if len(m.Items) == 0 {
		return fmt.Errorf("o menu deve conter pelo menos um item")
	}
	return nil
}

// TabConfig represents a single tab configuration.
type TabConfig struct {
	Name       string            `yaml:"name"`
	Label      string            `yaml:"label"`
	Components []ComponentConfig `yaml:"components"`
}

// TabsConfig represents a tabs configuration with multiple tabs.
type TabsConfig struct {
	Title string      `yaml:"title,omitempty"`
	Tabs  []TabConfig `yaml:"tabs"`
}

// Validate performs validation on the TabsConfig.
func (t *TabsConfig) Validate() error {
	if len(t.Tabs) == 0 {
		return fmt.Errorf("a configuração deve conter pelo menos uma aba")
	}

	for i, tab := range t.Tabs {
		if tab.Name == "" {
			return fmt.Errorf("aba %d: nome é obrigatório", i)
		}
		if tab.Label == "" {
			return fmt.Errorf("aba %d: label é obrigatório", i)
		}

		for j, comp := range tab.Components {
			if err := comp.Validate(); err != nil {
				return fmt.Errorf("aba %d, componente %d: %w", i, j, err)
			}
		}
	}

	return nil
}

// LoadFormConfig loads and validates a FormConfig from a YAML file.
// Returns an error with context if loading or validation fails.
func LoadFormConfig(filePath string) (*FormConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo de configuração: %w", err)
	}

	var config FormConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao analisar o YAML de configuração: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("erro de validação da configuração: %w", err)
	}

	return &config, nil
}

// LoadLayoutConfig loads and validates a LayoutConfig from a YAML file.
func LoadLayoutConfig(filePath string) (*LayoutConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo de configuração: %w", err)
	}

	var config LayoutConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao analisar o YAML de configuração: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("erro de validação da configuração: %w", err)
	}

	return &config, nil
}

// LoadMenuConfig loads and validates a MenuConfig from a YAML file.
func LoadMenuConfig(filePath string) (*MenuConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo de configuração: %w", err)
	}

	var config MenuConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao analisar o YAML de configuração: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("erro de validação da configuração: %w", err)
	}

	return &config, nil
}

// LoadTabsConfig loads and validates a TabsConfig from a YAML file.
func LoadTabsConfig(filePath string) (*TabsConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo de configuração: %w", err)
	}

	var config TabsConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao analisar o YAML de configuração: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("erro de validação da configuração: %w", err)
	}

	return &config, nil
}
