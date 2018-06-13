package main

import (
	"github.com/derek-elliott/url-shortner/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
