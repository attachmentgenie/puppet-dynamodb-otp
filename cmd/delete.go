package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/attachmentgenie/puppet-dynamodb-otp/internal/store"
)

var deleteCmd = &cobra.Command{
	Use:   "delete FQDN",
	Short: "Delete an OTP token.",
	Long:  "Delete an OTP token for use in puppet auto signing ceremony.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fqdn := args[0]

		ctx, cancel := GetCommandContext(cmd)
		defer cancel()

		s, err := store.NewDynamoDBStoreWithTableName(ctx, TableName)
		if err != nil {
			return fmt.Errorf("initializing storage client: %w", err)
		}

		if err := s.Delete(ctx, fqdn); err != nil {
			return err
		}

		cmd.Printf("Successfully deleted otp for %s\n", fqdn)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
