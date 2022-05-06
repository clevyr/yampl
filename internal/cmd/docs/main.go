package main

import (
	"github.com/clevyr/go-yampl/cmd"
	"github.com/spf13/cobra/doc"
	"log"
	"os"
)

func main() {
	var err error
	output := "./docs"
	log.Println(`generating docs in "` + output + `"`)

	log.Println("removing existing directory")
	err = os.RemoveAll(output)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("making directory")
	err = os.MkdirAll(output, 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("generating markdown")
	rootCmd := cmd.Command
	err = doc.GenMarkdownTree(rootCmd, output)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("finished")
}
