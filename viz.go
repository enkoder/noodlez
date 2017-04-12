package noodle

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"math/rand"
	"os"
)

const (
	MaxBrightness = 200
	Increasing    = 255
	Decreasing    = 1
	Up            = "up"
	Down          = "down"
	Left          = "left"
	Right         = "right"
	NumColors     = 24
)

var (
	Off     = colorful.Color{R: 0, G: 0, B: 0}
	Pallete []colorful.Color
)

func init() {
	var err error
	Pallete, err = colorful.HappyPalette(NumColors)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}

func RandomColor() colorful.Color {
	return Pallete[rand.Intn(NumColors)]
}

type Viz interface {
	Mutate(*Noodle)
	RefreshRate() float64
	String() string
}

type CircularViz struct {
	curStrip int
	curColor int
}

func NewCircularViz() Viz {
	viz := &CircularViz{}
	return viz
}

func (v *CircularViz) String() string {
	return fmt.Sprintf("CircularViz: curStrip: %d - curColor: %d", v.curStrip, v.curColor)
}

func (v *CircularViz) Mutate(n *Noodle) {
	color := RandomColor()
	for i := 0; i < LEDsPerStrip; i++ {
		n.Strips[v.curStrip].Pixels[i] = color
	}
	v.curStrip = (v.curStrip + 1) % NumStrips
}

func (v *CircularViz) RefreshRate() float64 {
	return .25
}

type SoftCircularViz struct {
	curStrip int
	curColor int
}

func NewSoftCircularViz() Viz {
	viz := &SoftCircularViz{}
	return viz
}

func (v *SoftCircularViz) String() string {
	return fmt.Sprintf("SoftCircularViz: curStrip: %d - curColor: %d", v.curStrip, v.curColor)
}

func (v *SoftCircularViz) Mutate(n *Noodle) {
	color := RandomColor()
	for strip := 0; strip < NumStrips; strip++ {
		for i := 0; i < LEDsPerStrip; i++ {
			if strip%2 == v.curStrip {
				n.Strips[strip].Pixels[i] = Off
			} else {
				n.Strips[strip].Pixels[i] = color
			}
		}
	}
	v.curStrip = (v.curStrip + 1) % 2
}

func (v *SoftCircularViz) RefreshRate() float64 {
	return 2
}

type SnakeViz struct {
	dir      string
	headx    uint
	heady    uint
	curColor colorful.Color
	count    int
}

func NewSnakeViz() Viz {
	viz := &SnakeViz{
		headx:    uint(rand.Intn(NumStrips)),
		heady:    uint(rand.Intn(LEDsPerStrip)),
		curColor: RandomColor(),
	}
	viz.dir = Up
	return viz
}

func (v *SnakeViz) GetNewLocation() (string, uint, uint) {
	if v.dir == Up {
		dir := []string{Left, Right, Up, Up, Up}[int(rand.Intn(5))]
		if dir == Left {
			return Left, (v.headx - 1) % NumStrips, v.heady
		} else if dir == Right {
			return Right, (v.headx + 1) % NumStrips, v.heady
		} else {
			return Up, v.headx, (v.heady + 1) % LEDsPerStrip
		}
	} else if v.dir == Down {
		dir := []string{Left, Right, Down, Up, Up}[int(rand.Intn(5))]
		if dir == Left {
			return Left, (v.headx - 1) % NumStrips, v.heady
		} else if dir == Right {
			return Right, (v.headx + 1) % NumStrips, v.heady
		} else {
			return Down, v.headx, (v.heady - 1) % LEDsPerStrip
		}
	} else if v.dir == Left {
		dir := []string{Up, Down, Left, Up, Up}[int(rand.Intn(5))]
		if dir == Left {
			return Left, (v.headx - 1) % NumStrips, v.heady
		} else if dir == Down {
			return Down, v.headx, (v.heady - 1) % LEDsPerStrip
		} else {
			return Left, v.headx, (v.heady + 1) % LEDsPerStrip
		}
	} else if v.dir == Right {
		dir := []string{Down, Up, Right, Up, Up}[int(rand.Intn(5))]
		if dir == Right {
			return Right, (v.headx + 1) % NumStrips, v.heady
		} else if dir == Down {
			return Down, v.headx, (v.heady - 1) % LEDsPerStrip
		} else {
			return Up, v.headx, (v.heady + 1) % LEDsPerStrip
		}
	}
	// really should return err here
	return "", 0, 0
}

func (v *SnakeViz) String() string {
	return fmt.Sprintf("SnakeViz: head: [%2d, %2d] - dir: %s",
		v.headx,
		v.heady,
		v.dir)
}

func (v *SnakeViz) Mutate(n *Noodle) {
	dir, x, y := v.GetNewLocation()
	n.Strips[x].Pixels[y] = v.curColor
	v.dir = dir
	v.headx = x
	v.heady = y
	v.count += 1
	if v.count > 100 {
		v.curColor = RandomColor()
		v.count = 0
	}
}

func (v *SnakeViz) RefreshRate() float64 {
	return .05
}

type SpiralViz struct {
	curStrip uint8
	curLED   uint8
	curColor colorful.Color
	count    int
}

func NewSpiralViz() Viz {
	return &SpiralViz{
		curLED:   0,
		curStrip: 0,
		count:    0,
		curColor: RandomColor(),
	}
}

func (v *SpiralViz) String() string {
	return fmt.Sprintf(
		"SpiralViz: %d.%d[%f, %f, %f]",
		v.curStrip,
		v.curLED,
		v.curColor.R,
		v.curColor.G,
		v.curColor.B,
	)
}

func (v *SpiralViz) Mutate(n *Noodle) {
	v.curStrip += 1
	if v.curStrip == NumStrips {
		v.curStrip = 0
		v.curLED = (v.curLED + 1) % LEDsPerStrip
	}

	v.count += 1
	n.Strips[v.curStrip].Pixels[v.curLED] = v.curColor
	if v.count > 20 {
		v.curColor = RandomColor()
		v.count = 0
	}
}

func (v SpiralViz) RefreshRate() float64 {
	return .04
}

type VertViz struct {
	mainColor colorful.Color
	midColor  colorful.Color
	midCount  int
	midPos    int
	dir       int
}

func NewVertViz() Viz {
	return &VertViz{
		mainColor: RandomColor(),
		midColor:  RandomColor(),
		midCount:  6,
		midPos:    LEDsPerStrip / 2,
		dir:       1,
	}
}

func (v *VertViz) String() string {
	return fmt.Sprintf("VertViz: midPos=%d", v.midPos)
}

func (v *VertViz) Mutate(n *Noodle) {
	for _, s := range n.Strips {
		for i := 0; i < LEDsPerStrip; i++ {
			s.SetColor(v.mainColor)
		}
	}
	for _, s := range n.Strips {
		for i := v.midPos - (v.midCount / 2); i < v.midPos+(v.midCount/2); i++ {
			s.Pixels[i] = v.midColor
		}
	}

	v.midPos += v.dir
	// top
	if v.midPos == LEDsPerStrip-(v.midCount/2) {
		v.dir = -1
		v.mainColor = v.midColor
		v.midColor = RandomColor()
		// bottom
	} else if v.midPos == (v.midCount / 2) {
		v.dir = 1
		v.mainColor = v.midColor
		v.midColor = RandomColor()
	}
}

func (v *VertViz) RefreshRate() float64 {
	return .05
}
