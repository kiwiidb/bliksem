package main

import (
	"io"
	"io/ioutil"
	"net"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func forward(localConn net.Conn, serverAddr string, remoteAddr string, config *ssh.ClientConfig) {
	// Setup sshClientConn (type *ssh.ClientConn)
	sshClientConn, err := ssh.Dial("tcp", serverAddr, config)
	if err != nil {
		logrus.Fatalf("ssh.Dial failed: %s", err)
	}

	// Setup sshConn (type net.Conn)
	sshConn, err := sshClientConn.Dial("tcp", remoteAddr)

	// Copy localConn.Reader to sshConn.Writer
	go func() {
		_, err = io.Copy(sshConn, localConn)
		if err != nil {
			logrus.Fatalf("io.Copy failed: %v", err)
		}
	}()

	// Copy sshConn.Reader to localConn.Writer
	go func() {
		_, err = io.Copy(localConn, sshConn)
		if err != nil {
			logrus.Fatalf("io.Copy failed: %v", err)
		}
	}()
}

func startSSHTunnel(conf Config) {
	// Setup SSH config (type *ssh.ClientConfig)
	key, err := ioutil.ReadFile(conf.SSHKeyFile)
	if err != nil {
		logrus.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logrus.Fatalf("unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: conf.SSHUsername,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Setup localListener (type net.Listener)
	localListener, err := net.Listen("tcp", conf.LNDLocalAddr)
	if err != nil {
		logrus.Fatalf("net.Listen failed: %v", err)
	}
	logrus.Info("Startin tunnel")
	for {
		// Setup localConn (type net.Conn)
		localConn, err := localListener.Accept()
		if err != nil {
			logrus.Fatalf("listen.Accept failed: %v", err)
		}
		go forward(localConn, conf.SSHServerAddr, conf.LNDRemoteAddr, sshConfig)
	}
}
