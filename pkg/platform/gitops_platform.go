package platform

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"nexus-sds.com/kpctl/pkg/core"
)

// SyncGitopsRepository synchronizes the git repository from gitopsURL to gitopsFolder.
// It creates gitopsFolder (and expands any ~ home path) if it does not exist.
// If the destination contains a .git repository, it runs 'git pull'.
// Otherwise, it runs 'git clone <gitopsURL> <expandedFolder>'.
func SyncGitopsRepository(gitopsURL, gitopsFolder string) error {
	if gitopsURL == "" {
		return errors.New("gitops_url is not configured")
	}
	if gitopsFolder == "" {
		return errors.New("gitops_folder is not configured")
	}

	expandedFolder, err := core.ExpandHomeDir(gitopsFolder)
	if err != nil {
		return fmt.Errorf("failed to resolve gitops folder path: %w", err)
	}

	if err := core.CreateDirectoryTree(expandedFolder); err != nil {
		return fmt.Errorf("failed to create gitops folder %s: %w", expandedFolder, err)
	}

	gitDir := filepath.Join(expandedFolder, ".git")
	if info, err := os.Stat(gitDir); err == nil && (info.IsDir() || !info.IsDir()) {
		// Existing git repository -> git pull
		cmd := exec.Command("git", "pull")
		cmd.Dir = expandedFolder
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git pull failed in %s: %w", expandedFolder, err)
		}
		return nil
	}

	// Not a git repository yet -> check if directory is empty for clone
	entries, err := os.ReadDir(expandedFolder)
	if err != nil {
		return fmt.Errorf("failed to read gitops folder %s: %w", expandedFolder, err)
	}

	if len(entries) == 0 {
		// Directory is empty -> clone directly into it
		cmd := exec.Command("git", "clone", gitopsURL, expandedFolder)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git clone %s into %s failed: %w", gitopsURL, expandedFolder, err)
		}
		return nil
	}

	// Directory is non-empty and has no .git -> init, set remote, fetch, checkout
	cmdInit := exec.Command("git", "init")
	cmdInit.Dir = expandedFolder
	cmdInit.Stdout = os.Stdout
	cmdInit.Stderr = os.Stderr
	if err := cmdInit.Run(); err != nil {
		return fmt.Errorf("git init failed in %s: %w", expandedFolder, err)
	}

	cmdRemote := exec.Command("git", "remote", "add", "origin", gitopsURL)
	cmdRemote.Dir = expandedFolder
	cmdRemote.Stdout = os.Stdout
	cmdRemote.Stderr = os.Stderr
	_ = cmdRemote.Run() // ignore error if origin already exists

	cmdPull := exec.Command("git", "pull", "origin", "HEAD")
	cmdPull.Dir = expandedFolder
	cmdPull.Stdout = os.Stdout
	cmdPull.Stderr = os.Stderr
	if err := cmdPull.Run(); err != nil {
		return fmt.Errorf("git pull origin HEAD failed in %s: %w", expandedFolder, err)
	}

	return nil
}
