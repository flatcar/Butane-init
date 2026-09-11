package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/flatcar/Butane-init/internal/transpile"
	"github.com/spf13/cobra"
)

var errUsage = errors.New("invalid command usage")

func main() {
	cmd := newCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		if errors.Is(err, errUsage) {
			os.Exit(2)
		}
		os.Exit(1)
	}
}

func newCommand(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	var outputPath string
	cmd := &cobra.Command{
		Use:           "bt [input-file]",
		Short:         "Transpile cloud-config YAML to Flatcar Butane YAML",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fmt.Errorf("%w: expected at most one input file", errUsage)
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			input, err := readInput(stdin, args)
			if err != nil {
				return err
			}
			output, err := transpile.Transpile(input)
			if err != nil {
				return err
			}
			if outputPath == "" || outputPath == "-" {
				_, err = stdout.Write(output)
				return err
			}
			return writeAtomically(outputPath, output)
		},
	}
	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return fmt.Errorf("%w: %v", errUsage, err)
	})
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "write output to a file instead of stdout")
	return cmd
}

func readInput(stdin io.Reader, args []string) ([]byte, error) {
	if len(args) == 0 || args[0] == "-" {
		input, err := io.ReadAll(stdin)
		if err != nil {
			return nil, fmt.Errorf("read stdin: %w", err)
		}
		return input, nil
	}
	input, err := os.ReadFile(args[0])
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", args[0], err)
	}
	return input, nil
}

func writeAtomically(path string, contents []byte) error {
	dir := filepath.Dir(path)
	temp, err := os.CreateTemp(dir, ".bt-*")
	if err != nil {
		return fmt.Errorf("create temporary output: %w", err)
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)

	if err := temp.Chmod(0o644); err != nil {
		temp.Close()
		return fmt.Errorf("set output permissions: %w", err)
	}
	if _, err := temp.Write(contents); err != nil {
		temp.Close()
		return fmt.Errorf("write output: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close output: %w", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace %q: %w", path, err)
	}
	return nil
}
