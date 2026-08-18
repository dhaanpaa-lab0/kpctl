package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func runGitCmdInCmdTest(t *testing.T, dir string, args ...string) {
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

func TestSyncGitopsCmd(t *testing.T) {
	remoteDir := t.TempDir()
	runGitCmdInCmdTest(t, remoteDir, "init", "--initial-branch=main")
	runGitCmdInCmdTest(t, remoteDir, "config", "user.name", "Test User")
	runGitCmdInCmdTest(t, remoteDir, "config", "user.email", "test@example.com")

	testFile := filepath.Join(remoteDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("sync cmd test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	runGitCmdInCmdTest(t, remoteDir, "add", ".")
	runGitCmdInCmdTest(t, remoteDir, "commit", "-m", "Initial commit")

	destFolder := filepath.Join(t.TempDir(), "cmd-dest")
	viper.Set("gitops_url", remoteDir)
	viper.Set("gitops_folder", destFolder)

	// Execute command via rootCmd
	rootCmd.SetArgs([]string{"syncGitops"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute(syncGitops) failed: %v", err)
	}

	copiedFile := filepath.Join(destFolder, "test.txt")
	if _, err := os.Stat(copiedFile); err != nil {
		t.Fatalf("expected file %s does not exist after syncGitops: %v", copiedFile, err)
	}
}
