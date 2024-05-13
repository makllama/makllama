package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	ps "github.com/mitchellh/go-ps"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func getProcessesByNames(names []string) ([]ps.Process, error) {
	allProcesses, err := ps.Processes()
	if err != nil {
		return nil, err
	}

	var filteredProcesses []ps.Process
	for _, process := range allProcesses {
		for _, name := range names {
			if strings.Contains(process.Executable(), name) {
				filteredProcesses = append(filteredProcesses, process)
			}
		}
	}

	return filteredProcesses, nil
}

func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Signal(syscall.SIGKILL)
	if err != nil {
		return err
	}

	return nil
}

var deleteCmd = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		status.Start("Deleting demo 🧹")

		// https://github.com/mitchellh/go-ps/issues/52
		// ollama_llama => ollama_llama_server
		processes, err := getProcessesByNames([]string{"containerd-shim-runc-v2", "containerd",
			"virtual-kubelet", "runm", "ollama_llama", "bronzewillow"})
		if err != nil {
			panic(err)
		}

		for _, process := range processes {
			killProcess(process.Pid())
		}

		os.Remove(filepath.Join(homeDir(), "bw.db"))

		time.Sleep(2 * time.Second)

		defer func() { status.End(err == nil) }()
	},
}
