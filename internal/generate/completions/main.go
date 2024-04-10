package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/clevyr/yampl/cmd"
)

func main() {
	if err := os.RemoveAll("completions"); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("completions", 0o777); err != nil {
		panic(err)
	}

	rootCmd := cmd.NewCommand()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	for _, shell := range []string{"bash", "zsh", "fish"} {
		rootCmd.SetArgs([]string{"--completion=" + shell})
		if err := rootCmd.Execute(); err != nil {
			panic(err)
		}

		f, err := os.Create(filepath.Join("completions", "yampl."+shell))
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(f, &buf); err != nil {
			panic(err)
		}

		if err := f.Close(); err != nil {
			panic(err)
		}
	}
}
