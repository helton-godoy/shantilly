package config

import (
	"fmt"
	"time"
)

// Config represents the complete application configuration with hierarchical structure
type Config struct {
	Global      GlobalConfig           `yaml:"global" json:"global"`
	Forms       []FormConfig           `yaml:"forms,omitempty" json:"forms,omitempty"`
	Layouts     []LayoutConfig         `yaml:"layouts,omitempty" json:"layouts,omitempty"`
	Tabs        []TabsConfig           `yaml:"tabs,omitempty" json:"tabs,omitempty"`
	Menus       []MenuConfig           `yaml:"menus,omitempty" json:"menus,omitempty"`
	Themes      map[string]ThemeConfig `yaml:"themes,omitempty" json:"themes,omitempty"`
	Validation  ValidationConfig       `yaml:"validation,omitempty" json:"validation,omitempty"`
	Logging     LoggingConfig          `yaml:"logging,omitempty" json:"logging,omitempty"`
	Performance PerformanceConfig      `yaml:"performance,omitempty" json:"performance,omitempty"`
	Security    SecurityConfig         `yaml:"security,omitempty" json:"security,omitempty"`
}

// GlobalConfig contains global application configuration
type GlobalConfig struct {
	AppName      string            `yaml:"app_name" json:"app_name"`
	Version      string            `yaml:"version" json:"version"`
	Environment  string            `yaml:"environment" json:"environment"`
	Debug        bool              `yaml:"debug" json:"debug"`
	LogLevel     string            `yaml:"log_level" json:"log_level"`
	DefaultTheme string            `yaml:"default_theme" json:"default_theme"`
	DefaultView  string            `yaml:"default_view" json:"default_view"`
	Metadata     map[string]string `yaml:"metadata" json:"metadata"`
	BuildTime    time.Time         `yaml:"build_time" json:"build_time"`
	GitCommit    string            `yaml:"git_commit" json:"git_commit"`
}

// ThemeConfig contains theme configuration
type ThemeConfig struct {
	BaseTheme    string                    `yaml:"base_theme" json:"base_theme"`
	Extends      []string                  `yaml:"extends,omitempty" json:"extends,omitempty"`
	CustomStyles map[string]StyleConfig    `yaml:"custom_styles" json:"custom_styles"`
	ColorPalette ColorPalette              `yaml:"color_palette" json:"color_palette"`
	Font         FontConfig                `yaml:"font" json:"font"`
	Spacing      SpacingConfig             `yaml:"spacing" json:"spacing"`
	Components   map[string]ComponentStyle `yaml:"components" json:"components"`
}

// StyleConfig contains style configuration
type StyleConfig struct {
	Foreground string `yaml:"foreground" json:"foreground"`
	Background string `yaml:"background" json:"background"`
	Bold       bool   `yaml:"bold" json:"bold"`
	Italic     bool   `yaml:"italic" json:"italic"`
	Underline  bool   `yaml:"underline" json:"underline"`
}

// ColorPalette contains color palette configuration
type ColorPalette struct {
	Primary   string `yaml:"primary" json:"primary"`
	Secondary string `yaml:"secondary" json:"secondary"`
	Success   string `yaml:"success" json:"success"`
	Warning   string `yaml:"warning" json:"warning"`
	Error     string `yaml:"error" json:"error"`
	Info      string `yaml:"info" json:"info"`
}

// FontConfig contains font configuration
type FontConfig struct {
	Family string `yaml:"family" json:"family"`
	Size   int    `yaml:"size" json:"size"`
	Style  string `yaml:"style" json:"style"`
}

// SpacingConfig contains spacing configuration
type SpacingConfig struct {
	Margin  int `yaml:"margin" json:"margin"`
	Padding int `yaml:"padding" json:"padding"`
	Gap     int `yaml:"gap" json:"gap"`
}

// ComponentStyle contains component-specific styling
type ComponentStyle struct {
	Border       StyleConfig `yaml:"border" json:"border"`
	BorderActive StyleConfig `yaml:"border_active" json:"border_active"`
	Foreground   string      `yaml:"foreground" json:"foreground"`
	Background   string      `yaml:"background" json:"background"`
}

// ValidationConfig contains validation configuration
type ValidationConfig struct {
	Component  ComponentValidation  `yaml:"component" json:"component"`
	CrossField CrossFieldValidation `yaml:"cross_field" json:"cross_field"`
	Business   BusinessValidation   `yaml:"business" json:"business"`
	Schema     SchemaValidation     `yaml:"schema" json:"schema"`
}

// ComponentValidation contains component-level validation settings
type ComponentValidation struct {
	StrictMode bool `yaml:"strict_mode" json:"strict_mode"`
	RealTime   bool `yaml:"real_time" json:"real_time"`
	DebounceMs int  `yaml:"debounce_ms" json:"debounce_ms"`
}

// CrossFieldValidation contains cross-field validation settings
type CrossFieldValidation struct {
	Enabled          bool `yaml:"enabled" json:"enabled"`
	ValidateOnSubmit bool `yaml:"validate_on_submit" json:"validate_on_submit"`
}

// BusinessValidation contains business logic validation settings
type BusinessValidation struct {
	Enabled         bool   `yaml:"enabled" json:"enabled"`
	RulesPath       string `yaml:"rules_path" json:"rules_path"`
	CustomValidator string `yaml:"custom_validator" json:"custom_validator"`
}

// SchemaValidation contains schema validation settings
type SchemaValidation struct {
	Enabled      bool   `yaml:"enabled" json:"enabled"`
	SchemaPath   string `yaml:"schema_path" json:"schema_path"`
	StrictSchema bool   `yaml:"strict_schema" json:"strict_schema"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level" json:"level"`
	Output     string `yaml:"output" json:"output"`
	FilePath   string `yaml:"file_path" json:"file_path"`
	MaxSize    int    `yaml:"max_size" json:"max_size"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`
	Compress   bool   `yaml:"compress" json:"compress"`
}

// PerformanceConfig contains performance configuration
type PerformanceConfig struct {
	EnableMetrics  bool `yaml:"enable_metrics" json:"enable_metrics"`
	SampleRate     int  `yaml:"sample_rate" json:"sample_rate"`
	BufferSize     int  `yaml:"buffer_size" json:"buffer_size"`
	TimeoutMs      int  `yaml:"timeout_ms" json:"timeout_ms"`
	MaxConcurrency int  `yaml:"max_concurrency" json:"max_concurrency"`
}

// SecurityConfig contains security configuration
type SecurityConfig struct {
	EnableCSRF     bool     `yaml:"enable_csrf" json:"enable_csrf"`
	TokenExpiry    int      `yaml:"token_expiry" json:"token_expiry"`
	AllowedOrigins []string `yaml:"allowed_origins" json:"allowed_origins"`
	RateLimit      int      `yaml:"rate_limit" json:"rate_limit"`
}

// Validate performs comprehensive validation on the Config
func (c *Config) Validate() error {
	// Validate global configuration
	if err := c.Global.Validate(); err != nil {
		return fmt.Errorf("erro na configuração global: %w", err)
	}

	// Validate forms
	for i, form := range c.Forms {
		if err := form.Validate(); err != nil {
			return fmt.Errorf("erro no formulário %d: %w", i, err)
		}
	}

	// Validate layouts
	for i, layout := range c.Layouts {
		if err := layout.Validate(); err != nil {
			return fmt.Errorf("erro no layout %d: %w", i, err)
		}
	}

	// Validate tabs
	for i, tabs := range c.Tabs {
		if err := tabs.Validate(); err != nil {
			return fmt.Errorf("erro nas abas %d: %w", i, err)
		}
	}

	// Validate themes
	for name, theme := range c.Themes {
		if err := theme.Validate(); err != nil {
			return fmt.Errorf("erro no tema %s: %w", name, err)
		}
	}

	return nil
}

// Validate performs validation on GlobalConfig
func (gc *GlobalConfig) Validate() error {
	if gc.AppName == "" {
		return fmt.Errorf("nome da aplicação é obrigatório")
	}
	if gc.Version == "" {
		return fmt.Errorf("versão é obrigatória")
	}
	if gc.Environment == "" {
		gc.Environment = "development" // Default
	}
	if gc.LogLevel == "" {
		gc.LogLevel = "info" // Default
	}
	if gc.DefaultTheme == "" {
		gc.DefaultTheme = "default" // Default
	}
	return nil
}

// Validate performs validation on ThemeConfig
func (tc *ThemeConfig) Validate() error {
	if tc.BaseTheme == "" {
		tc.BaseTheme = "default"
	}
	// Additional theme validation can be added here
	return nil
}

// DefaultConfig returns a default configuration for the application
func DefaultConfig() *Config {
	return &Config{
		Global: GlobalConfig{
			AppName:      "Shantilly",
			Version:      "1.0.0",
			Environment:  "development",
			Debug:        false,
			LogLevel:     "info",
			DefaultTheme: "default",
			DefaultView:  "form",
			Metadata: map[string]string{
				"description": "Terminal UI framework for interactive forms",
				"author":      "Shantilly Team",
			},
		},
		Validation: ValidationConfig{
			Component: ComponentValidation{
				StrictMode: false,
				RealTime:   true,
				DebounceMs: 300,
			},
			CrossField: CrossFieldValidation{
				Enabled:          true,
				ValidateOnSubmit: true,
			},
			Business: BusinessValidation{
				Enabled:   false,
				RulesPath: "",
			},
			Schema: SchemaValidation{
				Enabled:      false,
				SchemaPath:   "",
				StrictSchema: false,
			},
		},
		Logging: LoggingConfig{
			Level:      "info",
			Output:     "stdout",
			FilePath:   "",
			MaxSize:    100,
			MaxBackups: 3,
			Compress:   true,
		},
		Performance: PerformanceConfig{
			EnableMetrics:  false,
			SampleRate:     1000,
			BufferSize:     1000,
			TimeoutMs:      5000,
			MaxConcurrency: 10,
		},
		Security: SecurityConfig{
			EnableCSRF:     false,
			TokenExpiry:    3600,
			AllowedOrigins: []string{"*"},
			RateLimit:      100,
		},
	}
}
