package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/helton/shantilly/internal/components"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
)

// LayoutModel orchestrates components in a horizontal or vertical layout.
// It manages focus navigation and responsive resizing.
type LayoutModel struct {
	title       string
	description string
	layout      string // "horizontal" or "vertical"
	components  []components.Component
	focusIndex  int
	theme       *styles.Theme
	width       int
	height      int
	quitting    bool
}

// NewLayoutModel creates a new LayoutModel from configuration.
func NewLayoutModel(cfg *config.LayoutConfig, theme *styles.Theme) (*LayoutModel, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("erro de validação da configuração do layout: %w", err)
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

	m := &LayoutModel{
		title:       cfg.Title,
		description: cfg.Description,
		layout:      cfg.Layout,
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
func (m *LayoutModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *LayoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Propagate window size to all components for responsive layout
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
func (m *LayoutModel) View() string {
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

	// Render components according to layout
	var componentsView string
	if m.layout == "horizontal" {
		componentsView = m.renderHorizontal()
	} else {
		componentsView = m.renderVertical()
	}
	sections = append(sections, componentsView)

	// Navigation help
	sections = append(sections, m.theme.Help.Render("Tab/Shift+Tab: Navegar | Esc: Sair"))

	return m.theme.Border.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

// renderHorizontal renders components in horizontal layout.
// This is the ONLY place where borders are applied to components.
func (m *LayoutModel) renderHorizontal() string {
	var views []string
	for i, comp := range m.components {
		view := comp.View()

		// Apply border based on focus state
		if i == m.focusIndex && comp.CanFocus() {
			view = m.theme.BorderActive.Render(view)
		} else {
			view = m.theme.Border.Render(view)
		}
		views = append(views, view)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, views...)
}

// renderVertical renders components in vertical layout.
// This is the ONLY place where borders are applied to components.
func (m *LayoutModel) renderVertical() string {
	var views []string
	for i, comp := range m.components {
		view := comp.View()

		// Apply border based on focus state
		if i == m.focusIndex && comp.CanFocus() {
			view = m.theme.BorderActive.Render(view)
		} else {
			view = m.theme.Border.Render(view)
		}
		views = append(views, view)
	}
	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

// focusNext moves focus to the next focusable component.
func (m *LayoutModel) focusNext() {
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
func (m *LayoutModel) focusPrev() {
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
