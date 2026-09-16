package transpile_test

import (
	"os"
	"strings"
	"testing"

	"github.com/flatcar/Butane-init/internal/transpile"
)

func TestTranspileClusterAPIWorkerUser(t *testing.T) {
	input := readFixture(t, "cluster-api-supported-user.yaml")

	want := `variant: flatcar
version: 1.0.0
passwd:
  users:
    - name: foo
      password_hash: "!$6$REDACTED_TEST_HASH"
      ssh_authorized_keys:
        - ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIREDACTED fixture@example
      gecos: Foo B. Bar
      home_dir: /home/foo
      shell: /bin/false
`

	got, err := transpile.Transpile([]byte(input))
	if err != nil {
		t.Fatalf("Transpile() error = %v", err)
	}
	if string(got) != want {
		t.Fatalf("Transpile() output mismatch\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestTranspileRejectsInvalidOrUnsupportedInput(t *testing.T) {
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
		"Cluster API groups fixture": {
			input:   readFixture(t, "cluster-api-groups.yaml"),
			wantErr: "groups",
		},
		"Cluster API deferred fields fixture": {
			input:   readFixture(t, "cluster-api-deferred-fields.yaml"),
			wantErr: "inactive",
		},
		"empty optional field": {
			input:   "#cloud-config\nusers:\n  - name: alice\n    shell: \"\"\n",
			wantErr: "users[0].shell",
		},
		"duplicate username": {
			input:   "#cloud-config\nusers:\n  - name: alice\n  - name: alice\n",
			wantErr: "duplicate user",
		},
		"missing username has source location": {
			input:   "#cloud-config\nusers:\n  - shell: /bin/bash\n",
			wantErr: "[3:10] users[0].name: is required",
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

func readFixture(t *testing.T, name string) string {
	t.Helper()
	contents, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("read fixture %q: %v", name, err)
	}
	return string(contents)
}
