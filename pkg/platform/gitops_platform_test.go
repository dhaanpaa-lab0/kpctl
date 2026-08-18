package platform

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func runGitCmd(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v, output: %s", args, err, string(out))
	}
}

func setupTestRemoteRepo(t *testing.T) string {
	t.Helper()
	remoteDir := t.TempDir()

	runGitCmd(t, remoteDir, "init", "--initial-branch=main")
	runGitCmd(t, remoteDir, "config", "user.name", "Test User")
	runGitCmd(t, remoteDir, "config", "user.email", "test@example.com")

	testFile := filepath.Join(remoteDir, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test GitOps Repo\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	runGitCmd(t, remoteDir, "add", ".")
	runGitCmd(t, remoteDir, "commit", "-m", "Initial commit")

	return remoteDir
}

func TestSyncGitopsRepository_CloneAndPull(t *testing.T) {
	remoteDir := setupTestRemoteRepo(t)
	destBase := t.TempDir()
	destFolder := filepath.Join(destBase, "gitops-target")

	// 1. First sync - folder does not exist -> should create and clone
	err := SyncGitopsRepository(remoteDir, destFolder)
	if err != nil {
		t.Fatalf("SyncGitopsRepository failed on initial clone: %v", err)
	}

	readmePath := filepath.Join(destFolder, "README.md")
	data, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("failed to read cloned README.md: %v", err)
	}
	if string(data) != "# Test GitOps Repo\n" {
		t.Errorf("cloned content mismatch: got %q, want %q", string(data), "# Test GitOps Repo\n")
	}

	// 2. Add new commit in remote
	newFile := filepath.Join(remoteDir, "config.yaml")
	if err := os.WriteFile(newFile, []byte("key: value\n"), 0644); err != nil {
		t.Fatalf("failed to write new file in remote: %v", err)
	}
	runGitCmd(t, remoteDir, "add", ".")
	runGitCmd(t, remoteDir, "commit", "-m", "Add config")

	// 3. Second sync - destination already has .git -> should pull
	err = SyncGitopsRepository(remoteDir, destFolder)
	if err != nil {
		t.Fatalf("SyncGitopsRepository failed on pull: %v", err)
	}

	newFilePath := filepath.Join(destFolder, "config.yaml")
	data, err = os.ReadFile(newFilePath)
	if err != nil {
		t.Fatalf("failed to read pulled config.yaml: %v", err)
	}
	if string(data) != "key: value\n" {
		t.Errorf("pulled content mismatch: got %q, want %q", string(data), "key: value\n")
	}
}

func TestSyncGitopsRepository_TildePath(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	remoteDir := setupTestRemoteRepo(t)
	destFolder := "~/my-gitops-folder/nested"

	err := SyncGitopsRepository(remoteDir, destFolder)
	if err != nil {
		t.Fatalf("SyncGitopsRepository with tilde failed: %v", err)
	}

	expectedPath := filepath.Join(tempHome, "my-gitops-folder", "nested", "README.md")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Fatalf("expected file %s does not exist: %v", expectedPath, err)
	}
}

func TestSyncGitops_WithViperConfig(t *testing.T) {
	remoteDir := setupTestRemoteRepo(t)
	destFolder := filepath.Join(t.TempDir(), "viper-gitops")

	viper.Set("gitops_url", remoteDir)
	viper.Set("gitops_folder", destFolder)

	err := SyncGitops()
	if err != nil {
		t.Fatalf("SyncGitops failed: %v", err)
	}

	readmePath := filepath.Join(destFolder, "README.md")
	if _, err := os.Stat(readmePath); err != nil {
		t.Fatalf("expected file %s does not exist: %v", readmePath, err)
	}
}

func TestSyncGitops_MissingConfig(t *testing.T) {
	viper.Set("gitops_url", "")
	viper.Set("gitops_folder", "")

	err := SyncGitops()
	if err == nil {
		t.Fatal("expected error when config is missing, got nil")
	}
}
