package main

import (
	"log"
	"os"

	"github.com/clevyr/yampl/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	var err error
	output := "./docs"

	if err = os.RemoveAll(output); err != nil {
		log.Fatalf("failed to remove existing dir: %v", err)
	}

	if err = os.MkdirAll(output, 0o755); err != nil {
		log.Fatalf("failed to mkdir: %v", err)
	}

	rootCmd := cmd.NewCommand()
	if err = doc.GenMarkdownTree(rootCmd, output); err != nil {
		log.Fatalf("failed to generate markdown: %v", err)
	}
}
