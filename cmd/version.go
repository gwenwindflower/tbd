package cmd

import (
	"github.com/gwenwindflower/tbd/internal"
)

func init() {
	rootCmd.Version = internal.VERSION
}
