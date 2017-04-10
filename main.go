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

type Noodle struct {
	b hwio.Pin
	c *opc.Client
	m *opc.Message
}

func (n *Noodle) Solid(r uint8, g uint8, b uint8) error {
	for s := 0; s < Strips; s++ {
		for led := 0; led < LEDsPerStrip; led++ {
			fmt.Println(s*LEDsPerChannel + led)
			n.m.SetPixelColor((s*LEDsPerChannel)+led, r, g, b)
		}

	}
	return n.c.Send(n.m)
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

	return &Noodle{b, c, m}, nil
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
