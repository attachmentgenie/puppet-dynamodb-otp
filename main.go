package main

import (
	"fmt"

	"github.com/attachmentgenie/puppet-dynamodb-otp/internal/cmd"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the version identifier.",
	Long:  `Return the version identifier for this application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("puppet-dynamodb-otp %s, commit %s, built at %s", version, commit, date)
	},
}

func main() {
	cmd.RootCmd.AddCommand(versionCmd)

	// The puppet autosign config doesnt't allow for subcommands being specified
	// so we ll forward the command ourselves
	cmd.Execute("validate-csr")
}
