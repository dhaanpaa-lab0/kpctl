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
	"nexus-sds.com/kpctl/pkg/sysConfig"
)

// deleteProfileCmd represents the deleteProfile command
var deleteProfileCmd = &cobra.Command{
	Use:   "deleteProfile",
	Short: "Delete a configuration profile",
	Long:  `Delete the named configuration profile. If it is the active profile, the active profile reverts to "default".`,
	Run: func(cmd *cobra.Command, args []string) {
		sysConfig.DeleteProfile(helpers.GetFlagAsString(cmd, "name"))
		sysConfig.UpdateGlobalConfig()
	},
}

func init() {
	rootCmd.AddCommand(deleteProfileCmd)
	deleteProfileCmd.Flags().StringP("name", "n", "", "Profile name to delete")
	errMarkingNameFlagRequired := deleteProfileCmd.MarkFlagRequired("name")
	if errMarkingNameFlagRequired != nil {

	}
}
