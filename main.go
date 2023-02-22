package main

import (
	"github.com/clevyr/yampl/cmd"
	"github.com/clevyr/yampl/internal/util"
	"os"
)

func main() {
	os.Args = util.FixStringToStringNewlines(os.Args)
	if err := cmd.Command.Execute(); err != nil {
		os.Exit(1)
	}
}
