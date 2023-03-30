package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	net "github.com/THREATINT/go-net"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/spf13/cobra"
	"github.com/thanhpk/randstr"
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
		if net.IsFQDN(args[0]) {
			return nil
		}
		return fmt.Errorf("invalid fqdn specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		fqdn := args[0]
		ttl, _ := cmd.Flags().GetInt("ttl")
		expire_at_unix := time.Now().Unix() + int64(ttl)
		otp_token := randstr.Hex(16)

		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
		if err != nil {
			panic(err)
		}
		svc := dynamodb.NewFromConfig(cfg)
		_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String("puppet-dynamodb-otp"),
			Item: map[string]types.AttributeValue{
				"expire_at_unix": &types.AttributeValueMemberN{Value: strconv.FormatInt(expire_at_unix, 10)},
				"fqdn":           &types.AttributeValueMemberS{Value: fqdn},
				"otp_token":      &types.AttributeValueMemberS{Value: otp_token},
			},
		})

		fmt.Println("Successfully created otp for " + fqdn + " " + otp_token + " which expires at " + time.Unix(expire_at_unix, 0).Format(time.Kitchen) + "")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().Int("ttl", 300, "Token time to live (sec)")
}
