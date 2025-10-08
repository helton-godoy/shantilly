package components

import tea "github.com/charmbracelet/bubbletea/v2"

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
}
