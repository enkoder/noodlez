package noodle

import (
	"fmt"
	"time"

	"github.com/kellydunn/go-opc"
	"github.com/mrmorphic/hwio"
)

const (
	NumStrips      = 4
	LEDsPerStrip   = 37
	LEDsPerChannel = 64
	TotalLeds      = NumStrips * LEDsPerStrip
)

type Strip struct {
	Pixels []Pixel
}

type Noodle struct {
	button        hwio.Pin
	client        *opc.Client
	message       *opc.Message
	Strips        []Strip
	MaxBrightness uint8
	curViz        Viz
	prevViz       Viz
	vizs          []Viz
	vizi          int
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
	message.SetLength(uint16(LEDsPerChannel * NumStrips * 3))

	strips := make([]Strip, NumStrips)
	for i := 0; i < NumStrips; i++ {
		strips[i] = Strip{}
		strips[i].Pixels = make([]Pixel, LEDsPerStrip)
	}

	vizs := []Viz{
		NewSoftCircularViz(),
		NewSpiralViz(10),
		NewCircularViz(),
		NewSnakeViz()}

	return &Noodle{
		button:  button,
		client:  client,
		message: message,
		Strips:  strips,
		vizs:    vizs,
		prevViz: vizs[0],
		curViz:  vizs[0],
		vizi:    0,
	}, nil
}

// VizLoop runs forever as the main run thread calling Mutate on the Viz's
// and checking for input like buttons and maybe bluetooth
func (n *Noodle) VizLoop() {
	var err error
	lastRender := time.Now()
	lastButtonPress := time.Now()
	lastButtonRead := time.Now()
	buttval := false
	prevButtVal := false
	changed := false

	for {
		if time.Since(lastRender).Seconds() > n.curViz.RefreshRate() {
			fmt.Println(n.curViz.String())
			n.curViz.Mutate(n)
			n.Render()
			lastRender = time.Now()
		}

		// Read button press
		if time.Since(lastButtonRead) > time.Millisecond*10 {
			lastButtonRead = time.Now()
			buttval, err = n.ButtonPressed()
			if err != nil {
				fmt.Printf("Error during button read: %v", err)
				continue
			}

			// Someone just pressed a button
			if buttval && !prevButtVal {
				lastButtonPress = time.Now()
				prevButtVal = true
				// If its been pressed for a bit
			} else if buttval && prevButtVal {
				// change viz
				if time.Since(lastButtonPress) > 1*time.Second && !changed {
					changed = true
					n.prevViz = n.vizs[n.vizi]
					n.vizi = (n.vizi + 1) % len(n.vizs)
					n.curViz = n.vizs[n.vizi]
					n.Off()
				}
				// Shoot lights viz
			} else if !buttval && prevButtVal && time.Since(lastButtonRead) > 100*time.Millisecond {

			} else {
				changed = false
				prevButtVal = false
			}
		}
	}
}

func (n *Noodle) Render() error {
	for s := range n.Strips {
		for led := 0; led < LEDsPerStrip; led++ {
			n.message.SetPixelColor((s*LEDsPerChannel)+led, n.Strips[s].Pixels[led].R, n.Strips[s].Pixels[led].G, n.Strips[s].Pixels[led].B)
		}
	}
	return n.client.Send(n.message)
}

func (n *Noodle) Solid(r uint8, g uint8, b uint8) error {
	for s := range n.Strips {
		for led := 0; led < LEDsPerStrip; led++ {
			n.Strips[s].Pixels[led].R = r
			n.Strips[s].Pixels[led].G = g
			n.Strips[s].Pixels[led].B = b
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
