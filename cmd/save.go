package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
	"github.com/ghodss/yaml"
	"github.com/rsmitty/container-shifter/utils"
	"github.com/spf13/cobra"
)

var save = &cobra.Command{
	Use:   "save",
	Short: "Saves docker images to a desired directory",
	Long:  "Saves docker images to a desired directory",
	Run: func(cmd *cobra.Command, args []string) {
		imgPath, err := cmd.Flags().GetString("image-directory")
		utils.ErrorCheck(err)

		csConfig, err := cmd.Flags().GetString("config-file")
		utils.ErrorCheck(err)

		saveImages(imgPath, csConfig)
	},
}

func saveImages(imgPath string, csConfig string) {
	//Create a docker client to use
	cli, err := client.NewEnvClient()
	utils.ErrorCheck(err)

	//Read in config file
	config, err := ioutil.ReadFile(csConfig)
	utils.ErrorCheck(err)
	fmt.Print(string(config))
	var configContents containerConfig
	err = yaml.Unmarshal(config, &configContents)

	//Make img directory if it doesn't exist
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		err = os.MkdirAll(imgPath, 0755)
		utils.ErrorCheck(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(configContents.Containers))
	//Pull down images
	for _, image := range configContents.Containers {
		go func(img string) {
			defer wg.Done()

			//Create reader for image save
			log.Info("Saving " + img)
			out, err := cli.ImageSave(context.Background(), []string{img})
			utils.ErrorCheck(err)
			defer out.Close()

			//Make image name friendly, create tar file and write image
			img = strings.Replace(img, "/", "%", -1)
			imgFile, err := os.Create(imgPath + img + ".tar")
			utils.ErrorCheck(err)
			defer imgFile.Close()
			_, err = io.Copy(imgFile, out)
			utils.ErrorCheck(err)

		}(image)
	}
	wg.Wait()
}

func init() {
	pwd, err := os.Getwd()
	utils.ErrorCheck(err)
	imgPathDefault := pwd + "/img/"
	save.PersistentFlags().String("image-directory", imgPathDefault, "Path to export docker images")
}
