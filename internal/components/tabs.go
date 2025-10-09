package components

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
	"gopkg.in/yaml.v3"
)

// TabItem represents a single tab with its components.
type TabItem struct {
	Name       string
	Label      string
	Components []Component
}

// Tabs implements a tabs component that manages multiple tabs with their own components.
// This component orchestrates multiple components within tabs and handles tab navigation.
type Tabs struct {
	name       string
	label      string
	required   bool
	help       string
	tabs       []TabItem
	activeTab  int // Index of currently active tab
	theme      *styles.Theme
	errorMsg   string
	focused    bool
	initialTab int
}

// NewTabs creates a new Tabs component from configuration.
// This is different from other components as it needs to create child components for each tab.
func NewTabs(cfg config.TabsConfig, theme *styles.Theme) (*Tabs, error) {
	if len(cfg.Tabs) == 0 {
		return nil, fmt.Errorf("tabs deve conter pelo menos uma aba")
	}

	tabs := make([]TabItem, 0, len(cfg.Tabs))

	for _, tabCfg := range cfg.Tabs {
		// Validate required tab fields
		if tabCfg.Name == "" {
			return nil, fmt.Errorf("nome da aba é obrigatório")
		}
		if tabCfg.Label == "" {
			return nil, fmt.Errorf("label da aba é obrigatório")
		}

		// Create components for this tab using the factory
		components, err := NewComponents(tabCfg.Components, theme)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar componentes para aba %s: %w", tabCfg.Name, err)
		}

		tabItem := TabItem{
			Name:       tabCfg.Name,
			Label:      tabCfg.Label,
			Components: components,
		}

		tabs = append(tabs, tabItem)
	}

	t := &Tabs{
		name:       "tabs", // Tabs component has a fixed name
		label:      cfg.Title,
		required:   false, // Tabs itself doesn't have required validation
		help:       "",
		tabs:       tabs,
		activeTab:  0,
		theme:      theme,
		initialTab: 0,
	}

	return t, nil
}

// Init implements tea.Model.
func (t *Tabs) Init() tea.Cmd {
	// Initialize all components in the active tab
	var cmds []tea.Cmd
	for _, comp := range t.tabs[t.activeTab].Components {
		cmds = append(cmds, comp.Init())
	}
	return tea.Batch(cmds...)
}

// Update implements tea.Model.
func (t *Tabs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !t.focused {
		// Even if not focused, propagate messages to active tab components
		// for non-focus messages like WindowSizeMsg
		if t.activeTab >= 0 && t.activeTab < len(t.tabs) {
			for i, comp := range t.tabs[t.activeTab].Components {
				updated, cmd := comp.Update(msg)
				if updatedModel, ok := updated.(Component); ok {
					t.tabs[t.activeTab].Components[i] = updatedModel
					return t, cmd
				}
			}
		}
		return t, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if t.activeTab > 0 {
				t.activeTab--
				t.updateActiveTab()
			}
		case "right", "l":
			if t.activeTab < len(t.tabs)-1 {
				t.activeTab++
				t.updateActiveTab()
			}
		case "tab":
			// Tab navigation within the active tab
			t.focusNextInActiveTab()
			return t, nil
		case "shift+tab":
			// Reverse tab navigation within the active tab
			t.focusPrevInActiveTab()
			return t, nil
		case "ctrl+tab":
			// Ctrl+Tab: próxima aba
			t.nextTab()
			return t, nil
		case "ctrl+shift+tab":
			// Ctrl+Shift+Tab: aba anterior
			t.prevTab()
			return t, nil
		case "ctrl+1", "ctrl+2", "ctrl+3", "ctrl+4", "ctrl+5", "ctrl+6", "ctrl+7", "ctrl+8", "ctrl+9":
			// Ctrl+[número]: ir para aba específica
			if num := int(msg.String()[5] - '0'); num > 0 && num <= len(t.tabs) {
				t.activeTab = num - 1
				t.updateActiveTab()
			}
			return t, nil
		}
	}

	// Propagate message to focused component in active tab
	if t.activeTab >= 0 && t.activeTab < len(t.tabs) {
		tab := &t.tabs[t.activeTab]
		for i, comp := range tab.Components {
			if comp.CanFocus() {
				updated, cmd := comp.Update(msg)
				if updatedModel, ok := updated.(Component); ok {
					tab.Components[i] = updatedModel
					return t, cmd
				}
			}
		}
	}

	return t, nil
}

// View implements tea.Model.
func (t *Tabs) View() string {
	var b strings.Builder

	// Render tab headers
	tabHeaders := t.renderTabHeaders()
	b.WriteString(tabHeaders)
	b.WriteString("\n")

	// Render active tab content
	if t.activeTab >= 0 && t.activeTab < len(t.tabs) {
		activeTabContent := t.renderActiveTab()
		b.WriteString(activeTabContent)
	}

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

// renderTabHeaders renders the tab navigation header.
func (t *Tabs) renderTabHeaders() string {
	return t.renderTabHeadersWithErrors()
}

// renderActiveTab renders the content of the currently active tab.
func (t *Tabs) renderActiveTab() string {
	if t.activeTab < 0 || t.activeTab >= len(t.tabs) {
		return ""
	}

	tab := &t.tabs[t.activeTab]
	var componentsView []string

	for i, comp := range tab.Components {
		view := comp.View()

		// Apply border based on focus state (similar to other models)
		if i == t.getFocusedComponentIndex() && comp.CanFocus() {
			view = t.theme.BorderActive.Render(view)
		} else {
			view = t.theme.Border.Render(view)
		}

		componentsView = append(componentsView, view)
	}

	return lipgloss.JoinVertical(lipgloss.Left, componentsView...)
}

// getFocusedComponentIndex returns the index of the currently focused component in the active tab.
func (t *Tabs) getFocusedComponentIndex() int {
	if t.activeTab < 0 || t.activeTab >= len(t.tabs) {
		return -1
	}

	tab := &t.tabs[t.activeTab]
	for i, comp := range tab.Components {
		// This is a simplified approach - in a real implementation,
		// we'd need to track which component in which tab is focused
		if comp.CanFocus() {
			return i
		}
	}
	return -1
}

// focusNextInActiveTab moves focus to the next focusable component in the active tab.
func (t *Tabs) focusNextInActiveTab() {
	if t.activeTab < 0 || t.activeTab >= len(t.tabs) {
		return
	}

	tab := &t.tabs[t.activeTab]
	for _, comp := range tab.Components {
		if comp.CanFocus() {
			comp.SetFocus(true)
			// In a real implementation, we'd track this index
			break
		}
	}
}

// focusPrevInActiveTab moves focus to the previous focusable component in the active tab.
func (t *Tabs) focusPrevInActiveTab() {
	if t.activeTab < 0 || t.activeTab >= len(t.tabs) {
		return
	}

	tab := &t.tabs[t.activeTab]
	for i := len(tab.Components) - 1; i >= 0; i-- {
		if tab.Components[i].CanFocus() {
			tab.Components[i].SetFocus(true)
			// In a real implementation, we'd track this index
			break
		}
	}
}

// SetTheme implements Component.
func (t *Tabs) SetTheme(theme *styles.Theme) {
	t.theme = theme
	// Propagar mudança de tema para todos os componentes das abas
	for _, tab := range t.tabs {
		for _, comp := range tab.Components {
			comp.SetTheme(theme)
		}
	}
}

// Name implements Component.
func (t *Tabs) Name() string {
	return t.name
}

// CanFocus implements Component.
func (t *Tabs) CanFocus() bool {
	return true
}

// SetFocus implements Component.
func (t *Tabs) SetFocus(focused bool) {
	t.focused = focused
}

// IsValid implements Component.
func (t *Tabs) IsValid() bool {
	errors := t.ValidateWithContext(ValidationContext{
		ComponentValues: make(map[string]interface{}),
		GlobalConfig:    make(map[string]interface{}),
		ExternalData:    make(map[string]interface{}),
	})

	// Se há erros de validação, definir a primeira mensagem de erro
	if len(errors) > 0 {
		t.errorMsg = errors[0].Message
		return false
	}

	t.errorMsg = ""
	return true
}

// GetError implements Component.
func (t *Tabs) GetError() string {
	return t.errorMsg
}

// ValidateWithContext implements Component.
func (t *Tabs) ValidateWithContext(context ValidationContext) []ValidationError {
	var errors []ValidationError

	// Validar todas as abas e seus componentes
	for i, tab := range t.tabs {
		tabErrors := t.getTabErrors(i)

		// Adicionar erros de componentes individuais
		errors = append(errors, tabErrors...)

		// Validação específica da aba
		tabStatus := t.getTabValidationStatus(i)
		if !tabStatus.IsValid {
			errors = append(errors, ValidationError{
				Code:     "TAB_VALIDATION_FAILED",
				Message:  fmt.Sprintf("Aba '%s' contém erros de validação", tab.Label),
				Field:    fmt.Sprintf("tabs.%d", i),
				Severity: "error",
				Context: map[string]interface{}{
					"component":        "Tabs",
					"tab_name":         tab.Name,
					"tab_label":        tab.Label,
					"tab_index":        i,
					"valid_components": tabStatus.ValidComponents,
					"total_components": tabStatus.TotalComponents,
				},
			})
		}
	}

	// Validação cruzada entre abas (exemplo)
	if len(t.tabs) > 1 {
		// Verificar se há conflitos entre abas
		for i := 0; i < len(t.tabs); i++ {
			for j := i + 1; j < len(t.tabs); j++ {
				crossErrors := t.validateCrossTabDependencies(i, j, context)
				errors = append(errors, crossErrors...)
			}
		}
	}

	return errors
}

// validateCrossTabDependencies valida dependências entre abas
func (t *Tabs) validateCrossTabDependencies(tabIndex1, tabIndex2 int, context ValidationContext) []ValidationError {
	var errors []ValidationError

	// Exemplo de validação cruzada: verificar se há componentes relacionados
	tab1 := &t.tabs[tabIndex1]
	tab2 := &t.tabs[tabIndex2]

	// Verificar se há componentes com mesmo nome em abas diferentes (conflito potencial)
	for _, comp1 := range tab1.Components {
		for _, comp2 := range tab2.Components {
			if comp1.Name() == comp2.Name() && comp1.Value() != comp2.Value() {
				// Encontrou componentes com mesmo nome mas valores diferentes
				errors = append(errors, ValidationError{
					Code:     "CROSS_TAB_CONFLICT",
					Message:  fmt.Sprintf("Conflito entre abas: componente '%s' tem valores diferentes", comp1.Name()),
					Field:    fmt.Sprintf("tabs.%d.%s", tabIndex1, comp1.Name()),
					Severity: "warning",
					Context: map[string]interface{}{
						"component":      "Tabs",
						"tab1_name":      tab1.Name,
						"tab2_name":      tab2.Name,
						"component_name": comp1.Name(),
						"value1":         comp1.Value(),
						"value2":         comp2.Value(),
					},
				})
			}
		}
	}

	return errors
}

// SetError implements Component.
func (t *Tabs) SetError(msg string) {
	t.errorMsg = msg
}

// Value implements Component.
// Returns a map with tab names as keys and component values as values.
func (t *Tabs) Value() interface{} {
	result := make(map[string]interface{})

	for _, tab := range t.tabs {
		tabData := make(map[string]interface{})
		for _, comp := range tab.Components {
			tabData[comp.Name()] = comp.Value()
		}
		result[tab.Name] = tabData
	}

	return result
}

// SetValue implements Component.
func (t *Tabs) SetValue(value interface{}) error {
	// This is a simplified implementation
	// In a real implementation, we'd need to properly distribute values to components
	_, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("valor inválido: esperado map[string]interface{}, recebido %T", value)
	}

	// TODO: Implement proper value setting for tabs
	return nil
}

// Reset implements Component.
func (t *Tabs) Reset() {
	t.activeTab = t.initialTab
	t.errorMsg = ""
	t.focused = false

	// Reset all components in all tabs
	for _, tab := range t.tabs {
		for _, comp := range tab.Components {
			comp.Reset()
		}
	}
}

// Métodos avançados de navegação para Tabs

// updateActiveTab atualiza o estado quando a aba ativa muda
func (t *Tabs) updateActiveTab() {
	// Inicializar componentes da nova aba ativa
	if t.activeTab >= 0 && t.activeTab < len(t.tabs) {
		var cmds []tea.Cmd
		for _, comp := range t.tabs[t.activeTab].Components {
			cmds = append(cmds, comp.Init())
		}
		// Retorna um batch command se necessário
		if len(cmds) > 0 {
			tea.Batch(cmds...)
		}
	}
}

// nextTab vai para a próxima aba
func (t *Tabs) nextTab() {
	if t.activeTab < len(t.tabs)-1 {
		t.activeTab++
	} else {
		t.activeTab = 0 // Volta para a primeira aba
	}
	t.updateActiveTab()
}

// prevTab vai para a aba anterior
func (t *Tabs) prevTab() {
	if t.activeTab > 0 {
		t.activeTab--
	} else {
		t.activeTab = len(t.tabs) - 1 // Vai para a última aba
	}
	t.updateActiveTab()
}

// getTabErrors retorna erros da aba especificada
func (t *Tabs) getTabErrors(tabIndex int) []ValidationError {
	if tabIndex < 0 || tabIndex >= len(t.tabs) {
		return []ValidationError{}
	}

	var errors []ValidationError
	tab := &t.tabs[tabIndex]

	for _, comp := range tab.Components {
		compErrors := comp.ValidateWithContext(ValidationContext{
			ComponentValues: make(map[string]interface{}),
			GlobalConfig:    make(map[string]interface{}),
			ExternalData:    make(map[string]interface{}),
		})

		errors = append(errors, compErrors...)
	}

	return errors
}

// renderTabHeadersWithErrors renderiza cabeçalhos com indicadores de erro
func (t *Tabs) renderTabHeadersWithErrors() string {
	var headers []string

	for i, tab := range t.tabs {
		var header string
		tabErrors := t.getTabErrors(i)

		// Tab label
		if i == t.activeTab {
			header = t.theme.TabActive.Render(fmt.Sprintf(" %s ", tab.Label))
		} else {
			header = t.theme.TabInactive.Render(fmt.Sprintf(" %s ", tab.Label))
		}

		// Adicionar indicador de erro se houver
		if len(tabErrors) > 0 {
			errorIndicator := t.theme.Error.Render("⚠")
			header = header + " " + errorIndicator
		}

		headers = append(headers, header)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, headers...)
}

// getTabValidationStatus retorna o status de validação de uma aba
func (t *Tabs) getTabValidationStatus(tabIndex int) TabValidationStatus {
	if tabIndex < 0 || tabIndex >= len(t.tabs) {
		return TabValidationStatus{}
	}

	tab := &t.tabs[tabIndex]
	validComponents := 0
	totalComponents := len(tab.Components)

	for _, comp := range tab.Components {
		if comp.IsValid() {
			validComponents++
		}
	}

	return TabValidationStatus{
		TabIndex:        tabIndex,
		TabName:         tab.Name,
		TabLabel:        tab.Label,
		ValidComponents: validComponents,
		TotalComponents: totalComponents,
		IsValid:         validComponents == totalComponents,
		Errors:          t.getTabErrors(tabIndex),
	}
}

// TabValidationStatus representa o status de validação de uma aba
type TabValidationStatus struct {
	TabIndex        int               `json:"tab_index"`
	TabName         string            `json:"tab_name"`
	TabLabel        string            `json:"tab_label"`
	ValidComponents int               `json:"valid_components"`
	TotalComponents int               `json:"total_components"`
	IsValid         bool              `json:"is_valid"`
	Errors          []ValidationError `json:"errors"`
}

// GetMetadata implements Component.
func (t *Tabs) GetMetadata() ComponentMetadata {
	return ComponentMetadata{
		Version:      "1.0.0",
		Author:       "Shantilly Team",
		Description:  "Componente de abas para organização de múltiplos componentes",
		Dependencies: []string{},
		Examples: []ComponentExample{
			{
				Name:        "Abas básicas",
				Description: "Exemplo de configuração com múltiplas abas",
				Config: map[string]interface{}{
					"title": "Configurações",
					"tabs": []map[string]interface{}{
						{
							"name":  "geral",
							"label": "Geral",
							"components": []map[string]interface{}{
								{"name": "host", "type": "textinput", "label": "Host"},
							},
						},
					},
				},
			},
		},
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"title": map[string]interface{}{"type": "string"},
				"tabs": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name":       map[string]interface{}{"type": "string"},
							"label":      map[string]interface{}{"type": "string"},
							"components": map[string]interface{}{"type": "array"},
						},
						"required": []string{"name", "label"},
					},
				},
			},
			"required": []string{"tabs"},
		},
	}
}

// ExportToFormat implements Component.
func (t *Tabs) ExportToFormat(format ExportFormat) ([]byte, error) {
	data := t.Value()
	switch format {
	case FormatJSON:
		return json.Marshal(data)
	case FormatYAML:
		return yaml.Marshal(data)
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

// ImportFromFormat implements Component.
func (t *Tabs) ImportFromFormat(format ExportFormat, data []byte) error {
	var value interface{}
	switch format {
	case FormatJSON:
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
	case FormatYAML:
		if err := yaml.Unmarshal(data, &value); err != nil {
			return err
		}
	default:
		return fmt.Errorf("formato não suportado: %s", format)
	}

	return t.SetValue(value)
}

// GetDependencies implements Component.
func (t *Tabs) GetDependencies() []string {
	var deps []string
	for _, tab := range t.tabs {
		for _, comp := range tab.Components {
			compDeps := comp.GetDependencies()
			deps = append(deps, compDeps...)
		}
	}
	return deps
}
