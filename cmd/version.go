package cmd

import (
	"fmt"

	"github.com/gwenwindflower/tbd/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of tbd",
	Long:  `Wanna know what version of tbd you're running? Well I've got some great news for you.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tbd v%s\n", internal.VERSION)
	},
}
