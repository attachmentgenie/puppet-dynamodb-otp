package cmd

import (
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"

	otp "github.com/attachmentgenie/puppet-dynamodb-otp/internal/aws"
	"github.com/attachmentgenie/puppet-dynamodb-otp/internal/puppet"
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
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		inputReader := cmd.InOrStdin()
		csrPEM, err := io.ReadAll(inputReader)
		if err != nil {
			log.Fatal(err)
		}
		csrCP, err := puppet.GetChallengePassword(csrPEM)
		if err != nil {
			log.Fatal(err)
		}

		fqdn := args[0]
		client, err := otp.New()
		if err != nil {
			panic(err)
		}
		otp, err := client.Read(fqdn)
		if err != nil {
			log.Fatalf("unable to find otp token for %s", fqdn)
		}

		if otp.Token_table_item == csrCP {
			fmt.Println("Found otp for " + fqdn + "")
		} else {
			log.Fatalf("Unable to match otp token for %s", fqdn)
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
