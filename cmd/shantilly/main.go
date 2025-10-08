package main

import (
	"fmt"
	"os"

	"github.com/helton/shantilly/cmd/shantilly/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
}
