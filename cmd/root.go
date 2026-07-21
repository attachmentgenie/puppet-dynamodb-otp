package cmd

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/attachmentgenie/puppet-dynamodb-otp/internal/store"
)

var (
	TableName string
	Timeout   time.Duration
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:           "puppet-dynamodb-otp",
	Short:         "Validate puppet client CSRs.",
	Long:          `Manipulate OTP tokens for use in puppet auto signing ceremony.`,
	Args:          cobra.ExactArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
}

// GetCommandContext returns a context with timeout based on the --timeout flag.
func GetCommandContext(cmd *cobra.Command) (context.Context, context.CancelFunc) {
	parentCtx := cmd.Context()
	if parentCtx == nil {
		parentCtx = context.Background()
	}
	if Timeout > 0 {
		return context.WithTimeout(parentCtx, Timeout)
	}
	return parentCtx, func() {}
}

func Execute(defCmd string) error {
	cmd, _, err := RootCmd.Find(os.Args[1:])
	if len(os.Args[1:]) == 1 && err == nil && cmd.Use == RootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{defCmd}, os.Args[1:]...)
		RootCmd.SetArgs(args)
	}

	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&TableName, "table-name", "t", store.GetDefaultTableName(), "DynamoDB table name (env: DYNAMODB_TABLE_NAME)")
	RootCmd.PersistentFlags().DurationVar(&Timeout, "timeout", 30*time.Second, "Timeout duration for AWS operations (e.g. 10s, 30s)")
}
