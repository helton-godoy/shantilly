package models

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/helton/shantilly/internal/components"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
)

// FormModel orchestrates multiple components in a form layout.
// It manages focus navigation, validation aggregation, and JSON serialization.
type FormModel struct {
	title       string
	description string
	components  []components.Component
	focusIndex  int
	theme       *styles.Theme
	width       int
	height      int
	submitted   bool
	quitting    bool
}

// NewFormModel creates a new FormModel from configuration.
func NewFormModel(cfg *config.FormConfig, theme *styles.Theme) (*FormModel, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("erro de validação da configuração do formulário: %w", err)
	}

	// Create components using factory
	comps, err := components.NewComponents(cfg.Components, theme)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar componentes: %w", err)
	}

	// Find first focusable component
	focusIndex := -1
	for i, comp := range comps {
		if comp.CanFocus() {
			focusIndex = i
			break
		}
	}

	m := &FormModel{
		title:       cfg.Title,
		description: cfg.Description,
		components:  comps,
		focusIndex:  focusIndex,
		theme:       theme,
		width:       80,
		height:      24,
	}

	// Set initial focus
	if focusIndex >= 0 {
		m.components[focusIndex].SetFocus(true)
	}

	return m, nil
}

// Init implements tea.Model.
func (m *FormModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Propagate window size to all components
		m.width = msg.Width
		m.height = msg.Height
		for i := range m.components {
			if _, err := m.components[i].Update(msg); err != nil {
				return m, func() tea.Msg {
					return fmt.Errorf("erro ao atualizar componente %d com redimensionamento: %w", i, err)
				}
			}
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "tab":
			m.focusNext()
			return m, nil

		case "shift+tab":
			m.focusPrev()
			return m, nil

		case "enter":
			// Check if form can be submitted
			if m.CanSubmit() {
				m.submitted = true
				return m, tea.Quit
			}
			// If not valid, validate all to show errors
			m.validateAll()
			return m, nil
		}
	}

	// Propagate message to focused component
	if m.focusIndex >= 0 && m.focusIndex < len(m.components) {
		var cmd tea.Cmd
		updated, cmd := m.components[m.focusIndex].Update(msg)
		if updatedModel, ok := updated.(components.Component); ok {
			m.components[m.focusIndex] = updatedModel
			return m, cmd
		} else {
			// Log error and return unchanged model
			return m, func() tea.Msg {
				return fmt.Errorf("erro ao atualizar componente %d: modelo inválido retornado", m.focusIndex)
			}
		}
	}

	return m, nil
}

// View implements tea.Model.
func (m *FormModel) View() string {
	if m.quitting {
		return ""
	}

	var sections []string

	// Title
	if m.title != "" {
		sections = append(sections, m.theme.Title.Render(m.title))
	}

	// Description
	if m.description != "" {
		sections = append(sections, m.theme.Description.Render(m.description))
	}

	// Components (without individual borders - only the form container has a border)
	for i, comp := range m.components {
		view := comp.View()

		// Apply consistent border-based focus indicator (same as LayoutModel)
		if i == m.focusIndex && comp.CanFocus() {
			view = m.theme.BorderActive.Render(view)
		} else {
			view = m.theme.Border.Render(view)
		}

		sections = append(sections, view)
	}

	// Submit help
	if m.CanSubmit() {
		sections = append(sections, m.theme.Help.Render("Pressione Enter para submeter"))
	} else {
		sections = append(sections, m.theme.Error.Render("Complete todos os campos obrigatórios"))
	}

	// Navigation help
	sections = append(sections, m.theme.Help.Render("Tab/Shift+Tab: Navegar | Esc: Sair"))

	// Don't apply border to container since individual components now have borders
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// focusNext moves focus to the next focusable component.
func (m *FormModel) focusNext() {
	if m.focusIndex >= 0 {
		m.components[m.focusIndex].SetFocus(false)
	}

	// Find next focusable component
	start := m.focusIndex + 1
	for i := 0; i < len(m.components); i++ {
		idx := (start + i) % len(m.components)
		if m.components[idx].CanFocus() {
			m.focusIndex = idx
			m.components[idx].SetFocus(true)
			return
		}
	}
}

// focusPrev moves focus to the previous focusable component.
func (m *FormModel) focusPrev() {
	if m.focusIndex >= 0 {
		m.components[m.focusIndex].SetFocus(false)
	}

	// Find previous focusable component
	start := m.focusIndex - 1
	if start < 0 {
		start = len(m.components) - 1
	}

	for i := 0; i < len(m.components); i++ {
		idx := start - i
		if idx < 0 {
			idx += len(m.components)
		}
		if m.components[idx].CanFocus() {
			m.focusIndex = idx
			m.components[idx].SetFocus(true)
			return
		}
	}
}

// CanSubmit returns true if all components are valid.
func (m *FormModel) CanSubmit() bool {
	for _, comp := range m.components {
		if !comp.IsValid() {
			return false
		}
	}
	return true
}

// validateAll validates all components to trigger error display.
func (m *FormModel) validateAll() {
	for _, comp := range m.components {
		comp.IsValid()
	}
}

// Submitted returns true if the form was successfully submitted.
func (m *FormModel) Submitted() bool {
	return m.submitted
}

// ToJSON serializes the form data to JSON.
// Returns a JSON byte array with component names as keys and values.
func (m *FormModel) ToJSON() ([]byte, error) {
	data := make(map[string]interface{})

	for _, comp := range m.components {
		data[comp.Name()] = comp.Value()
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar dados: %w", err)
	}

	return jsonData, nil
}

// ToMap returns the form data as a map for programmatic access.
func (m *FormModel) ToMap() map[string]interface{} {
	data := make(map[string]interface{})

	for _, comp := range m.components {
		data[comp.Name()] = comp.Value()
	}

	return data
}
