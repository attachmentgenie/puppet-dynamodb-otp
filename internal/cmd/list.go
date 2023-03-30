package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	net "github.com/THREATINT/go-net"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

type otp_token struct {
	Expire_at_unix int64
	Fqdn           string
	Otp_token      string
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [FQDN]",
	Short: "List active OTP token(s).",
	Long:  "List active OTP token(s) for use in puppet auto signing ceremony.",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return err
		}
		if len(args) == 1 {
			if net.IsFQDN(args[0]) {
				return nil
			}
			return fmt.Errorf("invalid fqdn specified: %s", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tableName := "puppet-dynamodb-otp"

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"FQDN", "Expires at", "OTP token"})

		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
		if err != nil {
			panic(err)
		}
		svc := dynamodb.NewFromConfig(cfg)

		if len(args) == 1 {
			fqdn := args[0]
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
				t.AppendRow([]interface{}{otp_token.Fqdn, time.Unix(otp_token.Expire_at_unix, 0).Format(time.Kitchen), otp_token.Otp_token})
				t.Render()
			} else {
				log.Fatalf("unable to find otp token for %s", fqdn)
			}
		} else {
			out, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
				TableName: aws.String(tableName),
			})
			if err != nil {
				panic(err)
			}

			var tokens []otp_token
			err = attributevalue.UnmarshalListOfMaps(out.Items, &tokens)
			if err != nil {
				log.Fatalf("unable to unmarshal tokens: %v", err)
			}
			if out.Count > 0 {
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
	rootCmd.AddCommand(listCmd)
}
