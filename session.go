package main

import (
	"log"
	"net"
	"net/netip"
	"sync"
	"time"
)

var (
	mutex sync.RWMutex
	udpTable  = make(map[netip.AddrPort]*Session)
	icmpTable = make(map[netip.AddrPort]*Session)
)

type TableId bool
const UdpTable = TableId(false)
const IcmpTable = TableId(true)

type Session struct {
	activity   time.Time
	udpConn    *net.UDPConn
	udpAddr    netip.AddrPort // remote address on client side, local address on server side
	icmpAddr   netip.AddrPort
	icmpIPAddr net.IPAddr
	icmpSeq    uint16
}

func (s *Session) Purge() {
	mutex.Lock()
	delete(udpTable, s.udpAddr)
	delete(icmpTable, s.icmpAddr)
	mutex.Unlock()
	if s.udpConn != udpConn {
		_ = s.udpConn.Close()
	}
}

func (s *Session) Refresh() {
	s.activity = time.Now()
}

func (s *Session) IcmpId() uint16 {
	return s.icmpAddr.Port()
}

func (s *Session) activityMonitor() {
	for {
		delta := time.Until(s.activity.Add(cfgTimeout))
		if delta <= 0 {
			s.Purge()
			return
		}
		time.Sleep(delta)
	}
}

func getSession(tableId TableId, addr netip.AddrPort, c2s bool) *Session {
	var table map[netip.AddrPort]*Session
	if tableId == UdpTable {
		table = udpTable
	} else {
		table = icmpTable
	}
	mutex.RLock()
	session, ok := table[addr]
	mutex.RUnlock()
	if ok || !c2s {
		return session
	}
	session = new(Session)
	if tableId == UdpTable {
		session.udpConn = udpConn // use global UDPConn as server
		session.udpAddr = addr // client IP & port
		session.icmpIPAddr = net.IPAddr{IP: cfgIcmpAddr.AsSlice()}
	} else {
		var err error
		session.udpConn, err = net.DialUDP("udp", nil, net.UDPAddrFromAddrPort(cfgUdpAddrPort))
		if err != nil {
			log.Print(err)
			return nil
		}
		session.udpAddr = udpConnLocalAddrPort(session.udpConn) // client IP & port
		session.icmpAddr = addr // client IP & id
		session.icmpIPAddr = net.IPAddr{IP: addr.Addr().AsSlice()}
	}
	mutex.Lock()
	if _session, exist := table[addr]; exist {
		mutex.Unlock()
		return _session
	}
	if tableId == UdpTable {
		session.icmpAddr = netip.AddrPortFrom(cfgIcmpAddr, getIcmpId()) // server IP & id
	}
	udpTable[session.udpAddr] = session
	icmpTable[session.icmpAddr] = session
	mutex.Unlock()
	if tableId == IcmpTable {
		go udp2icmp(session.udpConn)
	}
	session.activity = time.Now()
	go session.activityMonitor()
	return session
}