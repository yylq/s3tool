package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yylq/s3tool/pkg/s3cmd"
)

var (
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
		rootCmd.OutOrStderr().Write([]byte(err.Error()))
		os.Exit(-1)
	}
}
func init() {
	rootCmd.AddCommand(s3cmd.GetCommads()...)
	rootCmd.SetArgs(os.Args[1:])
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}
