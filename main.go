package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func main() {
	service := &BliksemService{}
	service.Initialize()
	router := mux.NewRouter()
	router.HandleFunc("/addinvoice",
		func(w http.ResponseWriter, r *http.Request) { handleAddInvoice(w, r, service) })
	logrus.Info("Starting Server")
	logrus.Fatal(http.ListenAndServe("localhost:8081", router))
}

func handleAddInvoice(w http.ResponseWriter, r *http.Request, service *BliksemService) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	var inv Invoice
	err = json.Unmarshal(bytes, &inv)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	newInv := service.getNewInvoice(inv.Amount)
	toSendBytes, err := json.Marshal(newInv)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	w.Write(toSendBytes)
}
