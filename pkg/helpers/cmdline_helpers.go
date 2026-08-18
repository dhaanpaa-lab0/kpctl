package helpers

import "github.com/spf13/cobra"

func GetFlagAsString(cmd *cobra.Command, key string) string {
	return cmd.Flags().Lookup(key).Value.String()
}
