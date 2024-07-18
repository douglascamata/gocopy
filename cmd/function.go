package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/douglascamata/gocopy/pkg/gocopy"
)

func newFunctionCommand() *cobra.Command {
	var functionName string
	functionCmd := &cobra.Command{
		Use:   "function",
		Short: "Copy over a function from a Go source file.",
		Long: `Copy over a function from a Go source file.
Examples:

To copy the Decode function from the encoding/json package:
> gocopy function -u "https://golang.org/src/encoding/json/decode.go?m=text" -n Unmarshal
`,
		Run: func(cmd *cobra.Command, args []string) {
			url := cmd.Flag("url").Value.String()
			out, err := gocopy.CopyFunction(gocopy.HTTPFetcher{URL: url}, functionName)
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), err)
			}
			if out != "" {
				fmt.Fprintln(cmd.OutOrStdout(), out)
			}
		},
	}
	functionCmd.Flags().StringVarP(&functionName, "name", "n", "", "The name of the function to copy")
	return functionCmd
}
