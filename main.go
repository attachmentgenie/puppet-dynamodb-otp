package main

import (
	"github.com/attachmentgenie/puppet-dynamodb-otp/internal/cmd"
)

func main() {
	cmd.Execute("validate-csr")
}
