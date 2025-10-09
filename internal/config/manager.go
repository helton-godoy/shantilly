package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/helton/shantilly/internal/errors"
	"gopkg.in/yaml.v3"
)

// ConfigManager manages application configuration with advanced features
type ConfigManager struct {
	configs        map[string]*Config
	activeConfig   string
	configPaths    []string
	watchers       []ConfigWatcher
	autoReload     bool
	validateOnLoad bool
}

// ConfigWatcher defines the interface for configuration change listeners
type ConfigWatcher interface {
	OnConfigChanged(configName string, newConfig *Config)
	OnConfigError(configName string, err error)
}

// ValidationRule defines a validation rule for configuration
type ValidationRule struct {
	Field    string
	Rule     string
	Value    interface{}
	Message  string
	Severity errors.ErrorSeverity
}

// ConfigLoadOptions defines options for loading configuration
type ConfigLoadOptions struct {
	Validate      bool
	Watch         bool
	Environment   string
	ConfigPaths   []string
	DefaultConfig *Config
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configs:        make(map[string]*Config),
		configPaths:    getDefaultConfigPaths(),
		autoReload:     false,
		validateOnLoad: true,
	}
}

// LoadConfig loads configuration from a file with enhanced error handling and logging
func (cm *ConfigManager) LoadConfig(configPath string, options *ConfigLoadOptions) (*Config, error) {
	if options == nil {
		options = &ConfigLoadOptions{
			Validate: true,
			Watch:    false,
		}
	}

	log.Printf("Loading configuration from: %s", configPath)

	// Read configuration file with enhanced error context
	data, err := os.ReadFile(configPath)
	if err != nil {
		fileErr := errors.NewFileError(fmt.Sprintf("failed to read config file: %s", configPath), configPath)
		log.Printf("Config file read error: %v", fileErr)
		return nil, fileErr
	}

	// Parse YAML with enhanced error handling
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		parseErr := errors.NewConfigError(fmt.Sprintf("failed to parse YAML: %v", err))
		log.Printf("Config YAML parse error: %v", parseErr)
		return nil, parseErr
	}

	// Apply environment-specific overrides
	if options.Environment != "" {
		log.Printf("Applying environment overrides for: %s", options.Environment)
		cm.applyEnvironmentOverrides(&config, options.Environment)
	}

	// Validate configuration with enhanced error reporting
	if options.Validate {
		if err := cm.validateConfig(&config); err != nil {
			log.Printf("Config validation error: %v", err)
			return nil, err
		}
		log.Printf("Configuration validation successful")
	}

	// Store configuration
	configName := cm.getConfigName(configPath)
	cm.configs[configName] = &config
	log.Printf("Configuration stored as: %s", configName)

	// Set as active if it's the first config
	if cm.activeConfig == "" {
		cm.activeConfig = configName
		log.Printf("Set as active configuration: %s", configName)
	}

	// Notify watchers
	cm.notifyWatchers(configName, &config)

	log.Printf("Configuration loaded successfully from: %s", configPath)
	return &config, nil
}

// LoadConfigFromString loads configuration from a YAML string
func (cm *ConfigManager) LoadConfigFromString(yamlContent string, configName string) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal([]byte(yamlContent), &config); err != nil {
		return nil, errors.NewConfigError(fmt.Sprintf("failed to parse YAML string: %v", err))
	}

	if cm.validateOnLoad {
		if err := cm.validateConfig(&config); err != nil {
			return nil, err
		}
	}

	cm.configs[configName] = &config

	// Set as active if it's the first config
	if cm.activeConfig == "" {
		cm.activeConfig = configName
	}

	cm.notifyWatchers(configName, &config)

	return &config, nil
}

// GetConfig returns the configuration by name
func (cm *ConfigManager) GetConfig(configName string) (*Config, error) {
	if config, exists := cm.configs[configName]; exists {
		return config, nil
	}
	return nil, errors.NewConfigError(fmt.Sprintf("configuration not found: %s", configName))
}

// GetActiveConfig returns the currently active configuration
func (cm *ConfigManager) GetActiveConfig() (*Config, error) {
	if cm.activeConfig == "" {
		return nil, errors.NewConfigError("no active configuration set")
	}
	return cm.GetConfig(cm.activeConfig)
}

// SetActiveConfig sets the active configuration
func (cm *ConfigManager) SetActiveConfig(configName string) error {
	if _, exists := cm.configs[configName]; !exists {
		return errors.NewConfigError(fmt.Sprintf("configuration not found: %s", configName))
	}
	cm.activeConfig = configName
	return nil
}

// AddConfigPath adds a path to search for configuration files
func (cm *ConfigManager) AddConfigPath(path string) {
	cm.configPaths = append(cm.configPaths, path)
}

// AddWatcher adds a configuration watcher
func (cm *ConfigManager) AddWatcher(watcher ConfigWatcher) {
	cm.watchers = append(cm.watchers, watcher)
}

// SetAutoReload enables or disables automatic configuration reloading
func (cm *ConfigManager) SetAutoReload(enabled bool) {
	cm.autoReload = enabled
}

// validateConfig performs comprehensive configuration validation
func (cm *ConfigManager) validateConfig(config *Config) error {
	// Basic validation
	if err := config.Validate(); err != nil {
		return err
	}

	// Custom validation rules
	rules := cm.getValidationRules()
	for _, rule := range rules {
		if err := cm.applyValidationRule(config, rule); err != nil {
			return err
		}
	}

	return nil
}

// getValidationRules returns comprehensive validation rules for configuration
func (cm *ConfigManager) getValidationRules() []ValidationRule {
	return []ValidationRule{
		// Global configuration rules
		{
			Field:    "Global.AppName",
			Rule:     "required",
			Value:    nil,
			Message:  "Application name is required",
			Severity: errors.SeverityError,
		},
		{
			Field:    "Global.Version",
			Rule:     "required",
			Value:    nil,
			Message:  "Version is required",
			Severity: errors.SeverityError,
		},
		{
			Field:    "Global.LogLevel",
			Rule:     "one_of",
			Value:    []string{"debug", "info", "warn", "error"},
			Message:  "Log level must be one of: debug, info, warn, error",
			Severity: errors.SeverityWarning,
		},

		// Component validation rules
		{
			Field:    "Validation.Component.StrictMode",
			Rule:     "type",
			Value:    "bool",
			Message:  "Component strict mode must be a boolean",
			Severity: errors.SeverityWarning,
		},

		// Performance configuration rules
		{
			Field:    "Performance.EnableMetrics",
			Rule:     "type",
			Value:    "bool",
			Message:  "Performance metrics enable flag must be a boolean",
			Severity: errors.SeverityWarning,
		},

		// Security configuration rules
		{
			Field:    "Security.EnableCSRF",
			Rule:     "type",
			Value:    "bool",
			Message:  "CSRF protection flag must be a boolean",
			Severity: errors.SeverityWarning,
		},

		// Logging configuration rules
		{
			Field:    "Logging.Level",
			Rule:     "one_of",
			Value:    []string{"debug", "info", "warn", "error"},
			Message:  "Logging level must be one of: debug, info, warn, error",
			Severity: errors.SeverityWarning,
		},

		// Form validation rules
		{
			Field:    "Forms",
			Rule:     "min_length",
			Value:    1,
			Message:  "At least one form configuration is required",
			Severity: errors.SeverityError,
		},
	}
}

// applyValidationRule applies a validation rule to the configuration
func (cm *ConfigManager) applyValidationRule(config *Config, rule ValidationRule) error {
	value := cm.getFieldValue(config, rule.Field)
	if value == nil && rule.Rule == "required" {
		return errors.NewValidationError(rule.Message, rule.Field)
	}

	switch rule.Rule {
	case "required":
		if value == nil || cm.isEmpty(value) {
			return errors.NewValidationError(rule.Message, rule.Field)
		}

	case "one_of":
		if allowedValues, ok := rule.Value.([]string); ok {
			if strValue, ok := value.(string); ok {
				found := false
				for _, allowed := range allowedValues {
					if strValue == allowed {
						found = true
						break
					}
				}
				if !found {
					return errors.NewValidationError(rule.Message, rule.Field)
				}
			}
		}

	case "type":
		expectedType := rule.Value.(string)
		actualType := reflect.TypeOf(value).String()
		if actualType != expectedType {
			return errors.NewValidationError(
				fmt.Sprintf("%s: expected %s, got %s", rule.Message, expectedType, actualType),
				rule.Field,
			)
		}

	case "min_length":
		if minLen, ok := rule.Value.(int); ok {
			if reflect.ValueOf(value).Len() < minLen {
				return errors.NewValidationError(rule.Message, rule.Field)
			}
		}
	}

	return nil
}

// getFieldValue retrieves a field value using reflection
func (cm *ConfigManager) getFieldValue(config *Config, fieldPath string) interface{} {
	parts := strings.Split(fieldPath, ".")
	if len(parts) == 0 {
		return nil
	}

	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for _, part := range parts {
		if v.Kind() == reflect.Struct {
			v = v.FieldByName(cm.capitalizeFirst(part))
		} else if v.Kind() == reflect.Map {
			// Handle map access for nested structures
			return nil // Simplified for now
		} else {
			return nil
		}

		if !v.IsValid() {
			return nil
		}

		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}

	return v.Interface()
}

// capitalizeFirst capitalizes the first letter of a string
func (cm *ConfigManager) capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// isEmpty checks if a value is empty
func (cm *ConfigManager) isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case int, int32, int64:
		return v == 0
	case float32, float64:
		return v == 0.0
	default:
		return reflect.ValueOf(value).IsZero()
	}
}

// applyEnvironmentOverrides applies environment-specific configuration overrides
func (cm *ConfigManager) applyEnvironmentOverrides(config *Config, environment string) {
	// Apply environment-specific settings
	switch environment {
	case "development":
		config.Global.Debug = true
		config.Global.LogLevel = "debug"
		config.Logging.Level = "debug"

	case "production":
		config.Global.Debug = false
		config.Global.LogLevel = "info"
		config.Logging.Level = "info"
		config.Performance.EnableMetrics = true

	case "testing":
		config.Global.Debug = true
		config.Global.LogLevel = "debug"
		config.Logging.Level = "debug"
		config.Validation.Component.StrictMode = false
	}
}

// notifyWatchers notifies all watchers of configuration changes
func (cm *ConfigManager) notifyWatchers(configName string, config *Config) {
	for _, watcher := range cm.watchers {
		watcher.OnConfigChanged(configName, config)
	}
}

// getConfigName extracts the configuration name from the file path
func (cm *ConfigManager) getConfigName(configPath string) string {
	base := filepath.Base(configPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return name
}

// getDefaultConfigPaths returns the default configuration search paths
func getDefaultConfigPaths() []string {
	return []string{
		"./config.yaml",
		"./config.yml",
		"./configs/app.yaml",
		"./configs/app.yml",
		"./internal/config/app.yaml",
		"./internal/config/app.yml",
		"/etc/shantilly/config.yaml",
		"/etc/shantilly/config.yml",
	}
}

// LoadDefaultConfig loads the default configuration
func (cm *ConfigManager) LoadDefaultConfig() (*Config, error) {
	// Try to load from default paths
	for _, path := range cm.configPaths {
		if _, err := os.Stat(path); err == nil {
			return cm.LoadConfig(path, &ConfigLoadOptions{
				Validate: true,
				Watch:    cm.autoReload,
			})
		}
	}

	// If no config file found, use default configuration
	defaultConfig := DefaultConfig()
	cm.configs["default"] = defaultConfig
	cm.activeConfig = "default"

	cm.notifyWatchers("default", defaultConfig)

	return defaultConfig, nil
}

// MergeConfigs merges multiple configurations with precedence
func (cm *ConfigManager) MergeConfigs(baseConfig *Config, overlayConfigs ...*Config) *Config {
	merged := *baseConfig // Copy the base config

	for _, overlay := range overlayConfigs {
		cm.mergeConfigStruct(&merged, overlay)
	}

	return &merged
}

// mergeConfigStruct merges two configuration structs recursively
func (cm *ConfigManager) mergeConfigStruct(base, overlay *Config) {
	baseVal := reflect.ValueOf(base).Elem()
	overlayVal := reflect.ValueOf(overlay).Elem()

	for i := 0; i < baseVal.NumField(); i++ {
		baseField := baseVal.Field(i)
		overlayField := overlayVal.Field(i)

		if overlayField.IsValid() && !cm.isEmpty(overlayField.Interface()) {
			if baseField.CanSet() {
				if baseField.Kind() == reflect.Struct && overlayField.Kind() == reflect.Struct {
					// Recursively merge nested structs
					cm.mergeStructFields(baseField, overlayField)
				} else {
					baseField.Set(overlayField)
				}
			}
		}
	}
}

// mergeStructFields merges fields of two structs
func (cm *ConfigManager) mergeStructFields(base, overlay reflect.Value) {
	for i := 0; i < base.NumField(); i++ {
		baseField := base.Field(i)
		overlayField := overlay.Field(i)

		if overlayField.IsValid() && !cm.isEmpty(overlayField.Interface()) {
			if baseField.CanSet() {
				if baseField.Kind() == reflect.String && overlayField.Kind() == reflect.String {
					if overlayField.String() != "" {
						baseField.SetString(overlayField.String())
					}
				} else if baseField.Kind() == reflect.Bool && overlayField.Kind() == reflect.Bool {
					baseField.SetBool(overlayField.Bool())
				} else if baseField.Kind() == reflect.Int && overlayField.Kind() == reflect.Int {
					baseField.SetInt(overlayField.Int())
				}
			}
		}
	}
}

// ExportConfig exports the current configuration to a file
func (cm *ConfigManager) ExportConfig(configName, filePath string) error {
	config, exists := cm.configs[configName]
	if !exists {
		return errors.NewConfigError(fmt.Sprintf("configuration not found: %s", configName))
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return errors.NewConfigError(fmt.Sprintf("failed to marshal config: %v", err))
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return errors.NewFileError(fmt.Sprintf("failed to write config file: %s", filePath), filePath)
	}

	return nil
}

// ValidateAllConfigs validates all loaded configurations
func (cm *ConfigManager) ValidateAllConfigs() error {
	var validationErrors []error

	for name, config := range cm.configs {
		if err := cm.validateConfig(config); err != nil {
			validationErrors = append(validationErrors,
				fmt.Errorf("validation failed for config %s: %w", name, err))
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("multiple validation errors: %v", validationErrors)
	}

	return nil
}

// GetConfigSummary returns a summary of all loaded configurations
func (cm *ConfigManager) GetConfigSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"active_config":    cm.activeConfig,
		"loaded_configs":   len(cm.configs),
		"config_names":     cm.getConfigNames(),
		"auto_reload":      cm.autoReload,
		"validate_on_load": cm.validateOnLoad,
		"config_paths":     cm.configPaths,
		"watchers_count":   len(cm.watchers),
	}

	return summary
}

// getConfigNames returns the names of all loaded configurations
func (cm *ConfigManager) getConfigNames() []string {
	names := make([]string, 0, len(cm.configs))
	for name := range cm.configs {
		names = append(names, name)
	}
	return names
}

// CreateEnvironmentConfig creates a configuration optimized for a specific environment
func (cm *ConfigManager) CreateEnvironmentConfig(baseConfig *Config, environment string) *Config {
	config := *baseConfig // Copy the base config

	// Apply environment-specific optimizations
	switch environment {
	case "development":
		config.Global.Environment = "development"
		config.Global.Debug = true
		config.Global.LogLevel = "debug"
		config.Logging.Level = "debug"
		config.Validation.Component.StrictMode = false
		config.Performance.EnableMetrics = false

	case "testing":
		config.Global.Environment = "testing"
		config.Global.Debug = true
		config.Global.LogLevel = "debug"
		config.Logging.Level = "debug"
		config.Validation.Component.StrictMode = false
		config.Performance.EnableMetrics = false

	case "staging":
		config.Global.Environment = "staging"
		config.Global.Debug = false
		config.Global.LogLevel = "info"
		config.Logging.Level = "info"
		config.Validation.Component.StrictMode = true
		config.Performance.EnableMetrics = true

	case "production":
		config.Global.Environment = "production"
		config.Global.Debug = false
		config.Global.LogLevel = "warn"
		config.Logging.Level = "warn"
		config.Validation.Component.StrictMode = true
		config.Performance.EnableMetrics = true
		config.Security.EnableCSRF = true
	}

	return &config
}

// LoadConfigWithDefaults loads configuration with default fallbacks
func (cm *ConfigManager) LoadConfigWithDefaults(configPath string, defaults *Config) (*Config, error) {
	options := &ConfigLoadOptions{
		Validate:      true,
		DefaultConfig: defaults,
	}

	// Try to load the specified config
	if _, err := os.Stat(configPath); err == nil {
		if config, err := cm.LoadConfig(configPath, options); err == nil {
			return config, nil
		}
	}

	// Fall back to defaults if specified config doesn't exist or fails to load
	if defaults != nil {
		cm.configs["default"] = defaults
		cm.activeConfig = "default"
		cm.notifyWatchers("default", defaults)
		return defaults, nil
	}

	return nil, errors.NewConfigError(fmt.Sprintf("failed to load config from %s and no defaults provided", configPath))
}
