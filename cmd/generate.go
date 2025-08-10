package cmd

import (
	"github.com/heyuuu/go-lombok/internal/lombok"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

var generateFlags struct {
	dir      string
	excludes []string
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g", "gen"},
	Short:   "Generate lombok code",
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := filepath.Abs(generateFlags.dir)
		if err != nil {
			log.Fatalln(err)
		}

		lombok.Generate(dir, generateFlags.excludes)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.
	generateCmd.Flags().StringVarP(&generateFlags.dir, "dir", "d", "", "src code dir")
	generateCmd.Flags().StringSliceVarP(&generateFlags.excludes, "exclude", "e", nil, "exclude path")
}
