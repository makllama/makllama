package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/makllama/makllama/pkg/cmd"
)

var (
	rootCmd = &cobra.Command{}
	logger  = cmd.NewLogger()
	status  = cmd.StatusForLogger(logger)
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
