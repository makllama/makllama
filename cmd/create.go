package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/makllama/makllama/pkg/errors"
)

var (
	debugBW     = false
	overlayType string
	createCmd   = &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			ot := strings.ToLower(overlayType)
			if ot != "vxlan" && ot != "geneve" {
				return fmt.Errorf("invalid overlayType: %s", overlayType)
			}
			err := ensureContainerd()
			if err != nil {
				return err
			}

			localAddress, err := getLocalAddress()
			if err != nil {
				return err
			}

			err = ensureNodes(localAddress)
			if err != nil {
				return err
			}

			err = ensureNetwork(localAddress)
			if err != nil {
				return err
			}

			defer func() { status.End(err == nil) }()

			// actually create nodes
			return errors.UntilErrorConcurrent(nil)
		},
	}
)

func init() {
	createCmd.Flags().BoolVar(&debugBW, "debugBW", false, "enable debug logging for BW")
	createCmd.Flags().StringVar(&overlayType, "overlayType", "vxlan", "overlay type for BW (vxlan or geneve), default is vxlan")
	rootCmd.AddCommand(createCmd)
}

func executableRootPath() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(exePath)
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}

func ensureContainerd() error {
	status.Start("Starting containerd 🚢")
	err := exec.Command(filepath.Join(executableRootPath(), "containerd")).Start()
	if err != nil {
		status.End(false)
		return err
	}
	time.Sleep(2 * time.Second)
	return nil
}

func ensureNodes(localAddress string) error {
	status.Start("Preparing virtual nodes 📦")
	cmd := exec.Command(
		filepath.Join(executableRootPath(), "virtual-kubelet"),
		"--kubeconfig", filepath.Join(homeDir(), ".kube", "config"))
	cmd.Env = append(os.Environ(), fmt.Sprintf("VKUBELET_POD_IP=%s", localAddress))
	err := cmd.Start()
	if err != nil {
		status.End(false)
		return err
	}
	time.Sleep(2 * time.Second)
	return nil
}

func getClusterCIDR() (string, error) {
	cmd := exec.Command("bash", "-c", "kubectl get cm -n kube-system kubeadm-config -oyaml | grep podSubnet | awk '{print $2}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getPodCIDR() (string, error) {
	cmd := exec.Command("scutil", "--get", "LocalHostName")
	localHostName, err := cmd.Output()
	if err != nil {
		return "", err
	}
	nodeName := strings.TrimSpace(string(localHostName))

	cmd = exec.Command("bash", "-c", fmt.Sprintf("kubectl get no %s -oyaml | grep podCIDR | awk '{print $2}'", nodeName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getLocalAddress() (string, error) {
	iface, err := net.InterfaceByName("en0")
	if err != nil {
		return "", err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		ip := ipNet.IP
		if ip.To4() != nil {
			ipv4 := ip.String()
			return ipv4, err
		}
	}

	return "", errors.New("no IPv4 address found")
}

func ensureNetwork(localAddress string) error {
	status.Start("Creating network 🌐")
	var (
		clusterCIDR string
		podCIDR     string
		port        = "4789" // VXLAN default port
		err         error
		cmd         *exec.Cmd
	)

	clusterCIDR, err = getClusterCIDR()
	if err != nil {
		goto Error
	}

	podCIDR, err = getPodCIDR()
	if err != nil {
		goto Error
	}

	if strings.ToLower(overlayType) == "geneve" {
		port = "6081" // Geneve default port
	}

	cmd = exec.Command(
		filepath.Join(executableRootPath(), "bronzewillow"),
		"--local-address", localAddress,
		"--cluster-cidr", clusterCIDR,
		"--pod-cidr", podCIDR,
		"--overlay-type", overlayType,
		"--port", port,
	)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("DATABASE_URL=%s", filepath.Join(homeDir(), "bw.db")))
	if debugBW {
		cmd.Env = append(cmd.Env, "RUST_LOG=trace")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		err = cmd.Run()
	} else {
		err = cmd.Start()
	}
	if err != nil {
		goto Error
	}
	time.Sleep(2 * time.Second)
	return nil

Error:
	status.End(false)
	return err
}
