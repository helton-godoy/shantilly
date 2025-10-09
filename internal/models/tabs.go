package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/helton/shantilly/internal/components"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
)

// TabsModel orchestrates components in a tabbed interface.
// It manages tab navigation and component focus within tabs.
type TabsModel struct {
	name       string
	label      string
	tabs       []TabData
	activeTab  int // Index of currently active tab
	theme      *styles.Theme
	errorMsg   string
	focused    bool
	initialTab int
}

// TabData represents a single tab with its components
type TabData struct {
	Name       string
	Label      string
	Components []components.Component
}

// NewTabsModel creates a new TabsModel from configuration.
func NewTabsModel(cfg *config.TabsConfig, theme *styles.Theme) (*TabsModel, error) {
	if len(cfg.Tabs) == 0 {
		return nil, fmt.Errorf("tabs deve conter pelo menos uma aba")
	}

	tabs := make([]TabData, 0, len(cfg.Tabs))

	for _, tabCfg := range cfg.Tabs {
		// Create components for this tab using the factory
		components, err := components.NewComponents(tabCfg.Components, theme)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar componentes para aba %s: %w", tabCfg.Name, err)
		}

		tabData := TabData{
			Name:       tabCfg.Name,
			Label:      tabCfg.Label,
			Components: components,
		}

		tabs = append(tabs, tabData)
	}

	t := &TabsModel{
		name:       "tabs", // Tabs model has a fixed name
		label:      cfg.Title,
		tabs:       tabs,
		activeTab:  0,
		theme:      theme,
		initialTab: 0,
	}

	return t, nil
}

// Init implements tea.Model.
func (t *TabsModel) Init() tea.Cmd {
	// Initialize all components in the active tab
	var cmds []tea.Cmd
	for _, comp := range t.tabs[t.activeTab].Components {
		cmds = append(cmds, comp.Init())
	}
	return tea.Batch(cmds...)
}

// Update implements tea.Model.
func (t *TabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !t.focused {
		// Even if not focused, propagate messages to active tab components
		if t.activeTab >= 0 && t.activeTab < len(t.tabs) {
			for i, comp := range t.tabs[t.activeTab].Components {
				updated, cmd := comp.Update(msg)
				if updatedModel, ok := updated.(components.Component); ok {
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
			}
		case "right", "l":
			if t.activeTab < len(t.tabs)-1 {
				t.activeTab++
			}
		case "tab":
			// Tab navigation within the active tab
			t.focusNextInActiveTab()
			return t, nil
		case "shift+tab":
			// Reverse tab navigation within the active tab
			t.focusPrevInActiveTab()
			return t, nil
		}
	}

	// Propagate message to focused component in active tab
	if t.activeTab >= 0 && t.activeTab < len(t.tabs) {
		tab := &t.tabs[t.activeTab]
		for i, comp := range tab.Components {
			if comp.CanFocus() {
				updated, cmd := comp.Update(msg)
				if updatedModel, ok := updated.(components.Component); ok {
					tab.Components[i] = updatedModel
					return t, cmd
				}
			}
		}
	}

	return t, nil
}

// View implements tea.Model.
func (t *TabsModel) View() string {
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

	return b.String()
}

// renderTabHeaders renders the tab navigation header.
func (t *TabsModel) renderTabHeaders() string {
	var headers []string

	for i, tab := range t.tabs {
		var header string

		// Tab label
		if i == t.activeTab {
			header = t.theme.TabActive.Render(fmt.Sprintf(" %s ", tab.Label))
		} else {
			header = t.theme.TabInactive.Render(fmt.Sprintf(" %s ", tab.Label))
		}

		headers = append(headers, header)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, headers...)
}

// renderActiveTab renders the content of the currently active tab.
func (t *TabsModel) renderActiveTab() string {
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
func (t *TabsModel) getFocusedComponentIndex() int {
	if t.activeTab < 0 || t.activeTab >= len(t.tabs) {
		return -1
	}

	tab := &t.tabs[t.activeTab]
	for i, comp := range tab.Components {
		if comp.CanFocus() {
			return i
		}
	}
	return -1
}

// focusNextInActiveTab moves focus to the next focusable component in the active tab.
func (t *TabsModel) focusNextInActiveTab() {
	if t.activeTab < 0 || t.activeTab >= len(t.tabs) {
		return
	}

	tab := &t.tabs[t.activeTab]
	for _, comp := range tab.Components {
		if comp.CanFocus() {
			comp.SetFocus(true)
			break
		}
	}
}

// focusPrevInActiveTab moves focus to the previous focusable component in the active tab.
func (t *TabsModel) focusPrevInActiveTab() {
	if t.activeTab < 0 || t.activeTab >= len(t.tabs) {
		return
	}

	tab := &t.tabs[t.activeTab]
	for i := len(tab.Components) - 1; i >= 0; i-- {
		if tab.Components[i].CanFocus() {
			tab.Components[i].SetFocus(true)
			break
		}
	}
}
