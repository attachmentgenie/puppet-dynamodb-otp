package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	otp "github.com/attachmentgenie/puppet-dynamodb-otp/internal/aws"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete FQDN",
	Short: "Delete an OTP token.",
	Long:  "Delete an OTP token for use in puppet auto signing ceremony.",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fqdn := args[0]

		otp.Delete(fqdn)
		fmt.Println("Successfully deleted otp for " + fqdn + "")
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
