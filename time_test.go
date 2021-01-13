package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	now := time.Now()
	buffer := make([]byte, 100)
	putNowTime(now, buffer)
	getTime := getPacketTime(buffer)
	assert.True(t, now.Equal(getTime))
}
