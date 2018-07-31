package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/koding/multiconfig"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	configFile = "secrets/conf.json"
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

	go startSSHTunnel(conf)

	router := mux.NewRouter()
	router.HandleFunc("/addinvoice",
		func(w http.ResponseWriter, r *http.Request) { handleAddInvoice(w, r, service) })
	logrus.Info("Starting Server")
	logrus.Fatal(http.ListenAndServe(conf.Port, router))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func handleAddInvoice(w http.ResponseWriter, r *http.Request, service *BliksemService) {
	enableCors(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	fmt.Println(string(bytes))
	var invBody ReqInvoice
	err = json.Unmarshal(bytes, &invBody)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	fmt.Println(invBody.Body)
	var inv Invoice
	invBytes := "{" + invBody.Body + "}"
	fmt.Println(invBytes)
	err = json.Unmarshal([]byte(invBytes), &inv)
	if err != nil {
		logrus.WithError(err).Fatal("here")
	}
	newInv := service.getNewInvoice(inv.Amount)
	toSendBytes, err := json.Marshal(newInv)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	w.Write(toSendBytes)
}
