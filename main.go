//go:build linux

package main

import "log"

func main() {
	log.Fatal(work())
}

func work() error {
	err := icmpListen()
	if err != nil {
		return err
	}
	if cfgRole == Client {
		err = udpListen()
		if err != nil {
			return err
		}
	}
	go icmp2udp()
	select {}
}
