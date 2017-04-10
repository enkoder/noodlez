package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kellydunn/go-opc"
	"github.com/mrmorphic/hwio"
)

const (
	ButtonPin      = "gpio17"
	Strips         = 4
	LEDsPerStrip   = 37
	LEDsPerChannel = 64
	TotalLeds      = Strips * LEDsPerStrip
)

type Strip struct {
	R [LEDsPerStrip]uint8
	G [LEDsPerStrip]uint8
	B [LEDsPerStrip]uint8
}

type Noodle struct {
	b      hwio.Pin
	c      *opc.Client
	m      *opc.Message
	Strips []Strip
}

func (n *Noodle) Render() error {
	for s := range n.Strips {
		for led := 0; led < LEDsPerStrip; led++ {
			n.m.SetPixelColor((s*LEDsPerChannel)+led, n.Strips[s].R[led], n.Strips[s].G[led], n.Strips[s].B[led])
		}
	}
	return n.c.Send(n.m)
}

func (n *Noodle) Solid(r uint8, g uint8, b uint8) error {
	for s := range n.Strips {
		for led := 0; led < LEDsPerStrip; led++ {
			n.Strips[s].R[led] = r
			n.Strips[s].G[led] = g
			n.Strips[s].B[led] = b
		}
	}
	return n.Render()
}

// Turns off all leds
func (n *Noodle) Off() error {
	return n.Solid(0, 0, 0)
}

func (n *Noodle) Red() error {
	return n.Solid(255, 0, 0)
}

func (n *Noodle) Green() error {
	return n.Solid(0, 255, 0)
}

func (n *Noodle) Blue() error {
	return n.Solid(0, 0, 255)
}

func (n *Noodle) ButtonPressed() (bool, error) {
	value, err := hwio.DigitalRead(n.b)
	return value == 1, err
}

func NewNoodle(button_gpio string) (*Noodle, error) {
	b, err := hwio.GetPin(ButtonPin)
	if err != nil {
		return nil, fmt.Errorf("Error during button init: %v\n", err)
	}

	err = hwio.PinMode(b, hwio.INPUT)
	if err != nil {
		return nil, fmt.Errorf("Error during button mode set: %v\n", err)
	}

	// Create a client
	c := opc.NewClient()
	err = c.Connect("tcp", "localhost:7890")
	if err != nil {
		return nil, err
	}

	var m *opc.Message
	m = opc.NewMessage(0)
	m.SetLength(uint16(LEDsPerChannel * Strips * 3))

	strips := make([]Strip, Strips)
	for i := 0; i < Strips; i++ {
		strips[i] = Strip{}
	}

	return &Noodle{b, c, m, strips}, nil
}

func main() {
	noodle, err := NewNoodle(ButtonPin)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	err = noodle.Off()
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	for {
		time.Sleep(time.Millisecond * 500)

		value, err := noodle.ButtonPressed()
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Error during DigitalRead: %v\n", err))
			os.Exit(1)
		}

		if value {
			noodle.Blue()
		} else {
			noodle.Red()
		}

		fmt.Printf("Button Value: %t\n", value)
	}
}
