package cmd

import (
	"github.com/makllama/makllama/pkg/log"

	"github.com/makllama/makllama/pkg/internal/cli"
)

func StatusForLogger(l log.Logger) *cli.Status {
	return cli.StatusForLogger(l)
}
