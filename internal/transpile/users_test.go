package transpile_test

import (
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

func TestTranspileRejectsInvalidUsers(t *testing.T) {
	tests := map[string]struct {
		input   string
		wantErr string
	}{
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
