/*
This file is here to simply act as an entrypoint for the kube-client binary.
It doesn't do anything except execute the root cobra command.
*/
package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/rsmitty/container-shifter/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
