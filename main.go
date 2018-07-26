package main

import (
	"encoding/json"
	"fmt"
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
	err = json.Unmarshal([]byte("{"+invBody.Body+"}"), &inv)
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
