package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Builds the package and runs the resulting binary.
// Returns combined stdout+stderr and an error.
func buildAndRun(buildPkg, buildDir, runDir string, args []string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "grun-build-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	binName := "git-cc-test-bin"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(tmpDir, binName)

	buildCmd := exec.Command("go", "build", "-o", binPath, buildPkg)
	buildCmd.Dir = buildDir
	buildOut, err := buildCmd.CombinedOutput()
	if err != nil {
		return string(buildOut), err
	}

	runCmd := exec.Command(binPath, args...)
	if runDir != "" {
		runCmd.Dir = runDir
	}
	out, err := runCmd.CombinedOutput()
	return string(out), err
}

// pkgInfoFromCaller returns a build package path (relative to module) and the module root dir.
func pkgInfoFromCaller() (string, string) {
	_, file, _, _ := runtime.Caller(0)
	pkgDir := filepath.Dir(file)
	moduleRoot := filepath.Dir(filepath.Dir(pkgDir))
	return "./cmd/git-cc", moduleRoot
}

func TestUsage_NoArgs(t *testing.T) {
	buildPkg, moduleRoot := pkgInfoFromCaller()
	out, err := buildAndRun(buildPkg, moduleRoot, moduleRoot, nil)
	if err == nil {
		t.Fatalf("expected non-zero exit when no args provided; got success; output: %s", out)
	}
	if !strings.Contains(out, "Usage:") {
		t.Fatalf("expected usage text in output; got: %s", out)
	}
}

func TestHelpFlag(t *testing.T) {
	buildPkg, moduleRoot := pkgInfoFromCaller()
	out, _ := buildAndRun(buildPkg, moduleRoot, moduleRoot, []string{"-h"})
	if !strings.Contains(out, "Usage:") {
		t.Fatalf("expected usage text in help output; got: %s", out)
	}
}

func TestHelpContainsOptions(t *testing.T) {
	buildPkg, moduleRoot := pkgInfoFromCaller()
	out, _ := buildAndRun(buildPkg, moduleRoot, moduleRoot, []string{"-h"})

	wantFlags := []string{"-reset", "-dry-run", "-json", "-only-destructive", "-confirm", "-remote"}
	for _, f := range wantFlags {
		if !strings.Contains(out, f) {
			t.Fatalf("help output missing flag %q; full help:\n%s", f, out)
		}
	}
}

func TestJSONErrorInNonGitRepo(t *testing.T) {
	buildPkg, moduleRoot := pkgInfoFromCaller()

	tmpDir, err := os.MkdirTemp("", "git-cc-test-non-git-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	out, runErr := buildAndRun(buildPkg, moduleRoot, tmpDir, []string{"-json", "nonexistent-branch"})
	if runErr == nil {
		t.Fatalf("expected the command to fail when run outside a git repo; output: %s", out)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("expected JSON output on error when -json is provided; unmarshal error: %v; raw output: %s", err, out)
	}
	if _, ok := parsed["error"]; !ok {
		t.Fatalf("expected JSON output to contain \"error\" field; got: %v", parsed)
	}
}

func TestJSONErrorContainsMessage(t *testing.T) {
	buildPkg, moduleRoot := pkgInfoFromCaller()

	tmpDir, err := os.MkdirTemp("", "git-cc-test-non-git-2-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	out, runErr := buildAndRun(buildPkg, moduleRoot, tmpDir, []string{"-json", "nonexistent-branch"})
	if runErr == nil {
		t.Fatalf("expected non-zero exit; output: %s", out)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("failed to parse json error: %v; raw: %s", err, out)
	}
	errVal, _ := parsed["error"].(string)
	if errVal == "" {
		t.Fatalf("expected non-empty error message in JSON output; got: %v", parsed)
	}
	if !strings.Contains(strings.ToLower(errVal), "git") && !strings.Contains(strings.ToLower(errVal), "branch") {
		t.Fatalf("unexpected error message content: %s", errVal)
	}
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
	return string(out)
}

func createTestRepo(t *testing.T) (string, func()) {
	t.Helper()
	repoDir, err := os.MkdirTemp("", "gr-test-repo-")
	if err != nil {
		t.Fatalf("failed to create temp repo dir: %v", err)
	}

	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "test@example.com")
	runGit(t, repoDir, "config", "user.name", "Test User")

	if err := os.WriteFile(filepath.Join(repoDir, "README.md"), []byte("# test\n"), 0o644); err != nil {
		t.Fatalf("write readme: %v", err)
	}
	runGit(t, repoDir, "add", ".")
	runGit(t, repoDir, "commit", "-m", "init")

	runGit(t, repoDir, "checkout", "-b", "keep")
	runGit(t, repoDir, "checkout", "-b", "feature")

	bareDir, err := os.MkdirTemp("", "gr-test-bare-")
	if err != nil {
		os.RemoveAll(repoDir)
		t.Fatalf("failed to create bare dir: %v", err)
	}
	runGit(t, bareDir, "init", "--bare")
	runGit(t, repoDir, "remote", "add", "origin", bareDir)
	runGit(t, repoDir, "push", "--all", "origin")

	cleanup := func() {
		os.RemoveAll(repoDir)
		os.RemoveAll(bareDir)
	}
	return repoDir, cleanup
}

func TestDryRunSuccessLocalRepo(t *testing.T) {
	buildPkg, moduleRoot := pkgInfoFromCaller()
	repoDir, cleanup := createTestRepo(t)
	defer cleanup()

	out, err := buildAndRun(buildPkg, moduleRoot, repoDir, []string{"-dry-run", "feature"})
	if err != nil {
		t.Fatalf("expected success running dry-run in test repo; err: %v; out: %s", err, out)
	}
	if !strings.Contains(out, "completed successfully") && !strings.Contains(out, "✔ git-cc completed successfully") {
		t.Fatalf("expected success message for dry-run; got: %s", out)
	}
}

func TestJSONDryRunAndResetFlags(t *testing.T) {
	buildPkg, moduleRoot := pkgInfoFromCaller()
	repoDir, cleanup := createTestRepo(t)
	defer cleanup()

	out, err := buildAndRun(buildPkg, moduleRoot, repoDir, []string{"-json", "-dry-run", "-reset", "feature"})
	if err != nil {
		t.Fatalf("expected success running -json -dry-run -reset; err: %v; out: %s", err, out)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("expected JSON output; unmarshal error: %v; raw: %s", err, out)
	}
	// Expect top-level keys
	if _, ok := parsed["plan"]; !ok {
		t.Fatalf("json output missing plan: %v", parsed)
	}
	if _, ok := parsed["result"]; !ok {
		t.Fatalf("json output missing result: %v", parsed)
	}
	if success, ok := parsed["success"].(bool); !ok || !success {
		t.Fatalf("expected success=true in json output; got: %v", parsed["success"])
	}
}

func TestOnlyDestructiveFlagWithDryRun(t *testing.T) {
	buildPkg, moduleRoot := pkgInfoFromCaller()
	repoDir, cleanup := createTestRepo(t)
	defer cleanup()

	out, err := buildAndRun(buildPkg, moduleRoot, repoDir, []string{"-json", "-only-destructive", "-dry-run", "feature"})
	if err != nil {
		t.Fatalf("expected success running -only-destructive -dry-run; err: %v; out: %s", err, out)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("expected JSON output; unmarshal error: %v; raw: %s", err, out)
	}
	if success, ok := parsed["success"].(bool); !ok || !success {
		t.Fatalf("expected success=true in json output; got: %v", parsed["success"])
	}
}
