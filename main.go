package main

import (
	"github.com/derek-elliott/url-shortener/cmd"
	log "github.com/sirupsen/logrus"
)

var (
	gitCommit, date string
)

func main() {
	log.WithFields(log.Fields{
		"Git Commit": gitCommit,
		"Date":       date,
	}).Info("Snip Version")
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
