package cmd

import (
	"github.com/heyuuu/go-lombok/internal/lombok"
	"github.com/spf13/cobra"
	"os"
)

var generateFlags struct {
	dir string
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate lombok code",
	Run: func(cmd *cobra.Command, args []string) {
		dir := generateFlags.dir
		if dir == "" {
			dir, _ = os.Getwd()
		}

		lombok.RunTask(lombok.TaskGenerate, dir)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.
	generateCmd.Flags().StringVarP(&generateFlags.dir, "dir", "d", "", "src code dir")
}
