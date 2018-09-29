package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	serviceInvoiceRoute    = "/v1/invoices"
	serviceInvoiceSubRoute = "/v1/invoices/subscribe"
)

func newTLSClient(certFile string) *http.Client {
	caCert, _ := ioutil.ReadFile(certFile)
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	tlsConfig := tls.Config{
		RootCAs: pool,
	}

	transport := http.Transport{
		TLSClientConfig: &tlsConfig,
	}

	client := http.Client{
		Transport: &transport,
	}
	return &client
}

func getMacaroon(macaroonFile string) string {
	macaroonbytes, err := ioutil.ReadFile(macaroonFile)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	macaroon := hex.EncodeToString(macaroonbytes)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	return macaroon
}

//Initialize url, macaroon
func (service *BliksemService) Initialize(conf Config) {
	service.client = newTLSClient(conf.TLSPath)
	service.url = conf.LNDRestAddr
	service.macaroon = getMacaroon(conf.MacaroonPath)
	service.invoiceChan = make(chan Invoice, 1000)
}

func (service BliksemService) getNewInvoice(amount int64) Invoice {
	inv := Invoice{Amount: amount}
	logrus.WithField("invoice", inv).Info("Requesting invoice to LND backend")
	invoiceBytes, err := json.Marshal(inv)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	req, err := http.NewRequest("POST", service.url+serviceInvoiceRoute, bytes.NewBuffer(invoiceBytes))
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	req.Header.Set("Grpc-Metadata-macaroon", service.macaroon)
	res, err := service.client.Do(req)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	logrus.WithField("response", string(body)).Info("Response from LND backend")
	var receivedInv Invoice
	err = json.Unmarshal(body, &receivedInv)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	return receivedInv

}

func (service BliksemService) streamInvoicesToChannel() {
	logrus.Info("Starting invoice stream")
	req, err := http.NewRequest("GET", service.url+serviceInvoiceSubRoute, nil)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	req.Header.Set("Grpc-Metadata-macaroon", service.macaroon)
	res, err := service.client.Do(req)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	inv := LNDStreamInvoice{}
	decoder := json.NewDecoder(res.Body)
	for {
		err := decoder.Decode(&inv)
		if err != nil {
			logrus.WithError(err).Fatal()
		}
		if inv.Result.Settled {
			logrus.Info("Payment received")
			logrus.WithField("invoice", inv).Info("This invoice was recently paid")
			amt, err := strconv.Atoi(inv.Result.Amount)
			if err != nil {
				logrus.WithError(err).Fatal()
			}
			memo := inv.Result.Memo
			payR := inv.Result.PayReq

			invoice := Invoice{Amount: int64(amt), Memo: memo, PayReq: payR}
			service.invoiceChan <- invoice

		}
		inv = LNDStreamInvoice{}
	}
}
