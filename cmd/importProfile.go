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
	"nexus-sds.com/kpctl/pkg/helpers"
	"nexus-sds.com/kpctl/pkg/sysConfig"
)

// importProfileCmd represents the importProfile command
var importProfileCmd = &cobra.Command{
	Use:     "importProfile",
	Aliases: []string{"import-profile"},
	Short:   "Import a configuration profile from a YAML file",
	Long:    `Import a configuration profile from a YAML file, creating the profile if it does not already exist or updating it if it does.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := helpers.GetFlagAsString(cmd, "file")
		name := helpers.GetFlagAsString(cmd, "name")

		imported, err := sysConfig.ImportProfile(path, name)
		if err != nil {
			return err
		}
		sysConfig.UpdateGlobalConfig()
		fmt.Printf("Profile %q imported from %s\n", imported, path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importProfileCmd)
	importProfileCmd.Flags().StringP("file", "f", "", "YAML file to import the profile from")
	errMarkingFileFlagRequired := importProfileCmd.MarkFlagRequired("file")
	if errMarkingFileFlagRequired != nil {

	}
	importProfileCmd.Flags().StringP("name", "n", "", "Profile name to import as (default: name field in the file)")
}
