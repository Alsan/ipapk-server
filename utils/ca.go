package utils

import (
	"errors"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime/debug"
)

func LocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func InitCA() error {
	ip, err := LocalIP()
	if err != nil {
		return err
	}
	_, err = os.Stat(".ca")
	if os.IsNotExist(err) {
		os.Mkdir(".ca", 0755)
		if err := createPrivateCA("IPAPK Generated CA " + ip.String()); err != nil {
			return nil
		}
		if err := createServerCertKey(ip.String()); err != nil {
			return nil
		}
	}
	return nil
}

func createPrivateCA(certificateAuthorityName string) error {
	_, err := callCommand("openssl", "genrsa", "-out", ".ca/myCA.key", "2048")
	if err != nil {
		errors.New("could not create private CA key")
	}

	_, err = callCommand("openssl", "req", "-x509", "-new", "-key", ".ca/myCA.key", "-out", ".ca/myCA.cer", "-days", "730", "-subj", "/CN="+certificateAuthorityName)
	if err != nil {
		errors.New("could not create private CA certificate")
	}
	return nil
}

func createServerCertKey(host string) error {
	_, err := callCommand("openssl", "genrsa", "-out", ".ca/mycert1.key", "2048")
	if err != nil {
		errors.New("could not create private server key")
	}

	_, err = callCommand("openssl", "req", "-new", "-out", ".ca/mycert1.req", "-key", ".ca/mycert1.key", "-subj", "/CN="+host)
	if err != nil {
		errors.New("could not create private server certificate signing request")
	}

	_, err = callCommand("openssl", "x509", "-req", "-in", ".ca/mycert1.req", "-out", ".ca/mycert1.cer", "-CAkey", ".ca/myCA.key", "-CA", ".ca/myCA.cer", "-days", "365", "-CAcreateserial", "-CAserial", ".ca/serial")
	if err != nil {
		errors.New("could not create private server certificate")
	}
	return nil
}

func callCommand(command string, arg ...string) (string, error) {
	out, err := exec.Command(command, arg...).Output()

	if err != nil {
		log.Println("callCommand failed!")
		log.Println("")
		log.Println(string(debug.Stack()))
		return "", err
	}
	return string(out), nil
}
