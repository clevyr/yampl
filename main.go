package main

import (
	"os"

	"github.com/clevyr/yampl/cmd"
	"github.com/clevyr/yampl/internal/util"
)

//nolint:gochecknoglobals
var (
	version = "next"
	commit  = ""
)

func main() {
	os.Args = util.FixStringToStringNewlines(os.Args)
	rootCmd := cmd.NewCommand(version, commit)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
