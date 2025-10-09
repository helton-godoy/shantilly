package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/helton/shantilly/internal/components"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
)

// ViewType defines the different views/modes of the application
type ViewType int

const (
	FormView ViewType = iota
	LayoutView
	TabsView
	FilePickerView
	MenuView
)

// String returns the string representation of ViewType
func (vt ViewType) String() string {
	switch vt {
	case FormView:
		return "form"
	case LayoutView:
		return "layout"
	case TabsView:
		return "tabs"
	case FilePickerView:
		return "filepicker"
	case MenuView:
		return "menu"
	default:
		return "unknown"
	}
}

// AppMetadata contains comprehensive information about the application instance
type AppMetadata struct {
	Version      string    `json:"version"`
	BuildTime    string    `json:"build_time"`
	GitCommit    string    `json:"git_commit"`
	GoVersion    string    `json:"go_version"`
	Architecture string    `json:"architecture"`
	Platform     string    `json:"platform"`
	StartTime    time.Time `json:"start_time"`
	SessionID    string    `json:"session_id"`
}

// PerformanceMetrics tracks application performance indicators
type PerformanceMetrics struct {
	RenderCount     int           `json:"render_count"`
	TotalRenderTime time.Duration `json:"total_render_time"`
	AvgRenderTime   time.Duration `json:"avg_render_time"`
	LastUpdate      time.Time     `json:"last_update"`
	MemoryUsage     uint64        `json:"memory_usage"`
	GCCycles        uint32        `json:"gc_cycles"`
}

// ValidationState tracks the validation status across all components
type ValidationState struct {
	IsValid            bool                         `json:"is_valid"`
	TotalComponents    int                          `json:"total_components"`
	ValidComponents    int                          `json:"valid_components"`
	InvalidComponents  int                          `json:"invalid_components"`
	ComponentErrors    map[string][]ValidationError `json:"component_errors"`
	LastValidation     time.Time                    `json:"last_validation"`
	ValidationDuration time.Duration                `json:"validation_duration"`
}

// ValidationError represents a validation error with detailed information
type ValidationError struct {
	Code     string                 `json:"code"`
	Message  string                 `json:"message"`
	Field    string                 `json:"field"`
	Severity string                 `json:"severity"`
	Context  map[string]interface{} `json:"context"`
}

// AppModel is the central state management model for the entire application.
// It manages view transitions, global state, error handling, and coordinates
// between different orchestration models (FormModel, LayoutModel, TabsModel).
type AppModel struct {
	// Current application state
	currentView  ViewType `json:"current_view"`
	previousView ViewType `json:"previous_view"`

	// Active model instance (can be FormModel, LayoutModel, TabsModel, etc.)
	activeModel tea.Model `json:"-"` // Not serialized

	// Global application state
	config      *config.Config     `json:"config"`
	theme       *styles.Theme      `json:"theme"`
	metadata    AppMetadata        `json:"metadata"`
	performance PerformanceMetrics `json:"performance"`
	validation  ValidationState    `json:"validation"`

	// Error management
	errors      []AppError `json:"errors"`
	lastErrorID int        `json:"last_error_id"`

	// Component registry for dependency injection
	components map[string]components.Component `json:"-"` // Not serialized

	// Navigation state
	navigationHistory []ViewType `json:"navigation_history"`
	navigationIndex   int        `json:"navigation_index"`

	// Application lifecycle
	started  bool `json:"started"`
	quitting bool `json:"quitting"`
	debug    bool `json:"debug"`

	// Window and terminal state
	width         int  `json:"width"`
	height        int  `json:"height"`
	terminalReady bool `json:"terminal_ready"`
}

// AppError represents a structured application error with full context
type AppError struct {
	ID         string                 `json:"id"`
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Component  string                 `json:"component"`
	Severity   ErrorSeverity          `json:"severity"`
	Context    map[string]interface{} `json:"context"`
	StackTrace string                 `json:"stack_trace,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Resolved   bool                   `json:"resolved"`
}

// ErrorCode defines standardized error codes for the application
type ErrorCode int

const (
	ErrComponentNotFound ErrorCode = iota + 1000
	ErrValidationFailed
	ErrConfigInvalid
	ErrThemeLoadFailed
	ErrComponentCreationFailed
	ErrStateManagementFailed
	ErrSerializationFailed
	ErrFileOperationFailed
	ErrNetworkOperationFailed
	ErrPermissionDenied
	ErrResourceNotFound
	ErrTimeout
	ErrConcurrencyIssue
	ErrMemoryAllocationFailed
	ErrInvalidViewTransition
	ErrComponentDependencyFailed
)

// ErrorSeverity defines the severity levels for errors
type ErrorSeverity int

const (
	SeverityInfo ErrorSeverity = iota
	SeverityWarning
	SeverityError
	SeverityCritical
	SeverityFatal
)

// NewAppModel creates a new AppModel with the specified configuration
func NewAppModel(cfg *config.Config, theme *styles.Theme) (*AppModel, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("erro de validação da configuração: %w", err)
	}

	now := time.Now()

	app := &AppModel{
		currentView:  FormView, // Default to form view
		previousView: FormView,
		config:       cfg,
		theme:        theme,
		components:   make(map[string]components.Component),
		errors:       make([]AppError, 0),
		metadata: AppMetadata{
			Version:   cfg.Global.Version,
			StartTime: now,
			SessionID: generateSessionID(),
		},
		performance: PerformanceMetrics{
			LastUpdate: now,
		},
		validation: ValidationState{
			IsValid:         false,
			ComponentErrors: make(map[string][]ValidationError),
		},
		navigationHistory: make([]ViewType, 0),
		navigationIndex:   -1,
		width:             80,
		height:            24,
		debug:             cfg.Global.Debug,
	}

	// Initialize the first view based on configuration
	if err := app.initializeView(); err != nil {
		return nil, fmt.Errorf("erro ao inicializar visão inicial: %w", err)
	}

	return app, nil
}

// generateSessionID generates a unique session identifier
func generateSessionID() string {
	return fmt.Sprintf("shantilly_%d", time.Now().UnixNano())
}

// initializeView initializes the current view based on the application state
func (app *AppModel) initializeView() error {
	switch app.currentView {
	case FormView:
		if len(app.config.Forms) > 0 {
			formModel, err := NewFormModel(&app.config.Forms[0], app.theme)
			if err != nil {
				return fmt.Errorf("erro ao criar modelo de formulário: %w", err)
			}
			app.activeModel = formModel
		}

	case LayoutView:
		if len(app.config.Layouts) > 0 {
			layoutModel, err := NewLayoutModel(&app.config.Layouts[0], app.theme)
			if err != nil {
				return fmt.Errorf("erro ao criar modelo de layout: %w", err)
			}
			app.activeModel = layoutModel
		}

	case TabsView:
		if len(app.config.Tabs) > 0 {
			tabsModel, err := NewTabsModel(&app.config.Tabs[0], app.theme)
			if err != nil {
				return fmt.Errorf("erro ao criar modelo de abas: %w", err)
			}
			app.activeModel = tabsModel
		}

	default:
		return fmt.Errorf("tipo de visão não suportado: %s", app.currentView.String())
	}

	return nil
}

// Init implements tea.Model
func (app *AppModel) Init() tea.Cmd {
	if app.activeModel != nil {
		return app.activeModel.Init()
	}
	return nil
}

// Update implements tea.Model
func (app *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	startTime := time.Now()

	// Handle global messages first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		app.width = msg.Width
		app.height = msg.Height
		app.terminalReady = true

		// Propagate to active model
		if app.activeModel != nil {
			return app.activeModel.Update(msg)
		}
		return app, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			app.quitting = true
			return app, tea.Quit

		case "F1":
			// Debug information toggle
			app.debug = !app.debug
			return app, nil

		case "F2":
			// View navigation (for testing)
			return app, app.navigateToNextView()

		case "F12":
			// Show application statistics
			app.logAppStats()
			return app, nil
		}
	}

	// Update active model
	if app.activeModel != nil {
		updatedModel, cmd := app.activeModel.Update(msg)
		if updatedModel != nil {
			app.activeModel = updatedModel
		}

		// Update performance metrics
		updateDuration := time.Since(startTime)
		app.performance.TotalRenderTime += updateDuration
		app.performance.RenderCount++
		app.performance.AvgRenderTime = app.performance.TotalRenderTime / time.Duration(app.performance.RenderCount)
		app.performance.LastUpdate = time.Now()

		return app, cmd
	}

	return app, nil
}

// View implements tea.Model
func (app *AppModel) View() string {
	if app.quitting {
		return ""
	}

	if !app.terminalReady {
		return "Inicializando terminal..."
	}

	var sections []string

	// Application header with metadata
	if app.debug {
		sections = append(sections, app.renderDebugHeader())
	}

	// Main content from active model
	if app.activeModel != nil {
		// Type assertion to access View method
		if modelWithView, ok := app.activeModel.(interface{ View() string }); ok {
			content := modelWithView.View()
			sections = append(sections, content)
		}
	}

	// Error display
	if len(app.errors) > 0 {
		errorSection := app.renderErrors()
		sections = append(sections, errorSection)
	}

	// Debug footer
	if app.debug {
		sections = append(sections, app.renderDebugFooter())
	}

	// Global navigation help
	sections = append(sections, app.renderGlobalHelp())

	return app.theme.Border.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

// renderDebugHeader renders debug information in the header
func (app *AppModel) renderDebugHeader() string {
	header := fmt.Sprintf("Shantilly v%s | View: %s | Components: %d | Errors: %d",
		app.metadata.Version,
		app.currentView.String(),
		len(app.components),
		len(app.errors),
	)
	return fmt.Sprintf("\x1b[36m%s\x1b[0m", header) // Cyan color for debug info
}

// renderErrors renders current application errors
func (app *AppModel) renderErrors() string {
	var errorLines []string
	for _, err := range app.errors {
		if !err.Resolved {
			errorLine := fmt.Sprintf("✗ %s: %s", err.Code.String(), err.Message)
			errorLines = append(errorLines, errorLine)
		}
	}
	return app.theme.Error.Render(strings.Join(errorLines, "\n"))
}

// renderDebugFooter renders debug information in the footer
func (app *AppModel) renderDebugFooter() string {
	footer := fmt.Sprintf("Performance: %v avg render | Memory: %d MB | Validation: %s",
		app.performance.AvgRenderTime.Round(time.Millisecond),
		app.performance.MemoryUsage/1024/1024,
		map[bool]string{true: "✓", false: "✗"}[app.validation.IsValid],
	)
	return fmt.Sprintf("\x1b[36m%s\x1b[0m", footer) // Cyan color for debug info
}

// renderGlobalHelp renders global navigation help
func (app *AppModel) renderGlobalHelp() string {
	help := "F1: Debug | F2: Next View | F12: Stats | Esc: Quit"
	return app.theme.Help.Render(help)
}

// navigateToNextView handles view transitions
func (app *AppModel) navigateToNextView() tea.Cmd {
	app.previousView = app.currentView

	// Cycle through available views
	switch app.currentView {
	case FormView:
		if len(app.config.Layouts) > 0 {
			app.currentView = LayoutView
		} else if len(app.config.Tabs) > 0 {
			app.currentView = TabsView
		}
	case LayoutView:
		if len(app.config.Tabs) > 0 {
			app.currentView = TabsView
		} else if len(app.config.Forms) > 0 {
			app.currentView = FormView
		}
	case TabsView:
		if len(app.config.Forms) > 0 {
			app.currentView = FormView
		} else if len(app.config.Layouts) > 0 {
			app.currentView = LayoutView
		}
	}

	// Add to navigation history
	app.navigationHistory = append(app.navigationHistory, app.currentView)
	app.navigationIndex++

	// Reinitialize the view
	return func() tea.Msg {
		if err := app.initializeView(); err != nil {
			app.addError(ErrStateManagementFailed, "Failed to initialize view", "AppModel", SeverityError, map[string]interface{}{
				"from_view": app.previousView.String(),
				"to_view":   app.currentView.String(),
				"error":     err.Error(),
			})
		}
		return nil
	}
}

// addError adds a new error to the application error list
func (app *AppModel) addError(code ErrorCode, message, component string, severity ErrorSeverity, context map[string]interface{}) {
	app.lastErrorID++
	error := AppError{
		ID:        fmt.Sprintf("err_%d", app.lastErrorID),
		Code:      code,
		Message:   message,
		Component: component,
		Severity:  severity,
		Context:   context,
		Timestamp: time.Now(),
		Resolved:  false,
	}
	app.errors = append(app.errors, error)
}

// resolveError marks an error as resolved
func (app *AppModel) resolveError(errorID string) {
	for i := range app.errors {
		if app.errors[i].ID == errorID {
			app.errors[i].Resolved = true
			break
		}
	}
}

// logAppStats logs current application statistics for debugging
func (app *AppModel) logAppStats() {
	stats := map[string]interface{}{
		"current_view":     app.currentView.String(),
		"components_count": len(app.components),
		"errors_count":     len(app.errors),
		"performance":      app.performance,
		"validation":       app.validation,
		"navigation_depth": len(app.navigationHistory),
		"terminal_size":    fmt.Sprintf("%dx%d", app.width, app.height),
	}

	if data, err := json.MarshalIndent(stats, "", "  "); err == nil {
		// In a real implementation, this would be logged
		fmt.Printf("App Stats: %s\n", string(data))
	}
}

// GetCurrentView returns the current view type
func (app *AppModel) GetCurrentView() ViewType {
	return app.currentView
}

// GetActiveModel returns the currently active model
func (app *AppModel) GetActiveModel() tea.Model {
	return app.activeModel
}

// GetErrors returns all unresolved errors
func (app *AppModel) GetErrors() []AppError {
	var unresolved []AppError
	for _, err := range app.errors {
		if !err.Resolved {
			unresolved = append(unresolved, err)
		}
	}
	return unresolved
}

// IsQuitting returns true if the application is quitting
func (app *AppModel) IsQuitting() bool {
	return app.quitting
}

// String returns a string representation of the ErrorCode
func (ec ErrorCode) String() string {
	switch ec {
	case ErrComponentNotFound:
		return "COMPONENT_NOT_FOUND"
	case ErrValidationFailed:
		return "VALIDATION_FAILED"
	case ErrConfigInvalid:
		return "CONFIG_INVALID"
	case ErrThemeLoadFailed:
		return "THEME_LOAD_FAILED"
	case ErrComponentCreationFailed:
		return "COMPONENT_CREATION_FAILED"
	case ErrStateManagementFailed:
		return "STATE_MANAGEMENT_FAILED"
	case ErrSerializationFailed:
		return "SERIALIZATION_FAILED"
	case ErrFileOperationFailed:
		return "FILE_OPERATION_FAILED"
	case ErrNetworkOperationFailed:
		return "NETWORK_OPERATION_FAILED"
	case ErrPermissionDenied:
		return "PERMISSION_DENIED"
	case ErrResourceNotFound:
		return "RESOURCE_NOT_FOUND"
	case ErrTimeout:
		return "TIMEOUT"
	case ErrConcurrencyIssue:
		return "CONCURRENCY_ISSUE"
	case ErrMemoryAllocationFailed:
		return "MEMORY_ALLOCATION_FAILED"
	case ErrInvalidViewTransition:
		return "INVALID_VIEW_TRANSITION"
	case ErrComponentDependencyFailed:
		return "COMPONENT_DEPENDENCY_FAILED"
	default:
		return fmt.Sprintf("UNKNOWN_ERROR_%d", int(ec))
	}
}
