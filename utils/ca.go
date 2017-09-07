package utils

import (
	"net"
	"os/exec"
	"os"
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
		cmd := exec.Command("/bin/sh", "certificate", ip.String())
		_, err = cmd.Output()
		if err != nil {
			return err
		}
	}
	return nil
}
