// Package transpile converts the supported cloud-config subset to Butane.
package transpile

import (
	"fmt"
	"sort"
	"strings"

	butane "github.com/coreos/ignition/v2/butane/config"
	butanecommon "github.com/coreos/ignition/v2/butane/config/common"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

const (
	variant = "flatcar"
	version = "1.0.0"
)

type outputConfig struct {
	Variant string        `yaml:"variant"`
	Version string        `yaml:"version"`
	Passwd  *outputPasswd `yaml:"passwd,omitempty"`
}

type outputPasswd struct {
	Users []outputUser `yaml:"users"`
}

type validationErrors []string

func (e validationErrors) Error() string {
	return strings.Join(e, "\n")
}

// Transpile converts one cloud-config document into a Flatcar Butane config.
func Transpile(input []byte) ([]byte, error) {
	if err := validateCloudConfigHeader(input); err != nil {
		return nil, err
	}

	file, err := parser.ParseBytes(input, 0)
	if err != nil {
		return nil, fmt.Errorf("parse cloud-config: %s", yaml.FormatError(err, false, true))
	}
	if len(file.Docs) > 1 {
		return nil, fmt.Errorf("cloud-config must contain exactly one YAML document")
	}
	if err := rejectYAMLReferences(input); err != nil {
		return nil, err
	}

	var document map[string]any
	if len(file.Docs) == 0 {
		document = map[string]any{}
	} else if err := yaml.Unmarshal(input, &document); err != nil {
		return nil, fmt.Errorf("decode cloud-config: %s", yaml.FormatError(err, false, true))
	}

	var problems validationErrors
	for _, key := range sortedKeys(document) {
		if key != "users" {
			problems = append(problems, problem(file, key, "unsupported field"))
		}
	}

	users, userProblems := parseUsers(document, file)
	problems = append(problems, userProblems...)
	if len(problems) > 0 {
		return nil, problems
	}

	config := outputConfig{Variant: variant, Version: version}
	if len(users) > 0 {
		config.Passwd = &outputPasswd{Users: users}
	}
	out, err := yaml.MarshalWithOptions(config, yaml.Indent(2), yaml.IndentSequence(true))
	if err != nil {
		return nil, fmt.Errorf("encode Butane config: %w", err)
	}
	if _, _, err := butane.TranslateBytes(out, butanecommon.TranslateBytesOptions{}); err != nil {
		return nil, fmt.Errorf("validate generated Butane config: %w", err)
	}
	return out, nil
}

func validateCloudConfigHeader(input []byte) error {
	for _, line := range strings.Split(string(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "## template: jinja" {
			return fmt.Errorf("Jinja templates are unsupported; render the template before transpiling")
		}
		if line == "#cloud-config" {
			return nil
		}
		if line != "" {
			return fmt.Errorf("missing #cloud-config header")
		}
	}
	return fmt.Errorf("missing #cloud-config header")
}

func rejectYAMLReferences(input []byte) error {
	for _, tok := range lexer.Tokenize(string(input)) {
		if tok.Type == token.AnchorType || tok.Type == token.AliasType || tok.Type == token.MergeKeyType {
			return fmt.Errorf("[%d:%d] YAML aliases and anchors are unsupported", tok.Position.Line, tok.Position.Column)
		}
	}
	return nil
}

func sortedKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func problem(file *ast.File, path, message string) string {
	locationPath := path
	for locationPath != "" {
		yamlPath, err := yaml.PathString("$." + locationPath)
		if err == nil {
			if node, err := yamlPath.FilterFile(file); err == nil && node.GetToken() != nil {
				pos := node.GetToken().Position
				return fmt.Sprintf("[%d:%d] %s: %s", pos.Line, pos.Column, path, message)
			}
		}
		separator := strings.LastIndex(locationPath, ".")
		if separator == -1 {
			break
		}
		locationPath = locationPath[:separator]
	}
	return fmt.Sprintf("%s: %s", path, message)
}
