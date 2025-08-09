package cmd

import (
	"github.com/heyuuu/go-lombok/internal/lombok"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

var clearFlags struct {
	dir string
}

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear lombok code",
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := filepath.Abs(generateFlags.dir)
		if err != nil {
			log.Fatalln(err)
		}

		lombok.Clear(dir)
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)

	// Here you will define your flags and configuration settings.
	clearCmd.Flags().StringVarP(&clearFlags.dir, "dir", "d", "", "src code dir")
}
