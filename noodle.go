package noodlez

import (
	"fmt"
	"time"

	"github.com/currantlabs/gatt"
	"github.com/enkoder/noodlez/goodle"
	"github.com/kellydunn/go-opc"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mrmorphic/hwio"
)

const (
	NumStrips      = 4
	LEDsPerStrip   = 37
	LEDsPerChannel = 64
	TotalLeds      = NumStrips * LEDsPerStrip
)

type Strip struct {
	Pixels []colorful.Color
}

func (s *Strip) SetColor(color colorful.Color) {
	for i := 0; i < LEDsPerStrip; i++ {
		s.Pixels[i] = color
	}
}

type Noodle struct {
	button        hwio.Pin
	client        *opc.Client
	message       *opc.Message
	Strips        []*Strip
	MaxBrightness uint8
	vizs          []Viz
	prevViz       int
	curViz        int
}

func (n *Noodle) HandleWrite(r gatt.Request, data []byte) byte {
	fmt.Println("Wrote:", string(data))
	return gatt.StatusSuccess
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

	strips := make([]*Strip, NumStrips)
	for i := 0; i < NumStrips; i++ {
		strips[i] = &Strip{}
		strips[i].Pixels = make([]colorful.Color, LEDsPerStrip)
	}

	vizs := []Viz{
		NewSoftCircularViz(),
		NewSpiralViz(),
		NewSparkleViz(),
		NewVertSwapViz(),
		NewCircularViz(),
		NewVertViz(),
		NewSnakeViz(),
		NewLaserViz(),
	}

	noodle := &Noodle{
		button:  button,
		client:  client,
		message: message,
		Strips:  strips,
		vizs:    vizs,
		prevViz: 0,
		curViz:  0,
	}

	err = goodle.InitBluetooth(noodle.HandleWrite)
	if err != nil {
		return nil, fmt.Errorf("Error during InitBluetooth: %v\n", err)
	}

	return noodle, nil
}

func (n *Noodle) NextViz() {
	n.Off()
	n.prevViz = n.curViz
	n.curViz = (n.curViz + 1) % (len(n.vizs) - 1)
}

func (n *Noodle) FireTheLaser() {
	n.prevViz = n.curViz
	n.curViz = len(n.vizs) - 1
}

func (n *Noodle) StopHumpingTheLaser() {
	n.curViz = n.prevViz
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

	var viz Viz
	for {
		// sets the viz on each loop. yolo
		viz = n.vizs[n.curViz]

		if time.Since(lastRender).Seconds() > viz.RefreshRate() {
			fmt.Println(viz.String())
			viz.Mutate(n)
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

			// Someone just pressed a button, store that state
			if buttval && !prevButtVal {
				lastButtonPress = time.Now()
				prevButtVal = true

				// just depressed button, check for short press
			} else if !buttval && prevButtVal && !changed {
				prevButtVal = false
				// Shoot lights viz
				if time.Since(lastButtonPress) > 10*time.Millisecond {
					n.FireTheLaser()
				}

				// If its been pressed for a bit
			} else if buttval && prevButtVal {
				// change viz
				if time.Since(lastButtonPress) > 1*time.Second && !changed {
					changed = true
					n.NextViz()
				}

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
			n.message.SetPixelColor((s*LEDsPerChannel)+led,
				uint8(n.Strips[s].Pixels[led].R*MaxBrightness),
				uint8(n.Strips[s].Pixels[led].G*MaxBrightness),
				uint8(n.Strips[s].Pixels[led].B*MaxBrightness))
		}
	}
	return n.client.Send(n.message)
}

// Turns off all leds
func (n *Noodle) Off() error {
	for s := range n.Strips {
		for led := 0; led < LEDsPerStrip; led++ {
			n.Strips[s].Pixels[led].R = 0
			n.Strips[s].Pixels[led].G = 0
			n.Strips[s].Pixels[led].B = 0
		}
	}
	return n.Render()
}

func (n *Noodle) ButtonPressed() (bool, error) {
	value, err := hwio.DigitalRead(n.button)
	return value == 1, err
}
