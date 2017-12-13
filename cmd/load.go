package cmd

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
	"github.com/rsmitty/container-shifter/utils"
	"github.com/spf13/cobra"
)

var imageLoad = &cobra.Command{
	Use:   "load",
	Short: "Loads docker images to local system",
	Long:  "Loads docker images to local system",
	Run: func(cmd *cobra.Command, args []string) {
		imgPath, err := cmd.Flags().GetString("image-directory")
		utils.ErrorCheck(err)

		csConfig, err := cmd.Flags().GetString("config-file")
		utils.ErrorCheck(err)
		loadImages(imgPath, csConfig)
	},
}

func loadImages(imgPath string, csConfig string) {
	//Create a docker client to use
	cli, err := client.NewEnvClient()
	utils.ErrorCheck(err)

	//Find all tar files in the imgPath
	var tarFiles []string
	imgDirectory, err := os.Open(imgPath)
	utils.ErrorCheck(err)
	defer imgDirectory.Close()

	files, err := imgDirectory.Readdir(-1)
	utils.ErrorCheck(err)

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".tar" {
			tarFiles = append(tarFiles, imgPath+file.Name())
		}
	}

	//Load images
	var wg sync.WaitGroup
	wg.Add(len(tarFiles))
	for _, tarFile := range tarFiles {
		go func(tar string) {
			defer wg.Done()
			file, err := os.Open(tar)
			utils.ErrorCheck(err)
			out, err := cli.ImageLoad(context.Background(), io.Reader(file), false)
			utils.ErrorCheck(err)
			log.Info(out)
		}(tarFile)
	}
	wg.Wait()

}

func init() {
	pwd, err := os.Getwd()
	utils.ErrorCheck(err)
	imgPathDefault := pwd + "/img/"
	imageLoad.PersistentFlags().String("image-directory", imgPathDefault, "Path to load docker images")
}
