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

// useProfileCmd represents the useProfile command
var useProfileCmd = &cobra.Command{
	Use:   "useProfile",
	Short: "Set the active configuration profile",
	Long:  `Switch the active configuration profile used by other commands (e.g. syncGitops).`,
	Run: func(cmd *cobra.Command, args []string) {
		name := helpers.GetFlagAsString(cmd, "name")
		sysConfig.SetActiveProfileName(name)
		sysConfig.UpdateGlobalConfig()
		fmt.Println("Active profile set to", name)
	},
}

func init() {
	rootCmd.AddCommand(useProfileCmd)
	useProfileCmd.Flags().StringP("name", "n", "", "Profile name to activate")
	errMarkingNameFlagRequired := useProfileCmd.MarkFlagRequired("name")
	if errMarkingNameFlagRequired != nil {

	}
}
