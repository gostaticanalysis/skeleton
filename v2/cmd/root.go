package cmd

import (
	"github.com/spf13/cobra"

	_ "embed"
)

//go:embed Version.txt
var version string

var (
	flagChecker string
)

var rootCmd = &cobra.Command{
	Use:     "skeleton",
	Short:   "",
	Long:    ``,
	RunE:    runE,
	Version: version,
}

func init() {
	rootCmd.Flags().StringVarP(&flagChecker, "checker", "c", "unit", "[unit,single,multi]")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func runE(cmd *cobra.Command, args []string) error {
	var kind Kind
	if len(args) > 0 {
		kind = ParseKind(args[0])
	}

	s := &Skeleton{
		Kind:     kind,
		Checker:  flagChecker,
		Template: tmpl,
	}

	if err := s.Execute(); err != nil {
		return err
	}

	return nil
}
