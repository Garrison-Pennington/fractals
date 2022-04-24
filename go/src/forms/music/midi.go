package music

import (
	"fmt"

	midi "gitlab.com/gomidi/midi/v2"
)

func Test() {
	var b = []byte{0x91, 0x3C, 0x78}

	// convert to Message type
	msg := midi.Message(b)

	var channel, key, velocity uint8
	if msg.GetNoteOn(&channel, &key, &velocity) {
		fmt.Printf("got %s: channel: %v key: %v, velocity: %v\n", msg.Type(), channel, key, velocity)
	}
}
