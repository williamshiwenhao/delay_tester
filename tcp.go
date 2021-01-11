package main

import (
	"encoding/binary"
	"net"
)

// TCPConnection connection, can be server or client
type TCPConnection struct {
	socket      net.Conn
	receiveChan chan []byte
}

// CreateClient create a tcp connection to addr
func CreateClient(addr string) *TCPConnection {
	t := &TCPConnection{}
	sock, err := net.Dial("tcp", addr)
	if err != nil {
		logger.Fatalf("Create tcp connection failed, err=%+v", err)
	}
	t.socket = sock
	return t
}

// StartReceive start listen from socket
func (t *TCPConnection) StartReceive() <-chan []byte {
	if t.receiveChan != nil {
		logger.Fatalf("Try to start receive an socket already receving")
	}
	t.receiveChan = make(chan []byte, ChanSize)
	go t.receive()
	return t.receiveChan
}

// Send data out
func (t *TCPConnection) Send(data []byte) {
	outputBuffer := packetLen(data)
	n, err := t.socket.Write(outputBuffer)
	if err != nil {
		logger.Fatalf("Send error, err=%+v", err)
	} else if n != len(outputBuffer) {
		logger.Warnf("Sent length != data length")
	}
}

func packetLen(data []byte) []byte {
	const kUint32Size = 4
	buffer := make([]byte, len(data)+kUint32Size)
	binary.LittleEndian.PutUint32(buffer, uint32(len(data)))
	copy(data[kUint32Size:], data)
	return buffer
}

func (t *TCPConnection) receive() {
	var len uint32
	lenBuffer := make([]byte, 4)
	for {
		// Read packet length
		t.receiveAll(lenBuffer)
		len = binary.LittleEndian.Uint32(lenBuffer)
		outputBuffer := make([]byte, len)
		t.receiveAll(outputBuffer)
		t.receiveChan <- outputBuffer
	}
}

func (t *TCPConnection) receiveAll(data []byte) {
	var base int
	for base < len(data) {
		n, err := t.socket.Read(data[base:])
		if err != nil {
			logger.Warnf("Read from socket failed, err=%+v", err)
		}
		base += n
	}
}
