package main

import (
	"flag"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strings"
	"time"
)

type Role bool
const (
	Client Role = false
	Server Role = true
)

var (
	cfgRole        Role
	cfgUdpAddrPort netip.AddrPort
	cfgIcmpAddr    netip.Addr
	cfgIcmpTypeC2S uint
	cfgIcmpTypeS2C uint
	cfgIcmpCode    uint
	cfgTimeout     time.Duration
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), `Usage:
  client: %s [options] <UDP listen address:port> <ICMP server address>
  server: %s [options] <ICMP listen address> <UDP upstream address:port>
`, os.Args[0], os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.UintVar(&cfgIcmpTypeC2S, "icmp-type-c2s", uint(ICMPTypeAddressMaskRequest), "ICMP type (client to server)")
	flag.UintVar(&cfgIcmpTypeS2C, "icmp-type-s2c", uint(ICMPTypeAddressMaskReply), "ICMP type (server to client)")
	flag.UintVar(&cfgIcmpCode, "icmp-code", 255, "ICMP code")
	flag.DurationVar(&cfgTimeout, "timeout", time.Minute, "Session timeout")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 2 {
		usage()
		os.Exit(1)
	}

	if cfgIcmpTypeC2S > 255 || cfgIcmpTypeS2C > 255 {
		log.Fatal("invalid icmp-type")
	}
	if cfgIcmpCode > 255 {
		log.Fatal("invalid icmp-code")
	}

	// parse remaining arguments
	addr1, addr2 := flag.Arg(0), flag.Arg(1)
	isClient := strings.ContainsRune(addr1, ':')
	isServer := strings.ContainsRune(addr2, ':')
	if isClient == isServer {
		log.Fatal("invalid address")
	}
	if isClient {
		cfgRole = Client
		cfgUdpAddrPort = netip.MustParseAddrPort(addr1)
		cfgIcmpAddr = netip.MustParseAddr(addr2)
	} else {
		cfgRole = Server
		cfgIcmpAddr = netip.MustParseAddr(addr1)
		cfgUdpAddrPort = netip.MustParseAddrPort(addr2)
	}
}

// Returns ICMP type used by packets between client and server.
func cfgIcmpType(dir Direction) ICMPType {
	if isClientToServer(dir) {
		return ICMPType(cfgIcmpTypeC2S)
	} else {
		return ICMPType(cfgIcmpTypeS2C)
	}
}
