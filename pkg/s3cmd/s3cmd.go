package s3cmd

import "github.com/spf13/cobra"

func GetCommads() []*cobra.Command {
	return []*cobra.Command{objectCmd, BucketCmd, IpesCmd}
}
