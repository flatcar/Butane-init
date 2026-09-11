package transpile_test

import (
	"strings"
	"testing"

	"github.com/flatcar/Butane-init/internal/transpile"
)

func TestTranspileClusterAPIWorkerUser(t *testing.T) {
	input := `## template: jinja
#cloud-config
users:
  - name: alice
    passwd: "$6$hash"
    gecos: Alice Example
    homedir: /srv/alice
    shell: /bin/bash
    ssh_authorized_keys:
      - ssh-ed25519 AAAAC3Nza alice@example
`

	want := `variant: flatcar
version: 1.0.0
passwd:
  users:
    - name: alice
      password_hash: "!$6$hash"
      ssh_authorized_keys:
        - ssh-ed25519 AAAAC3Nza alice@example
      gecos: Alice Example
      home_dir: /srv/alice
      shell: /bin/bash
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
		"unsupported top-level field": {
			input:   "#cloud-config\nruncmd: []\n",
			wantErr: "runcmd",
		},
		"deferred user field": {
			input:   "#cloud-config\nusers:\n  - name: alice\n    groups: wheel\n",
			wantErr: "groups",
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
