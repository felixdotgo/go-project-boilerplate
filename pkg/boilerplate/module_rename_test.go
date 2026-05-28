package boilerplate

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestNormalizeGitRemote(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"git@github.com:felixdotgo/go-project-boilerplate.git":       "github.com/felixdotgo/go-project-boilerplate",
		"https://github.com/felixdotgo/go-project-boilerplate.git":   "github.com/felixdotgo/go-project-boilerplate",
		"ssh://git@github.com/felixdotgo/go-project-boilerplate.git": "github.com/felixdotgo/go-project-boilerplate",
	}

	for input, want := range tests {
		input := input
		want := want
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeGitRemote(input)
			if err != nil {
				t.Fatalf("normalizeGitRemote() error = %v", err)
			}
			if got != want {
				t.Fatalf("normalizeGitRemote() = %q, want %q", got, want)
			}
		})
	}
}

func TestRenameModuleExplicitTarget(t *testing.T) {
	rootDir := t.TempDir()

	mustWriteFile(t, filepath.Join(rootDir, "go.mod"), "module github.com/0x46656C6978/go-project-boilerplate\n\ngo 1.24.0\n")
	mustWriteFile(t, filepath.Join(rootDir, "buf.gen.yaml"), "value: github.com/0x46656C6978/go-project-boilerplate/pkg/rpc\n")
	mustWriteFile(t, filepath.Join(rootDir, "api/proto/example.proto"), "option go_package = \"github.com/0x46656C6978/go-project-boilerplate/pkg/rpc\";\n")
	mustWriteFile(t, filepath.Join(rootDir, "cmd/app/main.go"), "package main\n\nimport _ \"github.com/0x46656C6978/go-project-boilerplate/pkg/log\"\n")
	mustWriteFile(t, filepath.Join(rootDir, "README.md"), "github.com/0x46656C6978/go-project-boilerplate\n")
	mustWriteFile(t, filepath.Join(rootDir, "binary.ico"), "github.com/0x46656C6978/go-project-boilerplate\n")

	runGit(t, rootDir, "init")
	runGit(t, rootDir, "add", "go.mod", "buf.gen.yaml", "api/proto/example.proto", "cmd/app/main.go", "README.md", "binary.ico")

	result, err := RenameModule(RenameModuleOptions{
		RootDir:      rootDir,
		TargetModule: "github.com/example/forked-repo",
	})
	if err != nil {
		t.Fatalf("RenameModule() error = %v", err)
	}

	wantChanged := []string{"README.md", "api/proto/example.proto", "buf.gen.yaml", "cmd/app/main.go", "go.mod"}
	if !slices.Equal(result.ChangedFiles, wantChanged) {
		t.Fatalf("changed files = %v, want %v", result.ChangedFiles, wantChanged)
	}

	assertContains(t, filepath.Join(rootDir, "go.mod"), "module github.com/example/forked-repo")
	assertContains(t, filepath.Join(rootDir, "buf.gen.yaml"), "github.com/example/forked-repo/pkg/rpc")
	assertContains(t, filepath.Join(rootDir, "api/proto/example.proto"), "github.com/example/forked-repo/pkg/rpc")
	assertContains(t, filepath.Join(rootDir, "cmd/app/main.go"), "github.com/example/forked-repo/pkg/log")
	assertContains(t, filepath.Join(rootDir, "README.md"), "github.com/example/forked-repo")
	assertContains(t, filepath.Join(rootDir, "binary.ico"), "github.com/0x46656C6978/go-project-boilerplate")
}

func TestRenameModuleInfersOriginRemote(t *testing.T) {
	rootDir := t.TempDir()

	mustWriteFile(t, filepath.Join(rootDir, "go.mod"), "module github.com/0x46656C6978/go-project-boilerplate\n\ngo 1.24.0\n")

	runGit(t, rootDir, "init")
	runGit(t, rootDir, "remote", "add", "origin", "git@github.com:felixdotgo/go-project-boilerplate.git")
	runGit(t, rootDir, "add", "go.mod")

	result, err := RenameModule(RenameModuleOptions{RootDir: rootDir})
	if err != nil {
		t.Fatalf("RenameModule() error = %v", err)
	}

	if result.NewModule != "github.com/felixdotgo/go-project-boilerplate" {
		t.Fatalf("new module = %q", result.NewModule)
	}
	assertContains(t, filepath.Join(rootDir, "go.mod"), "module github.com/felixdotgo/go-project-boilerplate")
}

func mustWriteFile(t *testing.T, filename, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", filename, err)
	}
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", filename, err)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v error = %v\n%s", args, err, output)
	}
}

func assertContains(t *testing.T, filename, want string) {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", filename, err)
	}
	if !strings.Contains(string(data), want) {
		t.Fatalf("%s does not contain %q; got %q", filename, want, string(data))
	}
}
