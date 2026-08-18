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
	"maps"
	"slices"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"nexus-sds.com/kpctl/pkg/helpers"
	"nexus-sds.com/kpctl/pkg/sysConfig"
)

// listGlobalsCmd represents the listGlobals command
var listGlobalsCmd = &cobra.Command{
	Use:     "listGlobals",
	Aliases: []string{"listGlobal", "list-globals"},
	Short:   "List global configuration variables in a table format",
	Long:    `List all global configuration variables in a table format using pterm.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		globals := sysConfig.GetGlobalsConfigMap()
		tableData := pterm.TableData{
			{"Key", "Value"},
		}
		for _, k := range slices.Sorted(maps.Keys(globals)) {
			tableData = append(tableData, []string{k, globals[k]})
		}
		return helpers.RenderTable(tableData)
	},
}

func init() {
	rootCmd.AddCommand(listGlobalsCmd)
}
