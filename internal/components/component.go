package components

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/styles"
)

// ComponentMetadata contains comprehensive information about a component.
type ComponentMetadata struct {
	Version      string                 `json:"version"`
	Author       string                 `json:"author"`
	Description  string                 `json:"description"`
	Dependencies []string               `json:"dependencies"`
	Examples     []ComponentExample     `json:"examples"`
	Schema       map[string]interface{} `json:"schema"`
}

// ComponentExample represents an example usage of a component.
type ComponentExample struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Config         map[string]interface{} `json:"config"`
	ExpectedOutput interface{}            `json:"expected_output"`
}

// ValidationContext provides context for component validation.
type ValidationContext struct {
	ComponentValues map[string]interface{} `json:"component_values"`
	GlobalConfig    map[string]interface{} `json:"global_config"`
	ExternalData    map[string]interface{} `json:"external_data"`
}

// ValidationError represents a validation error with detailed information.
type ValidationError struct {
	Code     string                 `json:"code"`
	Message  string                 `json:"message"`
	Field    string                 `json:"field"`
	Severity string                 `json:"severity"`
	Context  map[string]interface{} `json:"context"`
}

// ExportFormat defines supported export formats for component data.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatYAML ExportFormat = "yaml"
	FormatXML  ExportFormat = "xml"
	FormatCSV  ExportFormat = "csv"
)

// Component defines the contract for all UI widgets in Shantilly.
// Every component must implement this interface completely to ensure
// interoperability with orchestration models (FormModel, LayoutModel, TabsModel).
//
// This interface extends tea.Model (Init, Update, View) and adds methods
// for focus management, validation, value serialization, and error handling.
//
// Design Principles:
// - Components must be self-contained and manage their own state
// - Orchestration models handle focus navigation (tab/shift+tab)
// - Components receive messages only when focused (except WindowSizeMsg)
// - Validation is declarative and based on ComponentConfig
type Component interface {
	// MVU architecture methods
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string

	// Name returns the unique identifier of the component.
	// This is used for JSON serialization in FormModel.ToJSON().
	// The name must match the 'name' field in ComponentConfig.
	Name() string

	// CanFocus returns true if the component is interactive and can receive focus.
	// Static components (e.g., Text labels) should return false.
	// This is used by LayoutModel/FormModel to determine focus order.
	CanFocus() bool

	// SetFocus enables or disables focus on the component.
	// When focused:
	// - The component should apply focused styling (theme.InputFocused)
	// - The component receives all tea.KeyMsg updates
	// When not focused:
	// - The component should apply unfocused styling (theme.Input)
	// - The component does not receive keyboard input (except global messages)
	SetFocus(focused bool)

	// IsValid returns true if the component's current state is valid.
	// Validation rules:
	// - If config.Required is true, value must not be empty
	// - All config.Options validation rules must pass (min_length, max_length, pattern, etc.)
	// - Custom business logic validation can be added
	//
	// This method is called by FormModel.CanSubmit() to determine if the form is submittable.
	IsValid() bool

	// GetError returns the current error message for the component.
	// Returns empty string if IsValid() is true.
	// The error should be stored in the component's errorMsg field.
	GetError() string

	// SetError manually sets an error message on the component.
	// This is used for server-side validation or business logic errors
	// that cannot be detected by the component itself.
	SetError(msg string)

	// Value returns the current value of the component for serialization.
	// Return types by component:
	// - TextInput/TextArea: string
	// - Checkbox: bool
	// - RadioGroup: string (selected item ID)
	// - Slider: float64 or int
	// - FilePicker: string (selected file path)
	//
	// This method is called by FormModel.ToJSON() to serialize form data.
	Value() interface{}

	// SetValue programmatically sets the component's value.
	// This is used for:
	// - Initialization from config default values
	// - Reset to initial state
	// - Pre-filling form data
	//
	// Returns an error if the value type is invalid or out of range.
	SetValue(value interface{}) error

	// Reset returns the component to its initial state.
	// This includes:
	// - Resetting value to default or empty
	// - Clearing error message
	// - Resetting internal state (cursor position, selection, etc.)
	Reset()

	// NOVOS métodos para arquitetura híbrida

	// GetMetadata returns comprehensive metadata about the component.
	// This includes version, author, description, dependencies, and examples.
	GetMetadata() ComponentMetadata

	// ValidateWithContext performs validation with additional context.
	// This allows for cross-field validation and business logic validation.
	ValidateWithContext(context ValidationContext) []ValidationError

	// ExportToFormat exports the component's data to a specific format.
	// Supported formats: JSON, YAML, XML, CSV
	ExportToFormat(format ExportFormat) ([]byte, error)

	// ImportFromFormat imports data from a specific format into the component.
	// This is used for data binding and initialization from external sources.
	ImportFromFormat(format ExportFormat, data []byte) error

	// GetDependencies returns a list of components that this component depends on.
	// This is used for dependency injection and component lifecycle management.
	GetDependencies() []string

	// SetTheme allows dynamic theme changes for the component.
	// This enables runtime theme switching without recreating the component.
	SetTheme(theme *styles.Theme)
}
