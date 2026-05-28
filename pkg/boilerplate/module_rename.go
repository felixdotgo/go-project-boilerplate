package boilerplate

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type RenameModuleOptions struct {
	RootDir      string
	TargetModule string
}

type RenameModuleResult struct {
	OldModule    string
	NewModule    string
	ChangedFiles []string
}

var textRewriteExts = map[string]struct{}{
	".go":    {},
	".mod":   {},
	".proto": {},
	".yaml":  {},
	".yml":   {},
	".md":    {},
}

func RenameModule(opts RenameModuleOptions) (*RenameModuleResult, error) {
	rootDir, err := findRepoRoot(opts.RootDir)
	if err != nil {
		return nil, err
	}

	currentModule, err := readModulePath(filepath.Join(rootDir, "go.mod"))
	if err != nil {
		return nil, err
	}

	targetModule := strings.TrimSpace(opts.TargetModule)
	if targetModule == "" {
		targetModule, err = inferModulePathFromOrigin(rootDir)
		if err != nil {
			return nil, err
		}
	}
	if err := validateModulePath(targetModule); err != nil {
		return nil, err
	}

	result := &RenameModuleResult{
		OldModule: currentModule,
		NewModule: targetModule,
	}
	if currentModule == targetModule {
		return result, nil
	}

	files, err := listTrackedFiles(rootDir)
	if err != nil {
		return nil, err
	}

	for _, relPath := range files {
		if !shouldRewriteFile(relPath) {
			continue
		}

		absPath := filepath.Join(rootDir, relPath)
		changed, err := rewriteModulePathInFile(absPath, currentModule, targetModule)
		if err != nil {
			return nil, err
		}
		if changed {
			result.ChangedFiles = append(result.ChangedFiles, relPath)
		}
	}

	slices.Sort(result.ChangedFiles)

	return result, nil
}

func findRepoRoot(start string) (string, error) {
	if start == "" {
		start = "."
	}

	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not locate repo root containing go.mod")
		}
		dir = parent
	}
}

func readModulePath(goModPath string) (string, error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			modulePath := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if modulePath == "" {
				break
			}
			return modulePath, nil
		}
	}

	return "", fmt.Errorf("could not find module path in %s", goModPath)
}

func inferModulePathFromOrigin(rootDir string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = rootDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not infer module path from git origin: %w", err)
	}

	modulePath, err := normalizeGitRemote(strings.TrimSpace(string(output)))
	if err != nil {
		return "", err
	}
	return modulePath, nil
}

func normalizeGitRemote(remote string) (string, error) {
	remote = strings.TrimSpace(remote)
	if remote == "" {
		return "", errors.New("git origin remote is empty")
	}

	if strings.HasPrefix(remote, "git@") {
		hostAndPath := strings.TrimPrefix(remote, "git@")
		parts := strings.SplitN(hostAndPath, ":", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("unsupported git remote: %s", remote)
		}
		return normalizeModuleParts(parts[0], parts[1])
	}

	parsed, err := url.Parse(remote)
	if err != nil {
		return "", fmt.Errorf("could not parse git remote %q: %w", remote, err)
	}
	if parsed.Host == "" || parsed.Path == "" {
		return "", fmt.Errorf("unsupported git remote: %s", remote)
	}

	return normalizeModuleParts(parsed.Host, parsed.Path)
}

func normalizeModuleParts(host, repoPath string) (string, error) {
	cleanPath := strings.TrimSpace(repoPath)
	cleanPath = strings.TrimPrefix(cleanPath, "/")
	cleanPath = strings.TrimSuffix(cleanPath, ".git")
	cleanPath = path.Clean(cleanPath)
	if cleanPath == "." || cleanPath == "" {
		return "", fmt.Errorf("unsupported repository path: %s", repoPath)
	}

	modulePath := host + "/" + cleanPath
	return modulePath, validateModulePath(modulePath)
}

func validateModulePath(modulePath string) error {
	if strings.TrimSpace(modulePath) == "" {
		return errors.New("module path must not be empty")
	}
	if strings.ContainsAny(modulePath, " \t\r\n") {
		return fmt.Errorf("module path contains whitespace: %q", modulePath)
	}
	if strings.HasPrefix(modulePath, "/") || !strings.Contains(modulePath, "/") {
		return fmt.Errorf("invalid module path: %q", modulePath)
	}
	return nil
}

func listTrackedFiles(rootDir string) ([]string, error) {
	cmd := exec.Command("git", "ls-files", "-z")
	cmd.Dir = rootDir
	output, err := cmd.Output()
	if err == nil {
		entries := bytes.Split(bytes.TrimRight(output, "\x00"), []byte{0})
		files := make([]string, 0, len(entries))
		for _, entry := range entries {
			if len(entry) == 0 {
				continue
			}
			files = append(files, string(entry))
		}
		return files, nil
	}

	var files []string
	walkErr := filepath.WalkDir(rootDir, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(rootDir, p)
		if err != nil {
			return err
		}
		files = append(files, relPath)
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	return files, nil
}

func shouldRewriteFile(relPath string) bool {
	base := filepath.Base(relPath)
	if base == "go.mod" || base == "buf.gen.yaml" {
		return true
	}

	ext := filepath.Ext(relPath)
	_, ok := textRewriteExts[ext]
	return ok
}

func rewriteModulePathInFile(filename, oldModule, newModule string) (bool, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}
	if !bytes.Contains(data, []byte(oldModule)) {
		return false, nil
	}

	updated := bytes.ReplaceAll(data, []byte(oldModule), []byte(newModule))
	info, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	if err := os.WriteFile(filename, updated, info.Mode()); err != nil {
		return false, err
	}

	return true, nil
}
