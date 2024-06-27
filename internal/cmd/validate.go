package cmd

import (
	"context"
	"fmt"
	"log"

	net "github.com/THREATINT/go-net"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/spf13/cobra"
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
		tableName := "puppet-dynamodb-otp"

		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
		if err != nil {
			panic(err)
		}
		svc := dynamodb.NewFromConfig(cfg)
		out, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"fqdn": &types.AttributeValueMemberS{Value: fqdn},
			},
		})
		if err != nil {
			panic(err)
		}

		otp_token := otp_token{}
		err = attributevalue.UnmarshalMap(out.Item, &otp_token)
		if err != nil {
			log.Fatalf("unmarshal failed, %v", err)
		}
		if len(out.Item) != 0 {
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
