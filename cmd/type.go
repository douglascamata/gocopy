package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/douglascamata/gocopy/pkg/gocopy"
)

func newTypeCommand() *cobra.Command {
	var typeName string
	var includeMethods bool

	typeCmd := &cobra.Command{
		Use:   "type",
		Short: "Copy over a type from a Go source file.",
		Long: `Copy over a type from a Go source file.
Examples:

To copy the Decoder type from the encoding/json package:

> gocopy type -u "https://golang.org/src/encoding/json/decode.go?m=text" -n Number

If you want to include all its methods, use the -m flag:
> gocopy type -u "https://golang.org/src/encoding/json/decode.go?m=text" -n Number -m
`,
		Run: func(cmd *cobra.Command, args []string) {
			url := cmd.Flag("url").Value.String()
			out, err := gocopy.CopyType(gocopy.HTTPFetcher{URL: url}, typeName, includeMethods)
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), err)
				return
			}
			if out != "" {
				fmt.Fprintln(cmd.OutOrStdout(), out)
			}
		},
	}

	typeCmd.Flags().StringVarP(&typeName, "name", "n", "", "The name of the type to copy")
	typeCmd.Flags().BoolVarP(&includeMethods, "methods", "m", false, "Include the type's methods")

	return typeCmd
}
