package commands

import (
	"fmt"
	"log"
	"os"
	"time"

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
	start := time.Now()
	log.Printf("[DEBUG] Iniciando execução do comando layout - arquivo: %s", args[0])

	configPath := args[0]

	// Load configuration with explicit error handling
	log.Printf("[DEBUG] Carregando configuração do layout: %s", configPath)
	cfg, err := config.LoadLayoutConfig(configPath)
	if err != nil {
		log.Printf("[ERROR] Falha ao carregar configuração do layout após %v: %v", time.Since(start), err)
		return fmt.Errorf("erro ao carregar configuração: %w", err)
	}
	log.Printf("[DEBUG] Configuração do layout carregada em %v", time.Since(start))

	// Create theme
	log.Printf("[DEBUG] Criando tema padrão")
	theme := styles.DefaultTheme()
	log.Printf("[DEBUG] Tema criado em %v", time.Since(start))

	// Create layout model
	log.Printf("[DEBUG] Criando modelo do layout")
	model, err := models.NewLayoutModel(cfg, theme)
	if err != nil {
		log.Printf("[ERROR] Falha ao criar modelo do layout após %v: %v", time.Since(start), err)
		return fmt.Errorf("erro ao criar modelo do layout: %w", err)
	}
	log.Printf("[DEBUG] Modelo do layout criado em %v", time.Since(start))

	// Create and run tea program
	log.Printf("[DEBUG] Criando programa tea.NewProgram para layout")

	// Configure program options based on environment
	var opts []tea.ProgramOption
	opts = append(opts, tea.WithAltScreen())

	// Check if we're in a non-TTY environment (CI, tests, etc.)
	// Use environment variables commonly set in CI environments
	isCI := os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" ||
		os.Getenv("GITLAB_CI") != "" || os.Getenv("TRAVIS") != "" ||
		os.Getenv("CIRCLECI") != "" || os.Getenv("JENKINS_URL") != ""

	// Also check for common test environment indicators
	isTestEnv := os.Getenv("GO_TEST_ENVIRONMENT") != "" ||
		os.Getenv("SHANTILLY_TEST") != "" ||
		len(os.Args) > 1 && (os.Args[1] == "-test.v" || os.Args[1] == "test")

	if isCI || isTestEnv {
		log.Printf("[DEBUG] Ambiente CI/teste detectado, configurando window size padrão")
		opts = append(opts, tea.WithWindowSize(80, 24))
		log.Printf("[DEBUG] Window size definido para 80x24 para ambiente: CI=%s, Test=%s", isCI, isTestEnv)
	}

	p := tea.NewProgram(model, opts...)
	log.Printf("[DEBUG] Programa criado em %v, iniciando execução", time.Since(start))

	if _, err := p.Run(); err != nil {
		log.Printf("[ERROR] Falha na execução da TUI de layout após %v: %v", time.Since(start), err)
		return fmt.Errorf("erro ao executar TUI de layout: %w", err)
	}
	log.Printf("[DEBUG] Comando layout concluído com sucesso em %v", time.Since(start))

	return nil
}
