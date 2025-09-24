package main

import (
	"io"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	//TODO: next add write to file
	if !debug {
		log.SetOutput(io.Discard)
	}
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(Storage.ServerLogLevel())
}
