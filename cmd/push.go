package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"sync"

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
		utils.ErrorCheck(err)

		csConfig, err := cmd.Flags().GetString("config-file")
		utils.ErrorCheck(err)

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

		//Create wait group for the push goroutines
		var wg sync.WaitGroup
		wg.Add(len(configContents.Containers))

		for _, image := range configContents.Containers {
			go func(img string) {
				defer wg.Done()

				//Retag image with current registry name
				imageName, err := tagImages(cli, img, registryName)
				utils.ErrorCheck(err)

				//Push image
				out, err := cli.ImagePush(context.Background(), imageName, types.ImagePushOptions{RegistryAuth: authString})
				utils.ErrorCheck(err)
				defer out.Close()
				io.Copy(os.Stdout, out)
			}(image)
		}
		wg.Wait()
	}

}

func init() {
	user, err := user.Current()
	utils.ErrorCheck(err)
	dockerDefault := user.HomeDir + "/.docker/config.json"
	push.PersistentFlags().String("docker-config", dockerDefault, "Path to docker config.")
}
