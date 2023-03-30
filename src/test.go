package main

import (
	"fmt"
	"log"
	"net"

	gm "github.com/c-seeger/mac-gen-go"
)

func main() {

	ips := []net.IP{
		net.ParseIP("10.0.0.0"),
		net.ParseIP("10.255.255.255"),
		net.ParseIP("192.168.12.127"),
	}

	for _, ip := range ips {
		// generate a random local administered unicast mac prefix
		prefix := gm.GenerateRandomLocalMacPrefix(true)

		// calculates the NIC Sufix by ip address
		sufix, err := gm.CalculateNICSufix(ip)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s - %s:%s\n", ip.String(), prefix, sufix)
	}
}
