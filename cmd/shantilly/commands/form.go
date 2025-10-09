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

var formCmd = &cobra.Command{
	Use:   "form [config.yaml]",
	Short: "Executa uma TUI de formulário interativo",
	Long: `Carrega um arquivo de configuração YAML e executa uma TUI de formulário
interativo. O resultado é serializado em JSON.`,
	Args: cobra.ExactArgs(1),
	RunE: runForm,
}

func runForm(cmd *cobra.Command, args []string) error {
	start := time.Now()
	log.Printf("[DEBUG] Iniciando execução do comando form - arquivo: %s", args[0])

	configPath := args[0]

	// Load configuration with explicit error handling
	log.Printf("[DEBUG] Carregando configuração do arquivo: %s", configPath)
	cfg, err := config.LoadFormConfig(configPath)
	if err != nil {
		log.Printf("[ERROR] Falha ao carregar configuração após %v: %v", time.Since(start), err)
		return fmt.Errorf("erro ao carregar configuração: %w", err)
	}
	log.Printf("[DEBUG] Configuração carregada com sucesso em %v", time.Since(start))

	// Create theme
	log.Printf("[DEBUG] Criando tema padrão")
	theme := styles.DefaultTheme()
	log.Printf("[DEBUG] Tema criado em %v", time.Since(start))

	// Create form model
	log.Printf("[DEBUG] Criando modelo do formulário")
	model, err := models.NewFormModel(cfg, theme)
	if err != nil {
		log.Printf("[ERROR] Falha ao criar modelo após %v: %v", time.Since(start), err)
		return fmt.Errorf("erro ao criar modelo do formulário: %w", err)
	}
	log.Printf("[DEBUG] Modelo criado com sucesso em %v", time.Since(start))

	// Create and run tea program
	log.Printf("[DEBUG] Criando programa tea.NewProgram")

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

	finalModel, err := p.Run()
	if err != nil {
		log.Printf("[ERROR] Falha na execução da TUI após %v: %v", time.Since(start), err)
		return fmt.Errorf("erro ao executar TUI: %w", err)
	}
	log.Printf("[DEBUG] TUI executada com sucesso em %v", time.Since(start))

	// Check if form was submitted
	log.Printf("[DEBUG] Verificando modelo final após %v", time.Since(start))
	formModel, ok := finalModel.(*models.FormModel)
	if !ok {
		log.Printf("[ERROR] Tipo de modelo inválido após %v", time.Since(start))
		return fmt.Errorf("erro interno: tipo de modelo inválido")
	}

	log.Printf("[DEBUG] Modelo verificado, status de submissão: %v (tempo: %v)", formModel.Submitted(), time.Since(start))

	if formModel.Submitted() {
		// Serialize to JSON
		log.Printf("[DEBUG] Serializando dados para JSON")
		jsonData, err := formModel.ToJSON()
		if err != nil {
			log.Printf("[ERROR] Falha na serialização após %v: %v", time.Since(start), err)
			return fmt.Errorf("erro ao serializar dados: %w", err)
		}
		log.Printf("[DEBUG] Dados serializados em %v", time.Since(start))

		// Write to stdout
		log.Printf("[DEBUG] Escrevendo saída JSON no stdout")
		if _, err := fmt.Fprintln(os.Stdout, string(jsonData)); err != nil {
			log.Printf("[ERROR] Falha na escrita do stdout após %v: %v", time.Since(start), err)
			return fmt.Errorf("erro ao escrever saída JSON no stdout: %w", err)
		}
		log.Printf("[DEBUG] Comando form concluído com sucesso em %v", time.Since(start))
	}

	return nil
}
