package main

import "golang.org/x/net/ipv4"

type ICMPType = ipv4.ICMPType

// https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-types
// net/ipv4/icmp.c, icmp_pointers[]
// net/netfilter/nf_conntrack_proto_icmp.c, invmap[]
const (
	ICMPTypeEchoReply              ICMPType = 0 // ping_rcv, invmap
	ICMPTypeUnassigned_1           ICMPType = 1 // icmp_discard
	ICMPTypeUnassigned_2           ICMPType = 2 // icmp_discard
	ICMPTypeDestinationUnreachable ICMPType = 3 // icmp_unreach
	ICMPTypeSourceQuench           ICMPType = 4 // icmp_unreach
	ICMPTypeRedirect               ICMPType = 5 // icmp_redirect
	ICMPTypeAlternateHostAddress   ICMPType = 6 // icmp_discard
	ICMPTypeUnassigned_7           ICMPType = 7 // icmp_discard
	ICMPTypeEchoRequest            ICMPType = 8 // icmp_echo, invmap
	ICMPTypeRouterAdvertisement    ICMPType = 9 // icmp_discard
	ICMPTypeRouterSolicitation     ICMPType = 10 // icmp_discard
	ICMPTypeTimeExceed             ICMPType = 11 // icmp_unreach
	ICMPTypeParameterProblem       ICMPType = 12 // icmp_unreach
	ICMPTypeTimestampRequest       ICMPType = 13 // icmp_timestamp, invmap
	ICMPTypeTimestampReply         ICMPType = 14 // icmp_discard, invmap
	ICMPTypeInformationRequest     ICMPType = 15 // icmp_discard, invmap
	ICMPTypeInformationReply       ICMPType = 16 // icmp_discard, invmap
	ICMPTypeAddressMaskRequest     ICMPType = 17 // icmp_discard, invmap
	ICMPTypeAddressMaskReply       ICMPType = 18 // icmp_discard
)
