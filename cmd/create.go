package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/thanhpk/randstr"

	"github.com/attachmentgenie/puppet-dynamodb-otp/internal/store"
)

var createCmd = &cobra.Command{
	Use:   "create FQDN [flags]",
	Short: "Create an OTP token.",
	Long:  "Create an OTP token for use in puppet auto signing ceremony.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fqdn := args[0]
		ttl, _ := cmd.Flags().GetInt("ttl")
		expireAtUnix := time.Now().Add(time.Duration(ttl) * time.Second).Unix()
		otpToken := randstr.Hex(16)

		ctx, cancel := GetCommandContext(cmd)
		defer cancel()

		s, err := store.NewDynamoDBStoreWithTableName(ctx, TableName)
		if err != nil {
			return fmt.Errorf("initializing storage client: %w", err)
		}

		if err := s.Create(ctx, expireAtUnix, fqdn, otpToken); err != nil {
			return err
		}

		expireTimeStr := time.Unix(expireAtUnix, 0).Format("2006-01-02 15:04:05 MST")
		cmd.Printf("Successfully created otp for %s %s which expires at %s\n", fqdn, otpToken, expireTimeStr)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().Int("ttl", 300, "Token time to live (sec)")
}
