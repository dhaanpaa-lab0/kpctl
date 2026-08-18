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
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/dynamic"
	"nexus-sds.com/kpctl/pkg/helpers"
	"nexus-sds.com/kpctl/pkg/platform"
	"nexus-sds.com/kpctl/pkg/sysConfig"
)

// externalSecretCmd represents the externalSecret command
var externalSecretCmd = &cobra.Command{
	Use:     "externalSecret",
	Aliases: []string{"external-secret"},
	Short:   "Create an ExternalSecret manifest, writing it to a file and/or applying it to the cluster",
	Long: `Build an ExternalSecret manifest from the given flags, prompting for any that
are not specified. Use --file to write the manifest as YAML, --apply to apply
it directly to the cluster (using the active profile's k8s_context and
k8s_cluster_ns), or both.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file := helpers.GetFlagAsString(cmd, "file")
		apply, err := cmd.Flags().GetBool("apply")
		if err != nil {
			return err
		}
		if file == "" && !apply {
			return errors.New("specify --file, --apply, or both")
		}

		profile := sysConfig.GetActiveProfile()

		name := helpers.GetFlagAsString(cmd, "name")
		if name == "" {
			name = helpers.GetInputAsString("ExternalSecret name", "")
		}

		namespace := helpers.GetFlagAsString(cmd, "namespace")
		if namespace == "" {
			namespace = helpers.GetInputAsString("K8S namespace", profile.K8sClusterNs)
		}

		storeKey := helpers.GetFlagAsString(cmd, "secret-store-key")
		if storeKey == "" {
			storeKey = helpers.GetInputAsString("Secret store key", "")
		}

		storeName := helpers.GetFlagAsString(cmd, "secret-store-name")
		if storeName == "" {
			storeName = helpers.GetInputAsString("Secret store name", "")
		}

		params := platform.ExternalSecretParams{
			ExternalSecretName: name,
			K8sNamespace:       namespace,
			SecretStoreKey:     storeKey,
			SecretStoreName:    storeName,
		}

		if file != "" {
			yamlText, err := platform.GenerateExternalSecretYAML(params)
			if err != nil {
				return err
			}
			if err := os.WriteFile(file, []byte(yamlText), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", file, err)
			}
			fmt.Printf("ExternalSecret manifest written to %s\n", file)
		}

		if apply {
			restConfig := platform.GetK8SPlatformRestConfigForContext(profile.K8sContext)
			dynClient, err := dynamic.NewForConfig(restConfig)
			if err != nil {
				return fmt.Errorf("failed to build k8s client: %w", err)
			}
			if _, err := platform.ApplyExternalSecret(cmd.Context(), dynClient, params, "kpctl"); err != nil {
				return err
			}
			fmt.Printf("ExternalSecret %q applied to namespace %q\n", name, namespace)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(externalSecretCmd)
	externalSecretCmd.Flags().StringP("name", "n", "", "ExternalSecret name (prompted if not set)")
	externalSecretCmd.Flags().String("namespace", "", "K8S namespace (defaults to the active profile's k8s_cluster_ns, prompted if not set)")
	externalSecretCmd.Flags().String("secret-store-key", "", "SecretStore key to extract (prompted if not set)")
	externalSecretCmd.Flags().String("secret-store-name", "", "SecretStore name to reference (prompted if not set)")
	externalSecretCmd.Flags().StringP("file", "f", "", "YAML file to write the ExternalSecret manifest to")
	externalSecretCmd.Flags().Bool("apply", false, "Apply the ExternalSecret to the cluster")
}
