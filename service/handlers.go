package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func handleAddInvoice(w http.ResponseWriter, r *http.Request, service *BliksemService) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Error("Error reading invoice request")
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}
	var invBody ReqInvoice
	err = json.Unmarshal(bytes, &invBody)
	if err != nil {
		logrus.WithError(err).Error("Error decoding invoice request")
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}
	var inv Invoice
	invBytes := "{" + invBody.Body + "}"
	err = json.Unmarshal([]byte(invBytes), &inv)
	if err != nil {
		logrus.WithError(err).Error("Error decoding request")
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}
	newInv := service.getNewInvoice(inv.Amount)
	toSendBytes, err := json.Marshal(newInv)
	if err != nil {
		logrus.WithError(err).Error("Error encoding request")
		http.Error(w, "Something horribly wrong", http.StatusInternalServerError)
		return
	}
	w.Write(toSendBytes)
}

func handleSettledInvoice(w http.ResponseWriter, r *http.Request, service *BliksemService, upgrader websocket.Upgrader) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Fatal("something wrong upgrader")
	}
	defer c.Close()
	for {
		_, payreq, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		toCompareInv := <-service.invoiceChan
		if toCompareInv.PayReq == string(payreq) {
			c.WriteMessage(websocket.TextMessage, []byte("true"))
		}
	}
}
