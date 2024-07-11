package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/thanhpk/randstr"

	otp "github.com/attachmentgenie/puppet-dynamodb-otp/internal/aws"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create FQDN [flags]",
	Short: "Create an OTP token.",
	Long:  "Create an OTP token for use in puppet auto signing ceremony.",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fqdn := args[0]
		ttl, _ := cmd.Flags().GetInt("ttl")
		expire_at_unix := time.Now().Unix() + int64(ttl)
		otp_token := randstr.Hex(16)

		client, err := otp.New()
		if err != nil {
			panic(err)
		}
		client.Create(expire_at_unix, fqdn, otp_token)
		fmt.Println("Successfully created otp for " + fqdn + " " + otp_token + " which expires at " + time.Unix(expire_at_unix, 0).Format(time.Kitchen) + "")
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().Int("ttl", 300, "Token time to live (sec)")
}
