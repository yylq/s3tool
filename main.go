package main

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)
var(
	rootCmd = &cobra.Command{
		Use:   "s3tool",
		Short: "s3tool operation files in jdcloud s3",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("'%s' is an invalid argument", args[0])
			}
			return nil
		},
	}

)

func main() {

	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
func init() {
	rootCmd.SetArgs(os.Args[1:])
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}