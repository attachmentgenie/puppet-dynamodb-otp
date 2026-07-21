package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCmd(t *testing.T) {
	SetVersionInfo("1.2.3", "abc1234", "2026-07-21")

	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetArgs([]string{"version"})

	err := RootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing version command: %v", err)
	}

	got := buf.String()
	want := "puppet-dynamodb-otp 1.2.3, commit abc1234, built at 2026-07-21"
	if !strings.Contains(got, want) {
		t.Errorf("version output = %q, want containing %q", got, want)
	}
}

func TestFixArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Strip validate-csr subcommand prefix",
			input:    []string{"validate-csr", "node.example.com"},
			expected: []string{"node.example.com"},
		},
		{
			name:     "Leave standard arguments untouched",
			input:    []string{"node.example.com"},
			expected: []string{"node.example.com"},
		},
		{
			name:     "Empty args",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fixArgs(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("fixArgs() len = %d, expected %d", len(got), len(tt.expected))
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("fixArgs()[%d] = %q, expected %q", i, got[i], tt.expected[i])
				}
			}
		})
	}
}
