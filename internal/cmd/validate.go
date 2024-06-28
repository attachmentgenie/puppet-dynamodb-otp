package cmd

import (
	"fmt"
	"log"

	"github.com/THREATINT/go-net"
	"github.com/spf13/cobra"

	otp "github.com/attachmentgenie/puppet-dynamodb-otp/internal/aws"
)

// validateCsrCmd represents the validateCsr command
var validateCsrCmd = &cobra.Command{
	Use:   "validate-csr FQDN",
	Short: "Validate puppet certificate signing request.",
	Long:  "Validate puppet certificate signing request in puppet auto signing ceremony.",
	Args: func(cmd *cobra.Command, args []string) error {
		// we need to undo the trick we setup in RootCmd
		args = fixArgs(args)

		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		if net.IsFQDN(args[0]) {
			return nil
		}
		return fmt.Errorf("invalid fqdn specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		// we need to undo the trick we setup in RootCmd
		args = fixArgs(args)
		fqdn := args[0]

		_, err := otp.Read(fqdn)
		if err != nil {
			fmt.Println("Found otp for " + fqdn + "")
		} else {
			log.Fatalf("unable to find otp token for %s", fqdn)
		}
	},
}

func init() {
	RootCmd.AddCommand(validateCsrCmd)
}

// The puppet autosign config doesnt't allow for subcommands being specified
// we need to undo the trick we setup in RootCmd
func fixArgs(args []string) []string {
	if len(args) > 0 && args[0] == "validate-csr" {
		args = args[1:]
	}
	return args
}
