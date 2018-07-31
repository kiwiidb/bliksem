package main

//Config holds basic config parameters for the bliksem rest api and the ssh tunnel
type Config struct {
	Port          string `default:":8081"`
	MacaroonPath  string `default:"secrets/admin.macaroon"`
	TLSPath       string `default:"secrets/tls.cert"`
	SSHUsername   string `json:"username"`
	SSHKeyFile    string `json:"keyfile"`
	SSHServerAddr string `json:"serveraddrstring"`
	LNDLocalAddr  string `json:"localaddrstring"`
	LNDRemoteAddr string `json:"remoteaddrstring"`
	LNDRestAddr   string `json:"lndrestaddr"`
}
