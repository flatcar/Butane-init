// Package transpile converts the supported cloud-config subset to Butane.
package transpile

import (
	"fmt"
	"sort"
	"strings"

	butane "github.com/coreos/butane/config"
	butanecommon "github.com/coreos/butane/config/common"
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

type outputUser struct {
	Name              string   `yaml:"name"`
	PasswordHash      *string  `yaml:"password_hash,omitempty"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	Gecos             *string  `yaml:"gecos,omitempty"`
	HomeDir           *string  `yaml:"home_dir,omitempty"`
	Shell             *string  `yaml:"shell,omitempty"`
}

type validationErrors []string

func (e validationErrors) Error() string {
	return strings.Join(e, "\n")
}

// Transpile converts one cloud-config document into a Flatcar Butane config.
func Transpile(input []byte) ([]byte, error) {
	if !hasCloudConfigHeader(input) {
		return nil, fmt.Errorf("missing #cloud-config header")
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

	users, problems := parseUsers(document, file)
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

func hasCloudConfigHeader(input []byte) bool {
	for _, line := range strings.Split(string(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "#cloud-config" {
			return true
		}
		if line != "" && line != "## template: jinja" {
			return false
		}
	}
	return false
}

func rejectYAMLReferences(input []byte) error {
	for _, tok := range lexer.Tokenize(string(input)) {
		if tok.Type == token.AnchorType || tok.Type == token.AliasType || tok.Type == token.MergeKeyType {
			return fmt.Errorf("[%d:%d] YAML aliases and anchors are unsupported", tok.Position.Line, tok.Position.Column)
		}
	}
	return nil
}

func parseUsers(document map[string]any, file *ast.File) ([]outputUser, validationErrors) {
	var problems validationErrors
	for _, key := range sortedKeys(document) {
		if key != "users" {
			problems = append(problems, problem(file, key, "unsupported field"))
		}
	}

	rawUsers, exists := document["users"]
	if !exists {
		return nil, problems
	}
	items, ok := rawUsers.([]any)
	if !ok {
		problems = append(problems, problem(file, "users", "must be a list"))
		return nil, problems
	}

	users := make([]outputUser, 0, len(items))
	seenUsers := map[string]struct{}{}
	for i, item := range items {
		path := fmt.Sprintf("users[%d]", i)
		fields, ok := item.(map[string]any)
		if !ok {
			problems = append(problems, problem(file, path, "must be an object"))
			continue
		}

		user, userProblems := parseUser(fields, path, file)
		problems = append(problems, userProblems...)
		if user.Name != "" {
			if _, duplicate := seenUsers[user.Name]; duplicate {
				problems = append(problems, problem(file, path+".name", fmt.Sprintf("duplicate user %q", user.Name)))
			} else {
				seenUsers[user.Name] = struct{}{}
			}
		}
		users = append(users, user)
	}
	return users, problems
}

func parseUser(fields map[string]any, path string, file *ast.File) (outputUser, validationErrors) {
	allowed := map[string]bool{
		"name": true, "passwd": true, "gecos": true, "homedir": true,
		"shell": true, "ssh_authorized_keys": true,
	}
	var problems validationErrors
	for _, key := range sortedKeys(fields) {
		if !allowed[key] {
			problems = append(problems, problem(file, path+"."+key, "unsupported field"))
		}
	}

	var user outputUser
	user.Name = requiredString(fields, "name", path, file, &problems)
	user.Gecos = optionalString(fields, "gecos", path, file, &problems)
	user.HomeDir = optionalString(fields, "homedir", path, file, &problems)
	user.Shell = optionalString(fields, "shell", path, file, &problems)
	if passwd := optionalString(fields, "passwd", path, file, &problems); passwd != nil {
		locked := lockPassword(*passwd)
		user.PasswordHash = &locked
	}
	user.SSHAuthorizedKeys = sshKeys(fields, path, file, &problems)
	return user, problems
}

func requiredString(fields map[string]any, key, path string, file *ast.File, problems *validationErrors) string {
	value := optionalString(fields, key, path, file, problems)
	if value == nil {
		if _, exists := fields[key]; !exists {
			*problems = append(*problems, problem(file, path+"."+key, "is required"))
		}
		return ""
	}
	return *value
}

func optionalString(fields map[string]any, key, path string, file *ast.File, problems *validationErrors) *string {
	raw, exists := fields[key]
	if !exists {
		return nil
	}
	value, ok := raw.(string)
	if !ok {
		*problems = append(*problems, problem(file, path+"."+key, "must be a string"))
		return nil
	}
	if value == "" {
		*problems = append(*problems, problem(file, path+"."+key, "must not be empty"))
		return nil
	}
	return &value
}

func sshKeys(fields map[string]any, path string, file *ast.File, problems *validationErrors) []string {
	raw, exists := fields["ssh_authorized_keys"]
	if !exists {
		return nil
	}
	items, ok := raw.([]any)
	if !ok {
		*problems = append(*problems, problem(file, path+".ssh_authorized_keys", "must be a list"))
		return nil
	}
	if len(items) == 0 {
		*problems = append(*problems, problem(file, path+".ssh_authorized_keys", "must not be empty"))
		return nil
	}

	keys := make([]string, 0, len(items))
	seen := map[string]struct{}{}
	for i, rawKey := range items {
		keyPath := fmt.Sprintf("%s.ssh_authorized_keys[%d]", path, i)
		key, ok := rawKey.(string)
		if !ok || key == "" {
			*problems = append(*problems, problem(file, keyPath, "must be a non-empty string"))
			continue
		}
		if _, duplicate := seen[key]; duplicate {
			*problems = append(*problems, problem(file, keyPath, "duplicate SSH authorized key"))
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	return keys
}

func lockPassword(hash string) string {
	if strings.HasPrefix(hash, "!") || strings.HasPrefix(hash, "*") {
		return hash
	}
	return "!" + hash
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
