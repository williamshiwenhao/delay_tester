package main

import (
	"encoding/binary"
	"log"
	"time"
)

type receiver struct {
	conn           Conn
	receivedPacket uint64
	maxDelay       time.Duration
	minDelay       time.Duration
	averageDelay   time.Duration
	totalDelay     time.Duration
	durations      chan time.Duration
}

func NewReceiver(conn Conn) *receiver {

	return &receiver{
		conn:      conn,
		minDelay:  time.Minute,
		durations: make(chan time.Duration, 65536),
	}
}

func (r *receiver) getDuration(data []byte) time.Duration {

	packTime := getPacketTime(data)
	duration := time.Now().Sub(packTime)
	return duration
}

func (r *receiver) Run() {
	reportTicker := time.NewTicker(time.Second)
	receiveChan := r.conn.StartReceive()
	var nowMax, nowMin time.Duration
	nowMin = time.Minute
	var nowReceived uint64
	for {
		select {
		case <-reportTicker.C:
			logger.Errorf("[Now] %v packets, min delay = %v, max delay = %v", nowReceived, nowMin, nowMax)
			if nowMin < r.minDelay {
				r.minDelay = nowMin
			}
			if nowMax > r.maxDelay {
				r.maxDelay = nowMax
			}
			r.receivedPacket += nowReceived
			nowMax = 0
			nowMin = time.Minute
			nowReceived = 0
			logger.Warnf("[Total] %v packets, min delay = %v, max delay = %v", r.receivedPacket, r.minDelay, r.maxDelay)
		case data := <-receiveChan:
			duration := r.getDuration(data)
			nowReceived++
			if duration < nowMin {
				nowMin = duration
			}
			if duration > nowMax {
				nowMax = duration
			}
		}
	}
}

func getPacketTime(data []byte) time.Time {
	length := binary.LittleEndian.Uint16(data)
	var t time.Time
	err := t.UnmarshalBinary(data[2 : 2+length])
	if err != nil {
		log.Fatalf("Unmarshal binary failed, err=%+v", err)
	}
	return t
}
