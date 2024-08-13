package main

import (
	"os"

	"github.com/clevyr/yampl/cmd"
)

var version = "beta"

func main() {
	root := cmd.New(cmd.WithVersion(version))
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
