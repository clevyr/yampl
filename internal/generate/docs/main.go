package main

import (
	"log"
	"os"
	"strings"

	"gabe565.com/utils/cobrax"
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

	rootCmd := cmd.New(cobrax.WithVersion("beta"))
	if i := strings.Index(rootCmd.Long, "\n\nFull reference at"); i != -1 {
		rootCmd.Long = rootCmd.Long[:i]
	}
	if err = doc.GenMarkdownTree(rootCmd, output); err != nil {
		log.Fatalf("failed to generate markdown: %v", err)
	}
}
