package cmd

import (
	"context"
	"fmt"

	net "github.com/THREATINT/go-net"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an OTP token.",
	Long:  "Delete an OTP token for use in puppet auto signing ceremony.",
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

		cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
			return nil
		})
		if err != nil {
			panic(err)
		}
		svc := dynamodb.NewFromConfig(cfg)
		_, err = svc.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
			TableName: aws.String("fleet-autosign-otp"),
			Key: map[string]types.AttributeValue{
				"fqdn": &types.AttributeValueMemberS{Value: fqdn},
			},
		})
		if err != nil {
			panic(err)
		}

		fmt.Println("Successfully deleted otp for " + fqdn + "")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
