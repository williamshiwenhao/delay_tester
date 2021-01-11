package main

import (
	"encoding/binary"
	"time"
)

type Sender struct {
	conn        Conn
	tick        *time.Ticker
	duration    time.Duration
	packLen     int
	packPreTick int
}

const fill uint8 = 0xff

// NewSender create a sender
func NewSender(conn Conn) *Sender {
	return &Sender{
		conn:        conn,
		duration:    time.Duration(Config.TickMs) * time.Millisecond,
		packLen:     Config.PacketLength,
		packPreTick: Config.PacketPreTick,
	}
}

// Run run the sender
func (s *Sender) Run() {
	if s.tick != nil {
		logger.Fatalf("Try to start a sender already running")
	}
	packet := make([]byte, s.packLen)
	// Fill the packet
	for i := 0; i < len(packet); i++ {
		packet[i] = fill
	}

	// Init and start ticker
	s.tick = time.NewTicker(s.duration)
	for ; ; <-s.tick.C {
		for i := 0; i < Config.PacketPreTick; i++ {
			putNowTime(packet)
			s.conn.Send(packet)
		}
	}
}

func putNowTime(data []byte) {
	now := time.Now()
	binary.BigEndian.PutUint64(data, uint64(now.UnixNano()))
}
