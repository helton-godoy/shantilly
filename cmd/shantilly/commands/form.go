package commands

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/models"
	"github.com/helton/shantilly/internal/styles"
	"github.com/spf13/cobra"
)

var formCmd = &cobra.Command{
	Use:   "form [config.yaml]",
	Short: "Executa uma TUI de formulário interativo",
	Long: `Carrega um arquivo de configuração YAML e executa uma TUI de formulário
interativo. O resultado é serializado em JSON.`,
	Args: cobra.ExactArgs(1),
	RunE: runForm,
}

func runForm(cmd *cobra.Command, args []string) error {
	configPath := args[0]

	// Load configuration with explicit error handling
	cfg, err := config.LoadFormConfig(configPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	// Create theme
	theme := styles.DefaultTheme()

	// Create form model
	model, err := models.NewFormModel(cfg, theme)
	if err != nil {
		return fmt.Errorf("erro ao criar modelo do formulário: %w", err)
	}

	// Create and run tea program
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("erro ao executar TUI: %w", err)
	}

	// Check if form was submitted
	formModel, ok := finalModel.(*models.FormModel)
	if !ok {
		return fmt.Errorf("erro interno: tipo de modelo inválido")
	}

	if formModel.Submitted() {
		// Serialize to JSON
		jsonData, err := formModel.ToJSON()
		if err != nil {
			return fmt.Errorf("erro ao serializar dados: %w", err)
		}

		// Write to stdout
		if _, err := fmt.Fprintln(os.Stdout, string(jsonData)); err != nil {
			return fmt.Errorf("erro ao escrever saída JSON no stdout: %w", err)
		}
	}

	return nil
}
