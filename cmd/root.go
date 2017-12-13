/*
This file implements the root command for cobra and adds the pods subcommand.
It also handles the setting of the log level with the --log-level flag
*/

package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/rsmitty/container-shifter/utils"
	"github.com/spf13/cobra"
)

type containerConfig struct {
	Containers []string
	Registries []string
}

//RootCmd adds 'kube-client' root command and handles the log-level flag for root and any subcommands
var RootCmd = &cobra.Command{
	Use:   "container-shifter",
	Short: "A golang binary to push/pull docker images",
	Long:  "This is a golang tool to support the mass push and pull of docker images from public to private registries.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		if err != nil {
			log.Fatal(err)
		}
		switch logLevel {
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		default:
			log.SetLevel(log.InfoLevel)
		}
	},
}

//Ensure log-level flag can be used on any subcommand and add subcommands
func init() {
	RootCmd.PersistentFlags().String("log-level", "info", "A desired logging level. Supported vals: debug, info, error")

	pwd, err := os.Getwd()
	utils.ErrorCheck(err)
	configDefault := pwd + "/config.yml"
	RootCmd.PersistentFlags().String("config-file", configDefault, "Path to container-shifter config file")
	RootCmd.AddCommand(pull)
	RootCmd.AddCommand(push)
	RootCmd.AddCommand(save)
	RootCmd.AddCommand(imageLoad)
}
