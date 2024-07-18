package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	var url string
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "gocopy",
		Short: "Copy over Go code.",
		Long: `Copy over some Go code from somewhere else in the internet.
See 'gocopy help' for more information on usage.`,
	}
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "The url of the source file")
	// rootCmd.AddCommand(newFunctionCommand(gocopy.CopyFunction))
	rootCmd.AddCommand(newFunctionCommand())
	rootCmd.AddCommand(newTypeCommand())
	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := newRootCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
