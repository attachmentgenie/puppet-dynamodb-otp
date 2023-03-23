package main

import (
	"github.com/seaplane-io/fleet-autosign-otp/internal/cmd"
)

func main() {
	cmd.Execute("validate-csr")
}
