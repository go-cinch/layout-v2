package biz_test

import (
	"context"
{{- if eq .Scaffold.proto_template "full" }}
	"fmt"
	"time"
{{- end }}
	"testing"

	"{{ .Computed.module_name_final }}/internal/biz"
	"{{ .Computed.module_name_final }}/internal/tests/mock"
)

func Test{{ .Computed.service_name_capitalized }}RepoGet(t *testing.T) {
	repo := {{ .Computed.service_name_camel }}Repo(t)
	ctx := context.Background()

{{- if eq .Scaffold.proto_template "full" }}
	seed := &biz.Create{{ .Computed.service_name_capitalized }}{
		Name: fmt.Sprintf("test-{{ .Computed.service_name_final }}-%d", time.Now().UnixNano()),
	}
	if err := repo.Create(ctx, seed); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := repo.Get(ctx, seed.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.ID != seed.ID {
		t.Fatalf("ID = %d, want %d", got.ID, seed.ID)
	}
	if got.Name != seed.Name {
		t.Fatalf("Name = %q, want %q", got.Name, seed.Name)
	}
{{- else }}
	got, err := repo.Get(ctx, 1)
	if err != nil {
		t.Skipf("skip: seed {{ .Computed.service_name_final }} record unavailable: %v", err)
	}
	if got.ID == 0 {
		t.Fatal("Get returned empty ID")
	}
	if got.Name == "" {
		t.Fatal("Get returned empty name")
	}
{{- end }}
}

func {{ .Computed.service_name_camel }}Repo(t *testing.T) biz.{{ .Computed.service_name_capitalized }}Repo {
	t.Helper()

	if err := mock.InitError(); err != nil {
		t.Skipf("skip: test database unavailable: %v", err)
	}

	repo := mock.{{ .Computed.service_name_capitalized }}Repo()
	if repo == nil {
		t.Skip("skip: {{ .Computed.service_name_final }} repo unavailable")
	}
	return repo
}
