package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "puppet-dynamodb-otp",
	Short: "Validate puppet client CSRs.",
	Long:  `Manipulate OTP tokens for use in puppet auto signing ceremony.`,
	Args:  cobra.ExactArgs(1),
}

func Execute(defCmd string) {
	cmd, _, err := RootCmd.Find(os.Args[1:])
	if len(os.Args[1:]) == 1 && err == nil && cmd.Use == RootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{defCmd}, os.Args[1:]...)
		RootCmd.SetArgs(args)
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {}
