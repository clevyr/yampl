package main

import (
	"github.com/clevyr/go-yampl/cmd"
	"os"
)

func main() {
	if err := cmd.Command.Execute(); err != nil {
		os.Exit(1)
	}
}
