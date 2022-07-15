package src

import (
	"log"
	"net"
	"os"
)

// fingerprinting function for database filtering
func FingerPrint() {
	ifaces, err := net.Interfaces()
	if err != nil {
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			hostname, err := os.Hostname()
			log.Println(ip, hostname)
			if err != nil {
			}
		}
	}
}
