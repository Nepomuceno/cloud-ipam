package main

import (
	"fmt"
	"os"

	"github.com/nepomuceno/cloud-ipam/cmd"
	"github.com/spf13/cobra/doc"
)

// go:generate go run main.go generate-docs

func main() {
	if len(os.Args) == 2 && os.Args[1] == "generate-docs" {
		fmt.Println("Generating docs...")
		doc.GenMarkdownTree(cmd.GetRootCmd(), "./docs")
	} else {
		cmd.Execute()
	}
}
