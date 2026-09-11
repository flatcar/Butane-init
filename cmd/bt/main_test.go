package main

import (
	"bytes"
	"testing"
)

func TestCommandTranspilesStdinToStdout(t *testing.T) {
	stdin := bytes.NewBufferString("#cloud-config\nusers: []\n")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newCommand(stdin, stdout, stderr)
	cmd.SetArgs(nil)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v; stderr = %q", err, stderr.String())
	}

	want := "variant: flatcar\nversion: 1.0.0\n"
	if got := stdout.String(); got != want {
		t.Fatalf("stdout mismatch\nwant:\n%s\ngot:\n%s", want, got)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q, want empty", got)
	}
}
