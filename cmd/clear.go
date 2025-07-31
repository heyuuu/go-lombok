package cmd

import (
	"github.com/heyuuu/go-lombok/internal/lombok"
	"os"

	"github.com/spf13/cobra"
)

var clearFlags struct {
	dir string
}

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear lombok code",
	Run: func(cmd *cobra.Command, args []string) {
		dir := clearFlags.dir
		if dir == "" {
			dir, _ = os.Getwd()
		}

		lombok.RunTask(lombok.TaskClear, dir)
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)

	// Here you will define your flags and configuration settings.
	clearCmd.Flags().StringVarP(&clearFlags.dir, "dir", "d", "", "src code dir")
}
