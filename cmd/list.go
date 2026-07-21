package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"github.com/attachmentgenie/puppet-dynamodb-otp/internal/store"
)

var listCmd = &cobra.Command{
	Use:   "list [FQDN]",
	Short: "List active OTP token(s).",
	Long:  "List active OTP token(s) for use in puppet auto signing ceremony.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"FQDN", "Expires at", "OTP token"})

		ctx, cancel := GetCommandContext(cmd)
		defer cancel()

		s, err := store.NewDynamoDBStoreWithTableName(ctx, TableName)
		if err != nil {
			return fmt.Errorf("initializing storage client: %w", err)
		}

		if len(args) == 1 {
			fqdn := args[0]
			otpToken, err := s.Read(ctx, fqdn)
			if err != nil {
				return fmt.Errorf("unable to find OTP token for %s: %w", fqdn, err)
			}

			expireTimeStr := time.Unix(otpToken.ExpireAtUnix, 0).Format("2006-01-02 15:04:05 MST")
			t.AppendRow([]interface{}{otpToken.FQDN, expireTimeStr, otpToken.TokenTableItem})
			t.Render()
		} else {
			tokens, err := s.ReadAll(ctx)
			if err != nil {
				return fmt.Errorf("fetching OTP tokens: %w", err)
			}

			if len(tokens) > 0 {
				for _, record := range tokens {
					expireTimeStr := time.Unix(record.ExpireAtUnix, 0).Format("2006-01-02 15:04:05 MST")
					t.AppendRow([]interface{}{record.FQDN, expireTimeStr, record.TokenTableItem})
				}
				t.Render()
			} else {
				cmd.Println("No otp tokens found.")
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
