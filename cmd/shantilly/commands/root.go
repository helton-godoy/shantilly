package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "shantilly",
	Short: "Construtor de TUI declarativo via YAML",
	Long: `Shantilly é uma ferramenta CLI moderna em Go que permite criar
Interfaces de Usuário de Terminal (TUI) ricas e interativas de forma
declarativa, utilizando arquivos de configuração YAML.

Construído sobre o ecossistema Charm (Bubble Tea, Lip Gloss, Bubbles).`,
	Version: version,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão do Shantilly",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Shantilly v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(formCmd)
	rootCmd.AddCommand(layoutCmd)
	// TODO: Add menu, tabs, serve commands
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
