package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/attachmentgenie/puppet-dynamodb-otp/internal/puppet"
	"github.com/attachmentgenie/puppet-dynamodb-otp/internal/store"
)

var validateCsrCmd = &cobra.Command{
	Use:   "validate-csr FQDN",
	Short: "Validate puppet certificate signing request.",
	Long:  "Validate puppet certificate signing request in puppet auto signing ceremony.",
	Args: func(cmd *cobra.Command, args []string) error {
		// Undo subcommand hack in RootCmd if necessary
		args = fixArgs(args)

		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		args = fixArgs(args)
		fqdn := args[0]

		inputReader := cmd.InOrStdin()
		csrPEM, err := io.ReadAll(inputReader)
		if err != nil {
			return fmt.Errorf("reading CSR PEM from input: %w", err)
		}

		csrCP, err := puppet.GetChallengePassword(csrPEM)
		if err != nil {
			return fmt.Errorf("extracting challenge password from CSR: %w", err)
		}

		ctx, cancel := GetCommandContext(cmd)
		defer cancel()

		s, err := store.NewDynamoDBStoreWithTableName(ctx, TableName)
		if err != nil {
			return fmt.Errorf("initializing storage client: %w", err)
		}

		token, err := s.Read(ctx, fqdn)
		if err != nil {
			return fmt.Errorf("unable to find OTP token for %s: %w", fqdn, err)
		}

		if token.TokenTableItem == csrCP {
			cmd.Printf("Found otp for %s\n", fqdn)
			return nil
		}

		return fmt.Errorf("unable to match OTP token for %s", fqdn)
	},
}

func init() {
	RootCmd.AddCommand(validateCsrCmd)
}

func fixArgs(args []string) []string {
	if len(args) > 0 && args[0] == "validate-csr" {
		args = args[1:]
	}
	return args
}
