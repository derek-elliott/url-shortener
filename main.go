package main

import (
	"github.com/derek-elliott/url-shortener/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
