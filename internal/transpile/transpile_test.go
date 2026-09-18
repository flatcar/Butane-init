package transpile_test

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/flatcar/Butane-init/internal/transpile"
)

func TestTranspileRejectsInvalidDocument(t *testing.T) {
	tests := map[string]struct {
		input   string
		wantErr string
	}{
		"missing header": {
			input:   "users: []\n",
			wantErr: "missing #cloud-config header",
		},
		"arbitrary comment before header": {
			input:   "# generated config\n#cloud-config\nusers: []\n",
			wantErr: "missing #cloud-config header",
		},
		"jinja template": {
			input:   "## template: jinja\n#cloud-config\nusers:\n  - name: alice\n    gecos: \"{{ ds.meta_data.name }}\"\n",
			wantErr: "Jinja templates are unsupported",
		},
		"unsupported top-level field": {
			input:   "#cloud-config\nruncmd: []\n",
			wantErr: "runcmd",
		},
		"yaml alias": {
			input:   "#cloud-config\nusers:\n  - &user\n    name: alice\n  - *user\n",
			wantErr: "YAML aliases and anchors are unsupported",
		},
		"multiple documents": {
			input:   "#cloud-config\nusers: []\n---\nusers: []\n",
			wantErr: "exactly one YAML document",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := transpile.Transpile([]byte(tt.input))
			if err == nil {
				t.Fatalf("Transpile() error = nil, want an error containing %q", tt.wantErr)
			}
			if got != nil {
				t.Fatalf("Transpile() output = %q, want nil on error", got)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Transpile() error = %q, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestTranspileReturnsStructuredValidationErrors(t *testing.T) {
	input := "#cloud-config\nruncmd: []\nusers:\n  - shell: \"\"\n"

	got, err := transpile.Transpile([]byte(input))
	if err == nil {
		t.Fatal("Transpile() error = nil, want validation errors")
	}
	if got != nil {
		t.Fatalf("Transpile() output = %q, want nil on error", got)
	}

	var problems transpile.ValidationErrors
	if !errors.As(err, &problems) {
		t.Fatalf("Transpile() error type = %T, want transpile.ValidationErrors", err)
	}
	want := transpile.ValidationErrors{
		{Path: "runcmd", Message: "unsupported field", Line: 2, Column: 9},
		{Path: "users[0].name", Message: "is required", Line: 4, Column: 10},
		{Path: "users[0].shell", Message: "must not be empty", Line: 4, Column: 12},
	}
	if !reflect.DeepEqual(problems, want) {
		t.Fatalf("Transpile() validation errors = %#v, want %#v", problems, want)
	}

	wantMessage := "[2:9] runcmd: unsupported field\n" +
		"[4:10] users[0].name: is required\n" +
		"[4:12] users[0].shell: must not be empty"
	if err.Error() != wantMessage {
		t.Fatalf("Transpile() error = %q, want %q", err, wantMessage)
	}
}

func readFixture(t *testing.T, name string) string {
	t.Helper()
	contents, err := os.ReadFile("testcases/" + name)
	if err != nil {
		t.Fatalf("read fixture %q: %v", name, err)
	}
	return string(contents)
}
