package main

type Direction bool
const (
	Send Direction = true
	Recv Direction = false
)

func isClientToServer(dir Direction) bool {
	return cfgRole == Client && dir == Send || cfgRole == Server && dir == Recv
}
