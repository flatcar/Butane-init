package main

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestWriteAtomicallyCreatesPrivateOutput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "butane.yaml")
	if err := writeAtomically(path, []byte("contents")); err != nil {
		t.Fatalf("writeAtomically() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat output: %v", err)
	}
	if got, want := info.Mode().Perm(), os.FileMode(0o600); got != want {
		t.Fatalf("output permissions = %o, want %o", got, want)
	}
}
