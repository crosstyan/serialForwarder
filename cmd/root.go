package cmd

import (
	"github.com/crosstyan/serialForwarder/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

func normalizeName(f *pflag.FlagSet, name string) pflag.NormalizedName {
	return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
}

var rootCmd = cobra.Command{
	Short: "Forward serial data to a TCP socket",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Sugar().Error(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetGlobalNormalizationFunc(normalizeName)
	// https://stackoverflow.com/questions/38105859/make-a-cobra-command-flag-required
	// https://github.com/spf13/cobra/issues/498
	rootCmd.AddCommand(&forwardCmd, &listCmd)
	forwardInit()
}
