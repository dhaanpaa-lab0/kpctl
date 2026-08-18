/*
Copyright © 2025 Daniel Haanpaa <djh@nexus-sds.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"nexus-sds.com/kpctl/pkg/platform"
	"nexus-sds.com/kpctl/pkg/sysConfig"
)

// syncGitopsCmd represents the syncGitops command
var syncGitopsCmd = &cobra.Command{
	Use:     "syncGitops",
	Aliases: []string{"sync", "sync-gitops"},
	Short:   "Sync GitOps repository into local folder",
	Long: `Sync from the git repository defined with the viper config key 'gitops_url'
into the folder defined with the viper config key 'gitops_folder', creating the folder if it does not exist.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := sysConfig.GetActiveProfile()
		if err := platform.SyncGitopsRepository(profile.GitopsURL, profile.GitopsFolder); err != nil {
			return fmt.Errorf("gitops sync failed: %w", err)
		}
		fmt.Println("GitOps repository synced successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncGitopsCmd)
}
