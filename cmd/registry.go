package cmd

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/rsmitty/container-shifter/utils"
	"github.com/spf13/cobra"
)

var registry = &cobra.Command{
	Use:   "registry",
	Short: "Creates a temporary registry",
	Long:  "Creates a temporary registry to host containers",
	Run: func(cmd *cobra.Command, args []string) {
		regPath, err := cmd.Flags().GetString("registry-dir")
		utils.ErrorCheck(err)

		runReg(regPath)
	},
}

func runReg(regPath string) {
	//Create a docker client to use
	cli, err := client.NewEnvClient()
	utils.ErrorCheck(err)

	containerConfig := &container.Config{
		Image: "library/registry:2",
		ExposedPorts: nat.PortSet{
			"5000/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		Binds: []string{
			regPath + ":/var/lib/registry",
		},
		PortBindings: nat.PortMap{
			"5000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "5000",
				},
			},
		},
	}
	resp, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig,
		nil, "docker-registry")
	utils.ErrorCheck(err)

	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	utils.ErrorCheck(err)

	out, err := cli.ContainerLogs(context.Background(), resp.ID,
		types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	utils.ErrorCheck(err)

	io.Copy(os.Stdout, out)

}

func init() {
	registry.PersistentFlags().String("registry-dir", "/opt/docker-registry", "Path for registry storage")
}
