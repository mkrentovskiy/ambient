package main

import (
	"fmt"
	"github.com/mkrentovskiy/ambient/devices"
	"github.com/tryphon/alsa-go"
	"time"
)

func sound_in(dest chan int) {
	handle := alsa.New()
	err := handle.Open("default", alsa.StreamTypeCapture, alsa.ModeBlock)
	if err != nil {
		fmt.Printf("Open failed. %s", err)
		return
	}

	handle.SampleFormat = alsa.SampleFormatU16LE
	handle.SampleRate = 11025
	handle.Channels = 1

	err = handle.ApplyHwParams()
	if err != nil {
		fmt.Printf("SetHwParams failed. %s\n", err)
		return
	}

	buflen := int(4096)
	buf := make([]uint8, buflen)
	for {
		n, err := handle.Read(buf)
		if err != nil {
			fmt.Printf("Read failed. %s\n", err)
		} else {
			max := uint16(buf[1])<<8 + uint16(buf[0])

			for i := 2; i < n; i += 2 {
				c := uint16(buf[i+1])<<8 + uint16(buf[i])
				if c > max {
					max = c
				}
			}
			dest <- int(max)
		}
	}
	handle.Close()
}

func pulse(src chan int, dest chan int, tres int, pulse_ms int) {
	k := 0
	for {
		go func() {
			v := <-src

			if v > tres {
				k = v
			}
		}()
		if k > 0 {
			dest <- k
			k = 0
		} else {
			dest <- 0
		}
		time.Sleep(time.Duration(pulse_ms) * time.Millisecond)
	}
}

func leds(source chan int, num_leds uint, min_val int, max_val int) {
	leds := make([]devices.RGB, num_leds)
	state := devices.InitWS2801(0, 1, num_leds)

	for {
		go func() {
			val := <-source
			for i := num_leds - 1; i > 0; i-- {
				leds[i] = leds[i-1]
			}
			leds[0] = val_to_color(val, min_val, max_val)
		}()
		state.Send(leds)
		time.Sleep(100 * time.Millisecond)
	}
	state.Done()
}

func val_to_color(val int, min_val int, max_val int) devices.RGB {
	if val < min_val || val > max_val {
		return devices.RGB{0, 0, 0}
	}

	var h float32 = float32(val-min_val) * 360 / float32(max_val-min_val)
	// fmt.Printf(" _ h = %f \n", h)

	if h >= 0 && h < 60 {
		return devices.RGB{byte(255), byte(255 * (1 - h/59)), byte(0)}
	} else if h >= 60 && h < 120 {
		return devices.RGB{byte(255 * (1 - h/59)), byte(255), byte(0)}
	} else if h >= 120 && h < 180 {
		return devices.RGB{byte(0), byte(255), byte(255 * (1 - h/59))}
	} else if h >= 180 && h < 240 {
		return devices.RGB{byte(0), byte(255 * (1 - h/59)), byte(255)}
	} else if h >= 240 && h < 300 {
		return devices.RGB{byte(255 * (1 - h/59)), byte(0), byte(255)}
	} else if h >= 300 && h <= 360 {
		return devices.RGB{byte(255), byte(0), byte(255 * (1 - h/59))}
	}
	return devices.RGB{0, 0, 0}
}

func main() {
	chp := make(chan int)
	chl := make(chan int)
	go pulse(chp, chl, 33523, 50)
	go leds(chl, 160, 33523, 65535)
	sound_in(chp)
}
