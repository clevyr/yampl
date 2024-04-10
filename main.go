package main

import (
	"os"

	"github.com/clevyr/yampl/cmd"
	"github.com/clevyr/yampl/internal/util"
)

func main() {
	os.Args = util.FixStringToStringNewlines(os.Args)
	rootCmd := cmd.NewCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
