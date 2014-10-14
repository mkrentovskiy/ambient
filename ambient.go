package main

import (
	"github.com/mkrentovskiy/ambient/devices"
	"time"
)

func main() {
	num_leds := uint(160)
	leds := make([]devices.RGB, num_leds)

	state := devices.InitWS2801(0, 1, num_leds)
	for j := byte(0); j < 255; j++ {
		for i, _ := range leds {
			if i%6 == 0 {
				leds[i] = devices.RGB{j % 32, 0, 0}
			} else if i%6 == 1 {
				leds[i] = devices.RGB{0, j % 32, 0}
			} else if i%6 == 2 {
				leds[i] = devices.RGB{0, 0, j % 32}
			} else if i%6 == 3 {
				leds[i] = devices.RGB{j % 32, j % 64, 0}
			} else if i%6 == 4 {
				leds[i] = devices.RGB{j % 32, 0, j % 64}
			} else if i%6 == 6 {
				leds[i] = devices.RGB{0, j % 32, j % 64}
			}

		}
		time.Sleep(120 * time.Millisecond)
		state.Send(leds)
	}
	state.Done()
}
