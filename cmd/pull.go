package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ghodss/yaml"
	"github.com/rsmitty/container-shifter/utils"
	"github.com/spf13/cobra"
)

var pull = &cobra.Command{
	Use:   "pull SUBCOMMAND",
	Short: "Allows for pulling docker images",
	Long:  "Allows for pulling docker images",
	Run: func(cmd *cobra.Command, args []string) {
		csConfig, err := cmd.Flags().GetString("config-file")
		if err != nil {
			log.Fatal(err)
		}
		pullImages(csConfig)
	},
}

func pullImages(csConfig string) {
	//Create a docker client to use for push/pulls
	cli, err := client.NewEnvClient()
	utils.ErrorCheck(err)

	//Read in config file
	config, err := ioutil.ReadFile(csConfig)
	utils.ErrorCheck(err)
	fmt.Print(string(config))
	var configContents containerConfig
	err = yaml.Unmarshal(config, &configContents)

	var wg sync.WaitGroup
	wg.Add(len(configContents.Containers))

	//Pull down images
	for _, image := range configContents.Containers {
		go func(img string) {
			defer wg.Done()

			log.Info("Pulling " + img)
			out, err := cli.ImagePull(context.Background(), img, types.ImagePullOptions{})
			utils.ErrorCheck(err)

			defer out.Close()
			io.Copy(os.Stdout, out)
		}(image)
	}
	wg.Wait()
}
