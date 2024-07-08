package cmd

import (
	"encoding/pem"
	"fmt"
	otp "github.com/attachmentgenie/puppet-dynamodb-otp/internal/aws"
	"io"
	"log"

	"github.com/micromdm/scep/v2/cryptoutil/x509util"
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
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		inputReader := cmd.InOrStdin()
		csrPEM, err := io.ReadAll(inputReader)
		if err != nil {
			log.Fatal(err)
		}
		block, _ := pem.Decode(csrPEM)
		if block == nil {
			log.Fatal("failed to decode PEM block")
		}

		// https://github.com/golang/go/issues/15995
		// https://github.com/micromdm/scep/pull/45
		//
		// The pem package is not able to parse challenge passwords yet,
		// so we need to obtain that through some parsing of our own.
		// Luckily someone did the hard work for us already
		cp, err := x509util.ParseChallengePassword(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(cp)

		fqdn := args[0]
		otp, err := otp.Read(fqdn)
		if err != nil {
			log.Fatalf("unable to find otp token for %s", fqdn)
		}

		if otp.Otp_token == cp {
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
