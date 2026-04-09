package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"part6/cmd"
)

func main() {
	root := &cobra.Command{
		Use:   "embed",
		Short: "Semantic code search using text embeddings",
	}

	root.AddCommand(cmd.NewSearchCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
