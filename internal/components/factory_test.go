package components

import (
	"fmt"
	"testing"

	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewComponent(t *testing.T) {
	theme := styles.DefaultTheme()

	tests := []struct {
		name         string
		cfg          config.ComponentConfig
		expectError  bool
		expectedType string
	}{
		{
			name: "create textinput component",
			cfg: config.ComponentConfig{
				Name: "username",
				Type: config.TypeTextInput,
			},
			expectError:  false,
			expectedType: "*components.TextInput",
		},
		{
			name: "create textarea component",
			cfg: config.ComponentConfig{
				Name: "description",
				Type: config.TypeTextArea,
			},
			expectError:  false,
			expectedType: "*components.TextArea",
		},
		{
			name: "create checkbox component",
			cfg: config.ComponentConfig{
				Name: "agree",
				Type: config.TypeCheckbox,
			},
			expectError:  false,
			expectedType: "*components.Checkbox",
		},
		{
			name: "create radiogroup component",
			cfg: config.ComponentConfig{
				Name: "gender",
				Type: config.TypeRadioGroup,
				Options: map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{"id": "male", "label": "Male"},
						map[string]interface{}{"id": "female", "label": "Female"},
					},
				},
			},
			expectError:  false,
			expectedType: "*components.RadioGroup",
		},
		{
			name: "create slider component",
			cfg: config.ComponentConfig{
				Name: "volume",
				Type: config.TypeSlider,
			},
			expectError:  false,
			expectedType: "*components.Slider",
		},
		{
			name: "create textlabel component",
			cfg: config.ComponentConfig{
				Name: "title",
				Type: config.TypeText,
			},
			expectError:  false,
			expectedType: "*components.TextLabel",
		},
		{
			name: "unsupported component type",
			cfg: config.ComponentConfig{
				Name: "invalid",
				Type: "unsupported",
			},
			expectError: true,
		},
		{
			name: "invalid configuration - empty name",
			cfg: config.ComponentConfig{
				Name: "",
				Type: config.TypeTextInput,
			},
			expectError: true,
		},
		{
			name: "invalid configuration - empty type",
			cfg: config.ComponentConfig{
				Name: "test",
				Type: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			component, err := NewComponent(tt.cfg, theme)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, component)
			} else {
				require.NoError(t, err)
				require.NotNil(t, component)

				// Verify component implements the interface correctly
				assert.NotEmpty(t, component.Name())
				assert.IsType(t, true, component.IsValid()) // Just checking it returns a boolean

				// Check specific type if provided
				if tt.expectedType != "" {
					componentType := fmt.Sprintf("%T", component)
					assert.Equal(t, tt.expectedType, componentType)
				}
			}
		})
	}
}

func TestNewComponent_ComponentSpecificValidation(t *testing.T) {
	theme := styles.DefaultTheme()

	t.Run("radiogroup without items should fail", func(t *testing.T) {
		cfg := config.ComponentConfig{
			Name: "test-radio",
			Type: config.TypeRadioGroup,
			// No items provided
		}

		component, err := NewComponent(cfg, theme)
		assert.Error(t, err)
		assert.Nil(t, component)
		assert.Contains(t, err.Error(), "radiogroup deve conter pelo menos um item")
	})

	t.Run("slider with invalid range should fail", func(t *testing.T) {
		cfg := config.ComponentConfig{
			Name: "test-slider",
			Type: config.TypeSlider,
			Options: map[string]interface{}{
				"min": 100.0,
				"max": 50.0, // min > max
			},
		}

		component, err := NewComponent(cfg, theme)
		assert.Error(t, err)
		assert.Nil(t, component)
		assert.Contains(t, err.Error(), "min deve ser menor que max")
	})
}

func TestNewComponents(t *testing.T) {
	theme := styles.DefaultTheme()

	t.Run("create multiple valid components", func(t *testing.T) {
		configs := []config.ComponentConfig{
			{
				Name: "username",
				Type: config.TypeTextInput,
			},
			{
				Name: "bio",
				Type: config.TypeTextArea,
			},
			{
				Name: "newsletter",
				Type: config.TypeCheckbox,
			},
			{
				Name:  "title",
				Type:  config.TypeText,
				Label: "User Information",
			},
		}

		components, err := NewComponents(configs, theme)
		require.NoError(t, err)
		require.Len(t, components, 4)

		// Verify each component
		assert.Equal(t, "username", components[0].Name())
		assert.Equal(t, "bio", components[1].Name())
		assert.Equal(t, "newsletter", components[2].Name())
		assert.Equal(t, "title", components[3].Name())

		// Verify types
		assert.IsType(t, &TextInput{}, components[0])
		assert.IsType(t, &TextArea{}, components[1])
		assert.IsType(t, &Checkbox{}, components[2])
		assert.IsType(t, &TextLabel{}, components[3])
	})

	t.Run("empty config slice should return empty components", func(t *testing.T) {
		configs := []config.ComponentConfig{}

		components, err := NewComponents(configs, theme)
		require.NoError(t, err)
		assert.Empty(t, components)
	})

	t.Run("single invalid component should fail entire batch", func(t *testing.T) {
		configs := []config.ComponentConfig{
			{
				Name: "valid1",
				Type: config.TypeTextInput,
			},
			{
				Name: "", // Invalid - empty name
				Type: config.TypeTextInput,
			},
			{
				Name: "valid2",
				Type: config.TypeCheckbox,
			},
		}

		components, err := NewComponents(configs, theme)
		assert.Error(t, err)
		assert.Nil(t, components)
		assert.Contains(t, err.Error(), "erro ao criar componente 1")
	})

	t.Run("component creation error should include context", func(t *testing.T) {
		configs := []config.ComponentConfig{
			{
				Name: "valid",
				Type: config.TypeTextInput,
			},
			{
				Name: "invalid-radio",
				Type: config.TypeRadioGroup,
				// No items - will cause error
			},
		}

		components, err := NewComponents(configs, theme)
		assert.Error(t, err)
		assert.Nil(t, components)
		assert.Contains(t, err.Error(), "erro ao criar componente 1 (invalid-radio)")
	})
}

func TestNewComponent_AllTypesCreation(t *testing.T) {
	theme := styles.DefaultTheme()

	// Test that all supported component types can be created
	configs := map[string]config.ComponentConfig{
		"textinput": {
			Name: "test-textinput",
			Type: config.TypeTextInput,
		},
		"textarea": {
			Name: "test-textarea",
			Type: config.TypeTextArea,
		},
		"checkbox": {
			Name: "test-checkbox",
			Type: config.TypeCheckbox,
		},
		"radiogroup": {
			Name: "test-radiogroup",
			Type: config.TypeRadioGroup,
			Options: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"id": "option1", "label": "Option 1"},
				},
			},
		},
		"slider": {
			Name: "test-slider",
			Type: config.TypeSlider,
		},
		"text": {
			Name: "test-text",
			Type: config.TypeText,
		},
	}

	for typeName, cfg := range configs {
		t.Run(typeName, func(t *testing.T) {
			component, err := NewComponent(cfg, theme)
			require.NoError(t, err, "Failed to create %s component", typeName)
			require.NotNil(t, component)

			// Test basic component interface
			assert.Equal(t, cfg.Name, component.Name())
			assert.IsType(t, true, component.CanFocus())
			assert.IsType(t, true, component.IsValid())
			assert.IsType(t, "", component.GetError())
			assert.NotNil(t, component.Value())
		})
	}
}

func TestNewComponent_ErrorMessaging(t *testing.T) {
	theme := styles.DefaultTheme()

	t.Run("configuration validation error should be wrapped", func(t *testing.T) {
		cfg := config.ComponentConfig{
			Name: "", // Will cause validation error
			Type: config.TypeTextInput,
		}

		component, err := NewComponent(cfg, theme)
		assert.Error(t, err)
		assert.Nil(t, component)
		assert.Contains(t, err.Error(), "erro de validação da configuração")
	})

	t.Run("unsupported type error should be clear", func(t *testing.T) {
		cfg := config.ComponentConfig{
			Name: "test",
			Type: "unknown-type",
		}

		component, err := NewComponent(cfg, theme)
		assert.Error(t, err)
		assert.Nil(t, component)
		assert.Contains(t, err.Error(), "tipo de componente inválido: unknown-type")
	})
}

func TestNewComponents_EdgeCases(t *testing.T) {
	theme := styles.DefaultTheme()

	t.Run("nil configs slice", func(t *testing.T) {
		components, err := NewComponents(nil, theme)
		require.NoError(t, err)
		assert.Empty(t, components)
	})

	t.Run("large number of components", func(t *testing.T) {
		configs := make([]config.ComponentConfig, 100)
		for i := 0; i < 100; i++ {
			configs[i] = config.ComponentConfig{
				Name: fmt.Sprintf("component_%d", i),
				Type: config.TypeTextInput,
			}
		}

		components, err := NewComponents(configs, theme)
		require.NoError(t, err)
		assert.Len(t, components, 100)

		// Verify first and last components
		assert.Equal(t, "component_0", components[0].Name())
		assert.Equal(t, "component_99", components[99].Name())
	})
}
