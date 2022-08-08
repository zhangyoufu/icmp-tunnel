package main

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"net/netip"

	"golang.org/x/net/ipv4"
)

const (
	MTU = 1500
	maxPayloadSize = MTU - 20 - 8
)

var (
	udpConn *net.UDPConn
	ipConn  *net.IPConn
)

func icmpListen() error {
	var err error
	ipConn, err = net.ListenIP("ip4:icmp", nil)
	if err != nil {
		return err
	}

	// set ICMP_FILTER, see raw(7)
	// do not use ipv4.NewRawConn, which sets IP_HDRINCL
	icmpFilter := ipv4.ICMPFilter{}
	icmpFilter.SetAll(true)
	icmpFilter.Accept(cfgIcmpType(Recv))
	ipv4.NewPacketConn(ipConn).SetICMPFilter(&icmpFilter)

	return nil
}

func udpListen() error {
	var err error
	udpConn, err = net.ListenUDP("udp", net.UDPAddrFromAddrPort(cfgUdpAddrPort))
	if err != nil {
		return err
	}
	go udp2icmp(udpConn)
	return nil
}

func udpConnLocalAddrPort(udpConn *net.UDPConn) netip.AddrPort {
	udpAddr := udpConn.LocalAddr().(*net.UDPAddr)
	addr, _ := netip.AddrFromSlice(udpAddr.IP)
	return netip.AddrPortFrom(addr.WithZone(udpAddr.Zone), uint16(udpAddr.Port))
}

func udp2icmp(udpConn *net.UDPConn) {
	var udpAddr netip.AddrPort
	if cfgRole == Server {
		// we have multiple UDPConn to upstream with different local address
		udpAddr = udpConnLocalAddrPort(udpConn)
	}
	for {
		var (
			buf [8+maxPayloadSize]byte
			n int
			err error
		)
		if cfgRole == Client {
			n, udpAddr, err = udpConn.ReadFromUDPAddrPort(buf[8:])
			if err != nil {
				log.Fatal(err)
			}
		}
		if cfgRole == Server {
			n, err = udpConn.Read(buf[8:])
		}
		session := getSession(UdpTable, udpAddr, cfgRole == Client)
		if err != nil { // implies cfgRole == Server
			if !errors.Is(err, net.ErrClosed) {
				log.Printf("udp:%s %s", udpAddr, err)
				if session != nil {
					session.Purge()
				}
			}
			return
		}
		if session == nil {
			continue
		}
		session.Refresh()
		buf[0] = byte(cfgIcmpType(Send))
		buf[1] = byte(cfgIcmpCode)
		buf[2], buf[3] = 0, 0 // Checksum
		binary.BigEndian.PutUint16(buf[4:6], session.IcmpId())
		binary.BigEndian.PutUint16(buf[6:8], session.icmpSeq)
		fillChecksum(buf[:n+8], 2)
		_, _ = ipConn.WriteToIP(buf[:n+8], &session.icmpIPAddr)
		if cfgRole == Client {
			log.Printf("udp:%s -> icmp:%s(seq=%d)", udpAddr, session.icmpAddr, session.icmpSeq)
		} else {
			log.Printf("icmp:%s(seq=%d) <-(udp:%s)-- upstream", session.icmpAddr, session.icmpSeq, udpAddr)
		}
		session.icmpSeq++
	}
}

func icmp2udp() {
	for {
		var buf [8+maxPayloadSize]byte
		n, ipAddr, err := ipConn.ReadFromIP(buf[:]) // without IPv4 header and options
		if err != nil {
			log.Fatal(err)
		}
		if n <= 8 {
			continue
		}
		icmpType := ICMPType(buf[0])
		if icmpType != cfgIcmpType(Recv) {
			continue
		}
		icmpCode := uint(buf[1])
		if icmpCode != cfgIcmpCode {
			continue
		}
		icmpId := binary.BigEndian.Uint16(buf[4:6])
		addr, _ := netip.AddrFromSlice(ipAddr.IP)
		icmpAddr := netip.AddrPortFrom(addr, icmpId)
		session := getSession(IcmpTable, icmpAddr, cfgRole == Server)
		if session == nil {
			continue
		}
		if icmpAddr != session.icmpAddr {
			continue
		}
		session.Refresh()
		icmpSeq := binary.BigEndian.Uint16(buf[6:8])
		if cfgRole == Client {
			_, _ = session.udpConn.WriteToUDPAddrPort(buf[8:n], session.udpAddr)
			log.Printf("udp:%s <- icmp:%s(seq=%d)", session.udpAddr, icmpAddr, icmpSeq)
		} else {
			_, _ = session.udpConn.Write(buf[8:n])
			log.Printf("icmp:%s(seq=%d) --(udp:%s)-> upstream", icmpAddr, icmpSeq, session.udpAddr)
		}
	}
}
