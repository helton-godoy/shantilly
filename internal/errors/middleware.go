package errors

import (
	"fmt"
	"log"
	"time"
)

// ErrorHandler defines the interface for error handling strategies
type ErrorHandler interface {
	Handle(err *AppError) error
	CanHandle(err *AppError) bool
	GetPriority() int
}

// RecoveryStrategy defines how to recover from errors
type RecoveryStrategy int

const (
	RecoveryRetry RecoveryStrategy = iota
	RecoveryFallback
	RecoverySkip
	RecoveryFail
)

// ErrorMiddleware provides error handling middleware functionality
type ErrorMiddleware struct {
	handlers       []ErrorHandler
	errorManager   *ErrorManager
	enableLogging  bool
	enableRecovery bool
	retryConfig    RetryConfig
}

// RetryConfig defines retry behavior for recoverable errors
type RetryConfig struct {
	MaxRetries      int
	BaseDelay       time.Duration
	MaxDelay        time.Duration
	Multiplier      float64
	RetryableErrors []ErrorCode
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: 3,
		BaseDelay:  time.Second,
		MaxDelay:   time.Minute,
		Multiplier: 2.0,
		RetryableErrors: []ErrorCode{
			ErrNetworkTimeout,
			ErrNetworkConnectionFailed,
			ErrFilePermissionDenied,
			ErrTimeout,
		},
	}
}

// NewErrorMiddleware creates a new ErrorMiddleware instance
func NewErrorMiddleware(manager *ErrorManager) *ErrorMiddleware {
	return &ErrorMiddleware{
		handlers:       make([]ErrorHandler, 0),
		errorManager:   manager,
		enableLogging:  true,
		enableRecovery: true,
		retryConfig:    DefaultRetryConfig(),
	}
}

// AddHandler adds an error handler to the middleware
func (em *ErrorMiddleware) AddHandler(handler ErrorHandler) {
	em.handlers = append(em.handlers, handler)

	// Sort handlers by priority (higher priority first)
	for i := len(em.handlers) - 1; i > 0; i-- {
		if em.handlers[i].GetPriority() > em.handlers[i-1].GetPriority() {
			em.handlers[i], em.handlers[i-1] = em.handlers[i-1], em.handlers[i]
		}
	}
}

// ProcessError processes an error through the middleware pipeline
func (em *ErrorMiddleware) ProcessError(err *AppError) error {
	// Add to error manager
	em.errorManager.AddError(err)

	// Log if enabled
	if em.enableLogging {
		em.logError(err)
	}

	// Try recovery if enabled
	if em.enableRecovery {
		if recovered := em.tryRecovery(err); recovered {
			err.Resolve()
			em.errorManager.ResolveError(err.ID)
			return nil
		}
	}

	// Try handlers
	for _, handler := range em.handlers {
		if handler.CanHandle(err) {
			if recoveryErr := handler.Handle(err); recoveryErr == nil {
				err.Resolve()
				em.errorManager.ResolveError(err.ID)
				return nil
			}
		}
	}

	return err
}

// tryRecovery attempts to recover from the error automatically
func (em *ErrorMiddleware) tryRecovery(err *AppError) bool {
	// Check if error is retryable
	if !err.IsRetryable() {
		return false
	}

	// Check if error code is in retryable list
	for _, retryableCode := range em.retryConfig.RetryableErrors {
		if err.Code == retryableCode {
			return em.executeRetryStrategy(err)
		}
	}

	return false
}

// executeRetryStrategy executes the retry strategy for an error
func (em *ErrorMiddleware) executeRetryStrategy(err *AppError) bool {
	// This would implement exponential backoff retry logic
	// For now, we'll implement a simple retry mechanism

	delay := em.retryConfig.BaseDelay
	for i := 0; i < em.retryConfig.MaxRetries; i++ {
		time.Sleep(delay)

		// In a real implementation, this would retry the operation
		// For now, we'll simulate a successful retry for network errors
		if err.Code == ErrNetworkTimeout || err.Code == ErrNetworkConnectionFailed {
			// Simulate network recovery
			if i == em.retryConfig.MaxRetries-1 {
				return true // Success on last retry
			}
		}

		delay = time.Duration(float64(delay) * em.retryConfig.Multiplier)
		if delay > em.retryConfig.MaxDelay {
			delay = em.retryConfig.MaxDelay
		}
	}

	return false
}

// logError logs an error with appropriate formatting
func (em *ErrorMiddleware) logError(err *AppError) {
	logLevel := getLogLevel(err.Severity)

	if err.IsCritical() {
		log.Printf("[%s] CRITICAL ERROR [%s] %s in %s: %s",
			logLevel, err.ID, err.Code.String(), err.Component, err.Message)
	} else {
		log.Printf("[%s] ERROR [%s] %s: %s",
			logLevel, err.ID, err.Code.String(), err.Message)
	}

	// Log context if present
	if len(err.Context) > 0 {
		log.Printf("  Context: %+v", err.Context)
	}

	// Log stack trace for critical errors
	if err.IsCritical() && err.StackTrace != "" {
		log.Printf("  Stack Trace:\n%s", err.StackTrace)
	}
}

// getLogLevel converts error severity to log level
func getLogLevel(severity ErrorSeverity) string {
	switch severity {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARN"
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

// RecoveryHandler implements automatic error recovery strategies
type RecoveryHandler struct {
	strategies map[ErrorCode]RecoveryStrategy
	priority   int
}

// NewRecoveryHandler creates a new recovery handler
func NewRecoveryHandler() *RecoveryHandler {
	return &RecoveryHandler{
		strategies: map[ErrorCode]RecoveryStrategy{
			ErrNetworkTimeout:          RecoveryRetry,
			ErrNetworkConnectionFailed: RecoveryRetry,
			ErrFilePermissionDenied:    RecoveryRetry,
			ErrValidationFailed:        RecoverySkip,
			ErrComponentNotFound:       RecoveryFallback,
			ErrConfigInvalid:           RecoveryFail,
		},
		priority: 100,
	}
}

// Handle implements ErrorHandler
func (rh *RecoveryHandler) Handle(err *AppError) error {
	strategy, exists := rh.strategies[err.Code]
	if !exists {
		return fmt.Errorf("no recovery strategy for error code: %s", err.Code.String())
	}

	switch strategy {
	case RecoveryRetry:
		return rh.handleRetry(err)
	case RecoveryFallback:
		return rh.handleFallback(err)
	case RecoverySkip:
		return nil // Skip the operation
	case RecoveryFail:
		return err // Fail fast
	default:
		return fmt.Errorf("unknown recovery strategy")
	}
}

// CanHandle implements ErrorHandler
func (rh *RecoveryHandler) CanHandle(err *AppError) bool {
	_, exists := rh.strategies[err.Code]
	return exists
}

// GetPriority implements ErrorHandler
func (rh *RecoveryHandler) GetPriority() int {
	return rh.priority
}

// handleRetry implements retry logic
func (rh *RecoveryHandler) handleRetry(err *AppError) error {
	// In a real implementation, this would retry the failing operation
	// For now, we'll simulate a retry delay
	time.Sleep(time.Second * 2)
	return nil
}

// handleFallback implements fallback logic
func (rh *RecoveryHandler) handleFallback(err *AppError) error {
	// In a real implementation, this would use a fallback mechanism
	// For example, use default values or alternative services
	return nil
}

// LoggingHandler implements error logging functionality
type LoggingHandler struct {
	logger   *log.Logger
	priority int
}

// NewLoggingHandler creates a new logging handler
func NewLoggingHandler(logger *log.Logger) *LoggingHandler {
	return &LoggingHandler{
		logger:   logger,
		priority: 50,
	}
}

// Handle implements ErrorHandler
func (lh *LoggingHandler) Handle(err *AppError) error {
	lh.logger.Printf("Error handled: %s", err.Error())
	return nil
}

// CanHandle implements ErrorHandler
func (lh *LoggingHandler) CanHandle(err *AppError) bool {
	return true // Handle all errors
}

// GetPriority implements ErrorHandler
func (lh *LoggingHandler) GetPriority() int {
	return lh.priority
}

// NotificationHandler implements error notification functionality
type NotificationHandler struct {
	notifier NotificationService
	priority int
}

// NotificationService defines the interface for notification services
type NotificationService interface {
	SendNotification(title, message string, severity ErrorSeverity) error
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notifier NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notifier: notifier,
		priority: 200,
	}
}

// Handle implements ErrorHandler
func (nh *NotificationHandler) Handle(err *AppError) error {
	if err.IsCritical() {
		title := fmt.Sprintf("Critical Error: %s", err.Code.String())
		message := fmt.Sprintf("Component: %s\nMessage: %s", err.Component, err.Message)
		return nh.notifier.SendNotification(title, message, err.Severity)
	}
	return nil
}

// CanHandle implements ErrorHandler
func (nh *NotificationHandler) CanHandle(err *AppError) bool {
	return err.IsCritical()
}

// GetPriority implements ErrorHandler
func (nh *NotificationHandler) GetPriority() int {
	return nh.priority
}

// ErrorContext provides contextual information for error handling
type ErrorContext struct {
	Operation string                 `json:"operation"`
	UserID    string                 `json:"user_id,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Component string                 `json:"component"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}

// NewErrorContext creates a new error context
func NewErrorContext(operation, component string) *ErrorContext {
	return &ErrorContext{
		Operation: operation,
		Component: component,
		Metadata:  make(map[string]interface{}),
		Timestamp: time.Now(),
	}
}

// WithUser adds user information to the context
func (ec *ErrorContext) WithUser(userID string) *ErrorContext {
	ec.UserID = userID
	return ec
}

// WithSession adds session information to the context
func (ec *ErrorContext) WithSession(sessionID string) *ErrorContext {
	ec.SessionID = sessionID
	return ec
}

// WithRequest adds request information to the context
func (ec *ErrorContext) WithRequest(requestID string) *ErrorContext {
	ec.RequestID = requestID
	return ec
}

// WithMetadata adds metadata to the context
func (ec *ErrorContext) WithMetadata(key string, value interface{}) *ErrorContext {
	if ec.Metadata == nil {
		ec.Metadata = make(map[string]interface{})
	}
	ec.Metadata[key] = value
	return ec
}

// WrapError wraps an error with additional context
func WrapError(err error, context *ErrorContext) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr.WithContextMap(map[string]interface{}{
			"operation": context.Operation,
			"metadata":  context.Metadata,
		})
	}

	// Create new error from standard error
	appErr := NewAppError(
		ErrComponentCreationFailed,
		err.Error(),
		context.Component,
		SeverityError,
	).WithContextMap(map[string]interface{}{
		"operation": context.Operation,
		"metadata":  context.Metadata,
	})

	return appErr
}
