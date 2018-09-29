package main

import (
	"net/http"

	"github.com/koding/multiconfig"

	"github.com/sirupsen/logrus"
)

const (
	configFile = "../secrets/conf.json"
)

func main() {
	conf := Config{}
	m := multiconfig.NewWithPath(configFile)
	err := m.Load(&conf)
	if err != nil {
		logrus.WithError(err).Fatal("Could not load config!")
	}
	service := &BliksemService{}
	service.Initialize(conf)

	go service.streamInvoicesToChannel()

	router := setupRouter(service, conf)

	logrus.Info("Starting Server")
	logrus.Fatal(http.ListenAndServe(conf.Port, router))
}
