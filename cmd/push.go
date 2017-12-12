package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ghodss/yaml"
	"github.com/rsmitty/container-shifter/utils"
	"github.com/spf13/cobra"
)

var push = &cobra.Command{
	Use:   "push SUBCOMMAND",
	Short: "Allows for pushing docker images",
	Long:  "Allows for pushing docker images",
	Run: func(cmd *cobra.Command, args []string) {
		dockerConfig, err := cmd.Flags().GetString("docker-config")
		if err != nil {
			log.Fatal(err)
		}
		csConfig, err := cmd.Flags().GetString("config-file")
		if err != nil {
			log.Fatal(err)
		}
		pushImages(dockerConfig, csConfig)
	},
}

func tagImages(clientPtr *client.Client, imageName string, privateRegistry string) (string, error) {
	imageSlice := strings.SplitN(imageName, "/", 2)
	imageRoot := strings.Join(imageSlice[1:], "")
	newImageName := privateRegistry + "/" + imageRoot
	err := clientPtr.ImageTag(context.Background(), imageName, newImageName)
	return newImageName, err
}

func pushImages(dockerConfig string, csConfig string) {
	//Create a docker client to use for push/pulls
	cli, err := client.NewEnvClient()
	utils.ErrorCheck(err)

	//Read in config file
	config, err := ioutil.ReadFile(csConfig)
	utils.ErrorCheck(err)
	fmt.Print(string(config))
	var configContents containerConfig
	err = yaml.Unmarshal(config, &configContents)

	//Push out images
	for _, registryName := range configContents.Registries {
		//Check if credentials already exist in docker configs for this registry
		authString := utils.RegistryAuth(registryName, dockerConfig)

		for _, image := range configContents.Containers {
			//Retag image with current registry name
			imageName, err := tagImages(cli, image, registryName)
			utils.ErrorCheck(err)

			//Push image
			out, err := cli.ImagePush(context.Background(), imageName, types.ImagePushOptions{RegistryAuth: authString})
			utils.ErrorCheck(err)
			defer out.Close()
			io.Copy(os.Stdout, out)
		}

	}

}

func init() {
	user, err := user.Current()
	utils.ErrorCheck(err)
	dockerDefault := user.HomeDir + "/.docker/config.json"
	push.PersistentFlags().String("docker-config", dockerDefault, "Path to docker config.")
}
