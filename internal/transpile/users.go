package transpile

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"
)

type outputUser struct {
	Name              string   `yaml:"name"`
	PasswordHash      *string  `yaml:"password_hash,omitempty"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	Gecos             *string  `yaml:"gecos,omitempty"`
	HomeDir           *string  `yaml:"home_dir,omitempty"`
	Shell             *string  `yaml:"shell,omitempty"`
}

func parseUsers(document map[string]any, file *ast.File) ([]outputUser, ValidationErrors) {
	var problems ValidationErrors
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

func parseUser(fields map[string]any, path string, file *ast.File) (outputUser, ValidationErrors) {
	allowed := map[string]bool{
		"name": true, "passwd": true, "gecos": true, "homedir": true,
		"shell": true, "ssh_authorized_keys": true,
	}
	var problems ValidationErrors
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

func requiredString(fields map[string]any, key, path string, file *ast.File, problems *ValidationErrors) string {
	value := optionalString(fields, key, path, file, problems)
	if value == nil {
		if _, exists := fields[key]; !exists {
			*problems = append(*problems, problem(file, path+"."+key, "is required"))
		}
		return ""
	}
	return *value
}

func optionalString(fields map[string]any, key, path string, file *ast.File, problems *ValidationErrors) *string {
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

func sshKeys(fields map[string]any, path string, file *ast.File, problems *ValidationErrors) []string {
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
