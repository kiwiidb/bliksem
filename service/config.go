package main

//Config holds basic config parameters for the bliksem rest api
type Config struct {
	Port         string `default:":8081"`
	MacaroonPath string `default:"../secrets/invoice.macaroon"`
	TLSPath      string `default:"../secrets/tls.cert"`
	LNDRestAddr  string `json:"lndrestaddr"`
}
