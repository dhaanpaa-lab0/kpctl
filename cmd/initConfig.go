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
	"github.com/spf13/cobra"
	"nexus-sds.com/kpctl/pkg/helpers"
	"nexus-sds.com/kpctl/pkg/platform"
	"nexus-sds.com/kpctl/pkg/sysConfig"
)

// initConfigCmd represents the initConfig command
var initConfigCmd = &cobra.Command{
	Use:   "initConfig",
	Short: "Setup or update K8S On Demand processing CLI",
	Long:  `Setup the initial configuration for K8S On Demand processing CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		profileName := helpers.GetInputAsString("Profile name to configure", sysConfig.GetActiveProfileName())

		sysConfig.SetProfileValuePromptedSelect(profileName, "k8s_context", "K8S Context to use", platform.GetAvailableContexts())
		sysConfig.SetProfileValuePromptedSelect(profileName, "k8s_cluster_ns", "K8S Cluster Namespace", platform.GetNamespaces())
		sysConfig.SetProfileValuePrompted(profileName, "gitops_url", "GitOps URL to use", "^")
		sysConfig.SetProfileValuePrompted(profileName, "gitops_folder", "GitOps Folder to use", "^")

		sysConfig.SetActiveProfileName(profileName)
		sysConfig.UpdateGlobalConfig()
	},
}

func init() {
	rootCmd.AddCommand(initConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
