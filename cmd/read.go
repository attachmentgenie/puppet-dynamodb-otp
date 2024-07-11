package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	otp "github.com/attachmentgenie/puppet-dynamodb-otp/internal/aws"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [FQDN]",
	Short: "List active OTP token(s).",
	Long:  "List active OTP token(s) for use in puppet auto signing ceremony.",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"FQDN", "Expires at", "OTP token"})

		client, err := otp.New()
		if err != nil {
			panic(err)
		}
		if len(args) == 1 {
			fqdn := args[0]
			otp_token, err := client.Read(fqdn)
			if err != nil {
				t.AppendRow([]interface{}{otp_token.Fqdn, time.Unix(otp_token.Expire_at_unix, 0).Format(time.Kitchen), otp_token.Otp_token})
				t.Render()
			} else {
				log.Fatalf("unable to find otp token for %s", fqdn)
			}
		} else {
			tokens := client.ReadAll()
			if len(tokens) > 0 {
				for _, record := range tokens {
					t.AppendRow([]interface{}{record.Fqdn, time.Unix(record.Expire_at_unix, 0).Format(time.Kitchen), record.Otp_token})
				}
				t.Render()
			} else {
				fmt.Println("No otp tokens found.")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
