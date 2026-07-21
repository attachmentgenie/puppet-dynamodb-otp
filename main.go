package main

import (
	"fmt"
	"os"

	"github.com/attachmentgenie/puppet-dynamodb-otp/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)

	// The puppet autosign config doesn't allow for subcommands being specified
	// so we forward the command ourselves.
	if err := cmd.Execute("validate-csr"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
