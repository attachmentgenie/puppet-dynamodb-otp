package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// SetVersionInfo sets the build version metadata.
func SetVersionInfo(ver, gitCommit, buildDate string) {
	if ver != "" {
		Version = ver
	}
	if gitCommit != "" {
		Commit = gitCommit
	}
	if buildDate != "" {
		Date = buildDate
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the version.",
	Long:  `Return the version for this application.`,
	Run: func(c *cobra.Command, args []string) {
		c.Printf("puppet-dynamodb-otp %s, commit %s, built at %s\n", Version, Commit, Date)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
