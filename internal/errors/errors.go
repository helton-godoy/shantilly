package errors

import (
	"fmt"
	"runtime"
	"strings"
	"time"
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

// String returns the string representation of ErrorSeverity
func (es ErrorSeverity) String() string {
	switch es {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	case SeverityCritical:
		return "CRITICAL"
	case SeverityFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ErrorCode defines standardized error codes for the application
type ErrorCode int

const (
	// Component errors (1000-1099)
	ErrComponentNotFound ErrorCode = iota + 1000
	ErrComponentCreationFailed
	ErrComponentValidationFailed
	ErrComponentDependencyFailed
	ErrComponentStateInvalid

	// Configuration errors (1100-1199)
	ErrConfigInvalid ErrorCode = iota + 1100
	ErrConfigLoadFailed
	ErrConfigValidationFailed
	ErrConfigNotFound

	// Theme errors (1200-1299)
	ErrThemeLoadFailed ErrorCode = iota + 1200
	ErrThemeValidationFailed
	ErrThemeNotFound

	// Model errors (1300-1399)
	ErrModelCreationFailed ErrorCode = iota + 1300
	ErrModelStateInvalid
	ErrModelNavigationFailed

	// File system errors (1400-1499)
	ErrFileOperationFailed ErrorCode = iota + 1400
	ErrFileNotFound
	ErrFilePermissionDenied
	ErrFileInvalidFormat

	// Network errors (1500-1599)
	ErrNetworkOperationFailed ErrorCode = iota + 1500
	ErrNetworkTimeout
	ErrNetworkConnectionFailed

	// Validation errors (1600-1699)
	ErrValidationFailed ErrorCode = iota + 1600
	ErrValidationCrossFieldFailed
	ErrValidationBusinessRuleFailed

	// Runtime errors (1700-1799)
	ErrMemoryAllocationFailed ErrorCode = iota + 1700
	ErrConcurrencyIssue
	ErrTimeout
	ErrResourceExhausted
)

// String returns the string representation of ErrorCode
func (ec ErrorCode) String() string {
	switch ec {
	case ErrComponentNotFound:
		return "COMPONENT_NOT_FOUND"
	case ErrComponentCreationFailed:
		return "COMPONENT_CREATION_FAILED"
	case ErrComponentValidationFailed:
		return "COMPONENT_VALIDATION_FAILED"
	case ErrComponentDependencyFailed:
		return "COMPONENT_DEPENDENCY_FAILED"
	case ErrComponentStateInvalid:
		return "COMPONENT_STATE_INVALID"
	case ErrConfigInvalid:
		return "CONFIG_INVALID"
	case ErrConfigLoadFailed:
		return "CONFIG_LOAD_FAILED"
	case ErrConfigValidationFailed:
		return "CONFIG_VALIDATION_FAILED"
	case ErrConfigNotFound:
		return "CONFIG_NOT_FOUND"
	case ErrThemeLoadFailed:
		return "THEME_LOAD_FAILED"
	case ErrThemeValidationFailed:
		return "THEME_VALIDATION_FAILED"
	case ErrThemeNotFound:
		return "THEME_NOT_FOUND"
	case ErrModelCreationFailed:
		return "MODEL_CREATION_FAILED"
	case ErrModelStateInvalid:
		return "MODEL_STATE_INVALID"
	case ErrModelNavigationFailed:
		return "MODEL_NAVIGATION_FAILED"
	case ErrFileOperationFailed:
		return "FILE_OPERATION_FAILED"
	case ErrFileNotFound:
		return "FILE_NOT_FOUND"
	case ErrFilePermissionDenied:
		return "FILE_PERMISSION_DENIED"
	case ErrFileInvalidFormat:
		return "FILE_INVALID_FORMAT"
	case ErrNetworkOperationFailed:
		return "NETWORK_OPERATION_FAILED"
	case ErrNetworkTimeout:
		return "NETWORK_TIMEOUT"
	case ErrNetworkConnectionFailed:
		return "NETWORK_CONNECTION_FAILED"
	case ErrValidationFailed:
		return "VALIDATION_FAILED"
	case ErrValidationCrossFieldFailed:
		return "VALIDATION_CROSS_FIELD_FAILED"
	case ErrValidationBusinessRuleFailed:
		return "VALIDATION_BUSINESS_RULE_FAILED"
	case ErrMemoryAllocationFailed:
		return "MEMORY_ALLOCATION_FAILED"
	case ErrConcurrencyIssue:
		return "CONCURRENCY_ISSUE"
	case ErrTimeout:
		return "TIMEOUT"
	case ErrResourceExhausted:
		return "RESOURCE_EXHAUSTED"
	default:
		return fmt.Sprintf("UNKNOWN_ERROR_%d", int(ec))
	}
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
	Retryable  bool                   `json:"retryable"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Severity.String(), e.Code.String(), e.Message)
}

// NewAppError creates a new AppError with the specified parameters
func NewAppError(code ErrorCode, message, component string, severity ErrorSeverity) *AppError {
	return &AppError{
		ID:         generateErrorID(),
		Code:       code,
		Message:    message,
		Component:  component,
		Severity:   severity,
		Context:    make(map[string]interface{}),
		StackTrace: captureStackTrace(),
		Timestamp:  time.Now(),
		Resolved:   false,
		Retryable:  isRetryable(code),
	}
}

// WithContext adds context information to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithContextMap adds multiple context values to the error
func (e *AppError) WithContextMap(context map[string]interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	for key, value := range context {
		e.Context[key] = value
	}
	return e
}

// Resolve marks the error as resolved
func (e *AppError) Resolve() {
	e.Resolved = true
}

// IsRetryable returns true if the error is retryable
func (e *AppError) IsRetryable() bool {
	return e.Retryable
}

// IsCritical returns true if the error is critical or fatal
func (e *AppError) IsCritical() bool {
	return e.Severity == SeverityCritical || e.Severity == SeverityFatal
}

// ErrorManager manages application errors with advanced features
type ErrorManager struct {
	errors      []AppError
	maxErrors   int
	listeners   []ErrorListener
	autoResolve bool
	filterLevel ErrorSeverity
}

// ErrorListener defines the interface for error listeners
type ErrorListener interface {
	OnError(error *AppError)
	OnErrorResolved(errorID string)
}

// NewErrorManager creates a new ErrorManager
func NewErrorManager(maxErrors int) *ErrorManager {
	return &ErrorManager{
		errors:      make([]AppError, 0),
		maxErrors:   maxErrors,
		listeners:   make([]ErrorListener, 0),
		autoResolve: false,
		filterLevel: SeverityInfo,
	}
}

// AddError adds a new error to the manager
func (em *ErrorManager) AddError(err *AppError) {
	// Filter by severity level
	if err.Severity < em.filterLevel {
		return
	}

	em.errors = append(em.errors, *err)

	// Limit the number of stored errors
	if len(em.errors) > em.maxErrors {
		em.errors = em.errors[len(em.errors)-em.maxErrors:]
	}

	// Notify listeners
	for _, listener := range em.listeners {
		listener.OnError(err)
	}
}

// ResolveError marks an error as resolved
func (em *ErrorManager) ResolveError(errorID string) {
	for i := range em.errors {
		if em.errors[i].ID == errorID {
			em.errors[i].Resolved = true

			// Notify listeners
			for _, listener := range em.listeners {
				listener.OnErrorResolved(errorID)
			}
			break
		}
	}
}

// GetErrors returns all errors, optionally filtered by severity
func (em *ErrorManager) GetErrors(minSeverity ErrorSeverity) []AppError {
	var filtered []AppError
	for _, err := range em.errors {
		if err.Severity >= minSeverity {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

// GetUnresolvedErrors returns all unresolved errors
func (em *ErrorManager) GetUnresolvedErrors() []AppError {
	var unresolved []AppError
	for _, err := range em.errors {
		if !err.Resolved {
			unresolved = append(unresolved, err)
		}
	}
	return unresolved
}

// ClearResolved removes all resolved errors
func (em *ErrorManager) ClearResolved() {
	var active []AppError
	for _, err := range em.errors {
		if !err.Resolved {
			active = append(active, err)
		}
	}
	em.errors = active
}

// AddListener adds an error listener
func (em *ErrorManager) AddListener(listener ErrorListener) {
	em.listeners = append(em.listeners, listener)
}

// SetFilterLevel sets the minimum severity level for error filtering
func (em *ErrorManager) SetFilterLevel(level ErrorSeverity) {
	em.filterLevel = level
}

// Helper functions

func generateErrorID() string {
	return fmt.Sprintf("err_%d", time.Now().UnixNano())
}

func captureStackTrace() string {
	// Skip this function and the caller
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(2, pcs[:])

	frames := runtime.CallersFrames(pcs[:n])
	var stack strings.Builder

	for {
		frame, more := frames.Next()
		stack.WriteString(fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function))

		if !more {
			break
		}
	}

	return stack.String()
}

func isRetryable(code ErrorCode) bool {
	switch code {
	case ErrNetworkTimeout, ErrNetworkConnectionFailed, ErrFilePermissionDenied:
		return true
	case ErrComponentValidationFailed, ErrValidationFailed:
		return false
	case ErrMemoryAllocationFailed, ErrResourceExhausted:
		return false
	default:
		return true // Most errors are retryable by default
	}
}

// Common error creation functions

// NewComponentError creates a component-related error
func NewComponentError(message, component string, severity ErrorSeverity) *AppError {
	code := ErrComponentCreationFailed
	if severity == SeverityError {
		code = ErrComponentValidationFailed
	}
	return NewAppError(code, message, component, severity)
}

// NewValidationError creates a validation-related error
func NewValidationError(message, field string) *AppError {
	return NewAppError(ErrValidationFailed, message, field, SeverityWarning)
}

// NewConfigError creates a configuration-related error
func NewConfigError(message string) *AppError {
	return NewAppError(ErrConfigInvalid, message, "config", SeverityError)
}

// NewFileError creates a file operation error
func NewFileError(message, filePath string) *AppError {
	return NewAppError(ErrFileOperationFailed, message, filePath, SeverityError)
}

// NewNetworkError creates a network-related error
func NewNetworkError(message string) *AppError {
	return NewAppError(ErrNetworkOperationFailed, message, "network", SeverityWarning)
}
