package bus

import (
	"fmt"
	"github.com/luismesas/goPi/ioctl"
	"os"
	"unsafe"
)

const SPIDEV = "/dev/spidev"
const SPI_IOC_MAGIC = 107
const (
	SPI_HARDWARE_ADDR = 0
	SPI_BUS           = 0
	SPI_CHIP          = 0
	SPI_DELAY         = 0
)

type SPI_IOC_TRANSFER struct {
	txBuf       uint64
	rxBuf       uint64
	length      uint32
	speedHz     uint32
	delayUsecs  uint16
	bitsPerWord uint8
	csChange    uint8
	pad         uint32
}

type SPIDev struct {
	Bus  int      // 0
	Chip int      // 0
	file *os.File // nil

	mode  uint8
	bpw   uint8
	speed uint32
}

// An SPI Device at /dev/spi<bus>.<chip_select>.
func NewSPIDev(bus int, chipSelect int) *SPIDev {
	spi := new(SPIDev)
	spi.Bus = bus
	spi.Chip = chipSelect

	return spi
}

// Opens SPI device
func (spi *SPIDev) Open() error {
	spiDevice := fmt.Sprintf("%s%d.%d", SPIDEV, spi.Bus, spi.Chip)

	var err error
	spi.file, err = os.OpenFile(spiDevice, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("I can't see %s. Have you enabled the SPI module?", spiDevice)
	}

	return nil
}

// Closes SPI device
func (spi *SPIDev) Close() error {
	err := spi.file.Close()
	if err != nil {
		return fmt.Errorf("Error closing spi", err)
	}
	return nil
}

// Sends bytes over SPI channel
func (spi *SPIDev) Send(bytes_to_send []byte) error {
	wBuffer := bytes_to_send
	rBuffer := [3]byte{}

	// generates message
	transfer := SPI_IOC_TRANSFER{}
	transfer.txBuf = uint64(uintptr(unsafe.Pointer(&wBuffer)))
	transfer.rxBuf = uint64(uintptr(unsafe.Pointer(&rBuffer)))
	transfer.length = uint32(unsafe.Sizeof(wBuffer))
	transfer.delayUsecs = SPI_DELAY
	transfer.bitsPerWord = spi.bpw
	transfer.speedHz = spi.speed

	// sends message over SPI
	err := ioctl.IOCTL(spi.file.Fd(), SPI_IOC_MESSAGE(1), uintptr(unsafe.Pointer(&transfer)))
	if err != nil {
		return fmt.Errorf("Error on sending: %s\n", err)
	}

	return nil
}

// Sends bytes over SPI channel and returns []byte response
func (spi *SPIDev) SendReplay(bytes_to_send [3]byte) ([]byte, error) {
	wBuffer := bytes_to_send
	rBuffer := [3]byte{}

	// generates message
	transfer := SPI_IOC_TRANSFER{}
	transfer.txBuf = uint64(uintptr(unsafe.Pointer(&wBuffer)))
	transfer.rxBuf = uint64(uintptr(unsafe.Pointer(&rBuffer)))
	transfer.length = uint32(unsafe.Sizeof(wBuffer))
	transfer.delayUsecs = SPI_DELAY
	transfer.bitsPerWord = spi.bpw
	transfer.speedHz = spi.speed

	// sends message over SPI
	err := ioctl.IOCTL(spi.file.Fd(), SPI_IOC_MESSAGE(1), uintptr(unsafe.Pointer(&transfer)))
	if err != nil {
		return nil, fmt.Errorf("Error on sending: %s\n", err)
	}

	// generates a valid response
	ret := make([]byte, unsafe.Sizeof(rBuffer))
	for i := range ret {
		ret[i] = rBuffer[i]
	}

	return ret, nil
}

// Write bytes over SPI
func (spi *SPIDev) Write(bytes []byte) error {
	_, err := spi.file.Write(bytes)
	return err
}

func (spi *SPIDev) SetMode(mode uint8) error {
	spi.mode = mode
	err := ioctl.IOCTL(spi.file.Fd(), SPI_IOC_WR_MODE(), uintptr(unsafe.Pointer(&mode)))
	if err != nil {
		return fmt.Errorf("Error setting mode: %s\n", err)
	}
	return nil
}

func (spi *SPIDev) SetBitsPerWord(bpw uint8) error {
	spi.bpw = bpw
	err := ioctl.IOCTL(spi.file.Fd(), SPI_IOC_WR_BITS_PER_WORD(), uintptr(unsafe.Pointer(&bpw)))
	if err != nil {
		return fmt.Errorf("Error setting bits per word: %s\n", err)
	}
	return nil
}

func (spi *SPIDev) SetSpeed(speed uint32) error {
	spi.speed = speed
	err := ioctl.IOCTL(spi.file.Fd(), SPI_IOC_WR_MAX_SPEED_HZ(), uintptr(unsafe.Pointer(&speed)))
	if err != nil {
		return fmt.Errorf("Error setting speed: %s\n", err)
	}
	return nil
}

// Read of SPI mode (SPI_MODE_0..SPI_MODE_3)
func SPI_IOC_RD_MODE() uintptr {
	return ioctl.IOR(SPI_IOC_MAGIC, 1, 1)
}

// Write of SPI mode (SPI_MODE_0..SPI_MODE_3)
func SPI_IOC_WR_MODE() uintptr {
	return ioctl.IOW(SPI_IOC_MAGIC, 1, 1)
}

// Read SPI bit justification
func SPI_IOC_RD_LSB_FIRST() uintptr {
	return ioctl.IOR(SPI_IOC_MAGIC, 2, 1)
}

// Write SPI bit justification
func SPI_IOC_WR_LSB_FIRST() uintptr {
	return ioctl.IOW(SPI_IOC_MAGIC, 2, 1)
}

// Read SPI device word length (1..N)
func SPI_IOC_RD_BITS_PER_WORD() uintptr {
	return ioctl.IOR(SPI_IOC_MAGIC, 3, 1)
}

// Write SPI device word length (1..N)
func SPI_IOC_WR_BITS_PER_WORD() uintptr {
	return ioctl.IOW(SPI_IOC_MAGIC, 3, 1)
}

// Read SPI device default max speed hz
func SPI_IOC_RD_MAX_SPEED_HZ() uintptr {
	return ioctl.IOR(SPI_IOC_MAGIC, 4, 4)
}

// Write SPI device default max speed hz
func SPI_IOC_WR_MAX_SPEED_HZ() uintptr {
	return ioctl.IOW(SPI_IOC_MAGIC, 4, 4)
}

// Write custom SPI message
func SPI_IOC_MESSAGE(n uintptr) uintptr {
	return ioctl.IOW(SPI_IOC_MAGIC, 0, uintptr(SPI_MESSAGE_SIZE(n)))
}

func SPI_MESSAGE_SIZE(n uintptr) uintptr {
	if (n * unsafe.Sizeof(SPI_IOC_TRANSFER{})) < (1 << ioctl.IOC_SIZEBITS) {
		return (n * unsafe.Sizeof(SPI_IOC_TRANSFER{}))
	}
	return 0
}
