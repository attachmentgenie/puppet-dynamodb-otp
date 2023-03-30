package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Item struct {
	Year   int
	Title  string
	Plot   string
	Rating float64
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "puppet-dynamodb-otp",
	Short: "Validate puppet client CSRs.",
	Long:  `Manipulate OTP tokens for use in puppet auto signing ceremony.`,
	Args:  cobra.ExactArgs(1),
}

func Execute(defCmd string) {
	// if no subcommand was provided forward to the validate command
	//
	// The puppet autosign config doesnt't allow for subcommands being specified
	// so we ll forward the command ourselves
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if len(os.Args[1:]) == 1 && err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{defCmd}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {}
