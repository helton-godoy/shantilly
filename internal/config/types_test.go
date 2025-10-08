package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponentConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ComponentConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid textinput",
			config: ComponentConfig{
				Type: TypeTextInput,
				Name: "username",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: ComponentConfig{
				Type: TypeTextInput,
			},
			wantErr: true,
			errMsg:  "nome do componente é obrigatório",
		},
		{
			name: "invalid type",
			config: ComponentConfig{
				Type: "invalid",
				Name: "test",
			},
			wantErr: true,
			errMsg:  "tipo de componente inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFormConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  FormConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid form",
			config: FormConfig{
				Title: "Test Form",
				Components: []ComponentConfig{
					{Type: TypeTextInput, Name: "field1"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty components",
			config: FormConfig{
				Title:      "Test Form",
				Components: []ComponentConfig{},
			},
			wantErr: true,
			errMsg:  "pelo menos um componente",
		},
		{
			name: "duplicate names",
			config: FormConfig{
				Components: []ComponentConfig{
					{Type: TypeTextInput, Name: "field1"},
					{Type: TypeTextInput, Name: "field1"},
				},
			},
			wantErr: true,
			errMsg:  "duplicado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadFormConfig(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	validYAML := `
title: "Test Form"
components:
  - type: textinput
    name: username
    required: true
`
	validPath := filepath.Join(tmpDir, "valid.yaml")
	require.NoError(t, os.WriteFile(validPath, []byte(validYAML), 0600))

	invalidPath := filepath.Join(tmpDir, "nonexistent.yaml")

	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid config",
			path:    validPath,
			wantErr: false,
		},
		{
			name:    "file not found",
			path:    invalidPath,
			wantErr: true,
			errMsg:  "erro ao ler o arquivo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadFormConfig(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, cfg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.Equal(t, "Test Form", cfg.Title)
				assert.Len(t, cfg.Components, 1)
			}
		})
	}
}

func TestLoadFormConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	invalidYAML := `
title: "Invalid Form"
components:
	 - type: textinput
	   name: username
	   invalid_field: true  # YAML inválido
`
	invalidPath := filepath.Join(tmpDir, "invalid.yaml")
	require.NoError(t, os.WriteFile(invalidPath, []byte(invalidYAML), 0600))

	cfg, err := LoadFormConfig(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao analisar o YAML")
	assert.Nil(t, cfg)
}

func TestLoadLayoutConfig_FileNotFound(t *testing.T) {
	invalidPath := filepath.Join(t.TempDir(), "nonexistent.yaml")

	cfg, err := LoadLayoutConfig(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao ler o arquivo")
	assert.Nil(t, cfg)
}

func TestLoadLayoutConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	invalidYAML := `
title: "Invalid Layout"
layout: horizontal
components:
	 - type: textinput
	   name: field1
	   invalid_field: true
`
	invalidPath := filepath.Join(tmpDir, "invalid.yaml")
	require.NoError(t, os.WriteFile(invalidPath, []byte(invalidYAML), 0600))

	cfg, err := LoadLayoutConfig(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao analisar o YAML")
	assert.Nil(t, cfg)
}

func TestLoadMenuConfig_FileNotFound(t *testing.T) {
	invalidPath := filepath.Join(t.TempDir(), "nonexistent.yaml")

	cfg, err := LoadMenuConfig(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao ler o arquivo")
	assert.Nil(t, cfg)
}

func TestLoadMenuConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	invalidYAML := `
title: "Invalid Menu"
items:
	 - item1
	 - item2
invalid_field: true
`
	invalidPath := filepath.Join(tmpDir, "invalid.yaml")
	require.NoError(t, os.WriteFile(invalidPath, []byte(invalidYAML), 0600))

	cfg, err := LoadMenuConfig(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao analisar o YAML")
	assert.Nil(t, cfg)
}

func TestLoadTabsConfig_FileNotFound(t *testing.T) {
	invalidPath := filepath.Join(t.TempDir(), "nonexistent.yaml")

	cfg, err := LoadTabsConfig(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao ler o arquivo")
	assert.Nil(t, cfg)
}

func TestLoadTabsConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	invalidYAML := `
title: "Invalid Tabs"
tabs:
	 - name: tab1
	   label: Tab 1
	   components:
	     - type: textinput
	       name: field1
invalid_field: true
`
	invalidPath := filepath.Join(tmpDir, "invalid.yaml")
	require.NoError(t, os.WriteFile(invalidPath, []byte(invalidYAML), 0600))

	cfg, err := LoadTabsConfig(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao analisar o YAML")
	assert.Nil(t, cfg)
}

func TestLayoutConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  LayoutConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid horizontal layout",
			config: LayoutConfig{
				Layout: "horizontal",
				Components: []ComponentConfig{
					{Type: TypeTextInput, Name: "field1"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid vertical layout",
			config: LayoutConfig{
				Layout: "vertical",
				Components: []ComponentConfig{
					{Type: TypeTextInput, Name: "field1"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid layout type",
			config: LayoutConfig{
				Layout: "diagonal",
				Components: []ComponentConfig{
					{Type: TypeTextInput, Name: "field1"},
				},
			},
			wantErr: true,
			errMsg:  "horizontal' ou 'vertical",
		},
		{
			name: "empty components",
			config: LayoutConfig{
				Layout:     "horizontal",
				Components: []ComponentConfig{},
			},
			wantErr: true,
			errMsg:  "pelo menos um componente",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMenuConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  MenuConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid menu",
			config: MenuConfig{
				Title: "Test Menu",
				Items: []string{"Option 1", "Option 2"},
			},
			wantErr: false,
		},
		{
			name: "empty items",
			config: MenuConfig{
				Title: "Test Menu",
				Items: []string{},
			},
			wantErr: true,
			errMsg:  "pelo menos um item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTabsConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  TabsConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid tabs",
			config: TabsConfig{
				Title: "Test Tabs",
				Tabs: []TabConfig{
					{
						Name:  "tab1",
						Label: "Tab 1",
						Components: []ComponentConfig{
							{Type: TypeTextInput, Name: "field1"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty tabs",
			config: TabsConfig{
				Title: "Test Tabs",
				Tabs:  []TabConfig{},
			},
			wantErr: true,
			errMsg:  "pelo menos uma aba",
		},
		{
			name: "tab without name",
			config: TabsConfig{
				Tabs: []TabConfig{
					{
						Label: "Tab without name",
					},
				},
			},
			wantErr: true,
			errMsg:  "nome é obrigatório",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
