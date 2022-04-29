package music

import smf "gitlab.com/gomidi/midi/v2/smf"

var NOTE_VALUES map[uint16]func(*smf.Track, smf.MetricTicks, uint8, uint8, uint8) = map[uint16]func(*smf.Track, smf.MetricTicks, uint8, uint8, uint8){
	1:  Play16th,
	2:  Play8th,
	4:  Play4th,
	8:  PlayHalf,
	16: PlayWhole,
}
