package main

// Conn an connection
type Conn interface {
	Send([]byte)
	StartReceive() <-chan []byte
}
