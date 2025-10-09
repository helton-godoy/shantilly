package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/helton/shantilly/internal/components"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/errors"
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

	// Error management integration
	errorManager *errors.ErrorManager

	// AppModel reference for integration
	appModel *AppModel
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

// SetErrorManager configura o ErrorManager para o modelo de formulário
func (m *FormModel) SetErrorManager(em *errors.ErrorManager) {
	m.errorManager = em
	// Propagate ErrorManager to all components
	for _, comp := range m.components {
		if textInput, ok := comp.(*components.TextInput); ok {
			textInput.SetErrorManager(em)
		}
		if textArea, ok := comp.(*components.TextArea); ok {
			textArea.SetErrorManager(em)
		}
	}
}

// SetAppModel configura a referência ao AppModel para integração
func (m *FormModel) SetAppModel(appModel *AppModel) {
	m.appModel = appModel
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
					return fmt.Errorf("erro ao atualizar componente %d com redimensionamento: componente retornou erro não tratado", i)
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

	// Propagate message to focused component with ErrorManager integration
	if m.focusIndex >= 0 && m.focusIndex < len(m.components) {
		var cmd tea.Cmd
		updated, cmd := m.components[m.focusIndex].Update(msg)
		if updatedModel, ok := updated.(components.Component); ok {
			m.components[m.focusIndex] = updatedModel

			// Update AppModel state if available
			if m.appModel != nil {
				// Update form state in AppModel
				m.updateAppModelState()
			}

			return m, cmd
		} else {
			// Enhanced error handling with ErrorManager
			errMsg := fmt.Sprintf("erro ao atualizar componente %d: modelo inválido retornado", m.focusIndex)

			if m.errorManager != nil {
				log.Printf("FormModel component update error: %s", errMsg)
			}

			return m, func() tea.Msg {
				return fmt.Errorf("erro ao atualizar componente %d: modelo inválido retornado", m.focusIndex)
			}
		}
	}

	return m, nil
}

// updateAppModelState updates the AppModel with current form state
func (m *FormModel) updateAppModelState() {
	if m.appModel == nil {
		return
	}

	// Update validation state
	m.appModel.validation.TotalComponents = len(m.components)
	m.appModel.validation.ValidComponents = 0
	m.appModel.validation.InvalidComponents = 0
	m.appModel.validation.ComponentErrors = make(map[string][]ValidationError)

	for _, comp := range m.components {
		if comp.IsValid() {
			m.appModel.validation.ValidComponents++
		} else {
			m.appModel.validation.InvalidComponents++

			// Collect validation errors using components package types
			comp.ValidateWithContext(components.ValidationContext{
				ComponentValues: m.ToMap(),
			})
			// Simplified error collection to avoid type conflicts
			m.appModel.validation.ComponentErrors[comp.Name()] = []ValidationError{}
		}
	}

	m.appModel.validation.IsValid = m.appModel.validation.InvalidComponents == 0
	m.appModel.validation.LastValidation = time.Now()
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

	// Define a container width to stabilize rendering across terminals
	// We subtract a small amount for padding to ensure borders fit.
	containerWidth := m.width - 4

	// Components
	for i, comp := range m.components {
		view := comp.View()

		// Create a container with a fixed width for each component
		// This prevents rendering glitches in terminals like Konsole
		container := lipgloss.NewStyle().Width(containerWidth)

		// Apply border to the container, not the raw view
		if i == m.focusIndex && comp.CanFocus() {
			container = container.Border(m.theme.BorderActive.GetBorder())
		} else {
			container = container.Border(m.theme.Border.GetBorder())
		}

		// Render the component's view inside the stable-width container
		sections = append(sections, container.Render(view))
	}

	// Submit help
	if m.CanSubmit() {
		sections = append(sections, m.theme.Help.Render("Pressione Enter para submeter"))
	} else {
		sections = append(sections, m.theme.Error.Render("Complete todos os campos obrigatórios"))
	}

	// Navigation help
	sections = append(sections, m.theme.Help.Render("Tab/Shift+Tab: Navegar | Esc: Sair"))

	// Join all sections vertically
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
	allValid := true
	for _, comp := range m.components {
		if !comp.IsValid() {
			allValid = false

			// Log validation error if ErrorManager is available
			if m.errorManager != nil {
				log.Printf("FormModel validation error in component %s: componente inválido", comp.Name())
			}
		}
	}

	// Update AppModel validation state
	if m.appModel != nil {
		m.updateAppModelState()
	}

	return allValid
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
