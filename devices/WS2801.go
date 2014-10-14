package devices

import (
	"github.com/mkrentovskiy/ambient/bus"
)

type WS2801 struct {
	dev    *bus.SPIDev
	num    uint
	buffer []byte
}

type RGB struct {
	R, G, B byte
}

func InitWS2801(bus_id int, chip_id int, n uint) *WS2801 {
	state := new(WS2801)
	state.num = n
	state.buffer = make([]byte, n*3)
	state.dev = bus.NewSPIDev(bus_id, chip_id)
	state.dev.Open()
	state.dev.SetMode(0)
	state.dev.SetBitsPerWord(8)
	state.dev.SetSpeed(475000)
	state.Clear()

	return state
}

func (state *WS2801) Clear() error {
	for i, _ := range state.buffer {
		state.buffer[i] = 0
	}
	return state.dev.Write(state.buffer)
}

func (state *WS2801) SetBusSpeed(s uint32) {
	state.dev.SetSpeed(s)
}

func (state *WS2801) Send(colors []RGB) error {
	j := 0
	for i := uint(0); i < state.num*3; i += 3 {
		state.buffer[i] = colors[j].B
		state.buffer[i+1] = colors[j].G
		state.buffer[i+2] = colors[j].R
		j++
	}
	return state.dev.Write(state.buffer)
}

func (state *WS2801) Done() {
	state.dev.Close()
}
