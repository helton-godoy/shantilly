package commands

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/models"
	"github.com/helton/shantilly/internal/styles"
	"github.com/spf13/cobra"
)

var layoutCmd = &cobra.Command{
	Use:   "layout [config.yaml]",
	Short: "Executa uma TUI com layout estruturado",
	Long: `Carrega um arquivo de configuração YAML e executa uma TUI com layout
horizontal ou vertical.`,
	Args: cobra.ExactArgs(1),
	RunE: runLayout,
}

func runLayout(cmd *cobra.Command, args []string) error {
	configPath := args[0]

	// Load configuration with explicit error handling
	cfg, err := config.LoadLayoutConfig(configPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	// Create theme
	theme := styles.DefaultTheme()

	// Create layout model
	model, err := models.NewLayoutModel(cfg, theme)
	if err != nil {
		return fmt.Errorf("erro ao criar modelo do layout: %w", err)
	}

	// Create and run tea program
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("erro ao executar TUI de layout: %w", err)
	}

	return nil
}
