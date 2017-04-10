package noodle

import (
	"fmt"

	"github.com/kellydunn/go-opc"
	"github.com/mrmorphic/hwio"
)

const (
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
	button  hwio.Pin
	client  *opc.Client
	message *opc.Message
	Strips  []Strip
	curViz  Viz
	prevViz Viz
}

func (n *Noodle) Render() error {
	for s := range n.Strips {
		for led := 0; led < LEDsPerStrip; led++ {
			n.message.SetPixelColor((s*LEDsPerChannel)+led, n.Strips[s].R[led], n.Strips[s].G[led], n.Strips[s].B[led])
		}
	}
	return n.client.Send(n.message)
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
	value, err := hwio.DigitalRead(n.button)
	return value == 1, err
}

func NewNoodle(button_gpio string) (*Noodle, error) {
	button, err := hwio.GetPin(button_gpio)
	if err != nil {
		return nil, fmt.Errorf("Error during button init: %v\n", err)
	}

	err = hwio.PinMode(button, hwio.INPUT)
	if err != nil {
		return nil, fmt.Errorf("Error during button mode set: %v\n", err)
	}

	// Create a client
	client := opc.NewClient()
	err = client.Connect("tcp", "localhost:7890")
	if err != nil {
		return nil, err
	}

	var message *opc.Message
	message = opc.NewMessage(0)
	message.SetLength(uint16(LEDsPerChannel * Strips * 3))

	strips := make([]Strip, Strips)
	for i := 0; i < Strips; i++ {
		strips[i] = Strip{}
	}

	return &Noodle{
		button:  button,
		client:  client,
		message: message,
		Strips:  strips,
		prevViz: nil,
		curViz:  nil,
	}, nil
}
