package main

import "unsafe"

// checksum is endian-neutral thanks to its design (carries in a circle)
// assume 2 <= offset < len(data) <= 65536
func fillChecksum(data []byte, offset uint) {
	_ = data[offset+1]
	p := unsafe.Pointer(&data[0])
	l := len(data) - 1
	s := uint32(0)
	for i := 0; i < l; i += 2 {
		s += uint32(*((*uint16)(unsafe.Add(p, i))))
	}
	if l & 1 == 0 {
		s += uint32(*((*uint16)(unsafe.Add(p, l))))
	}
	s = s>>16 + s&0xffff
	s = s + s>>16
	*(*uint16)(unsafe.Add(p, offset)) = ^uint16(s)
}
