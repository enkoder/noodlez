package noodle

import (
	"fmt"
	"math/rand"
)

const (
	MaxBrightness = 200
	Increasing    = 255
	Decreasing    = 1
	Up            = "up"
	Down          = "down"
	Left          = "left"
	Right         = "right"
)

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
	for i := 0; i < LEDsPerStrip; i++ {
		n.Strips[v.curStrip].Pixels[i] = NamedColors[v.curColor]
	}
	v.curStrip = (v.curStrip + 1) % NumStrips
	v.curColor = (v.curColor + 1) % len(NamedColors)
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
	for strip := 0; strip < NumStrips; strip++ {
		for i := 0; i < LEDsPerStrip; i++ {
			if strip%2 == v.curStrip {
				n.Strips[strip].Pixels[i] = Off
			} else {
				n.Strips[strip].Pixels[i] = NamedColors[v.curColor]
			}
		}
	}
	v.curStrip = (v.curStrip + 1) % 2
	v.curColor = (v.curColor + 1) % len(NamedColors)
}

func (v *SoftCircularViz) RefreshRate() float64 {
	return 2
}

type SnakeViz struct {
	dir      string
	headx    uint
	heady    uint
	curColor Pixel
	count    int
}

func NewSnakeViz() Viz {
	viz := &SnakeViz{
		headx:    uint(rand.Intn(NumStrips)),
		heady:    uint(rand.Intn(LEDsPerStrip)),
		curColor: Red,
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
		v.curColor = NamedColors[rand.Intn(len(NamedColors))]
		v.count = 0
	}
}

func (v *SnakeViz) RefreshRate() float64 {
	return .05
}

type SpiralViz struct {
	curStrip uint8
	curLED   uint8
	stateR   uint8
	stateG   uint8
	stateB   uint8
	r        uint8
	g        uint8
	b        uint8
	step     uint8
}

func NewSpiralViz(step uint8) Viz {
	return &SpiralViz{
		r:      MaxBrightness,
		g:      0,
		b:      0,
		stateR: Decreasing,
		stateG: Increasing,
		stateB: 0,
		step:   step,
	}
}

func (v *SpiralViz) String() string {
	return fmt.Sprintf(
		"SpiralViz: %d.%d[%d, %d, %d]",
		v.curStrip,
		v.curLED,
		v.r,
		v.g,
		v.b,
	)
}

func (v *SpiralViz) Mutate(n *Noodle) {
	v.curStrip += 1
	if v.curStrip == NumStrips {
		v.curStrip = 0
		v.curLED = (v.curLED + 1) % LEDsPerStrip
	}

	if v.stateR == Increasing && v.r+v.step >= MaxBrightness {
		v.stateR = Decreasing
		v.stateG = Increasing
		v.stateB = 0
	} else if v.stateG == Increasing && v.g+v.step >= MaxBrightness {
		v.stateR = 0
		v.stateG = Decreasing
		v.stateB = Increasing
	} else if v.stateB == Increasing && v.b+v.step >= MaxBrightness {
		v.stateR = Increasing
		v.stateG = 0
		v.stateB = Decreasing
	}

	// Change values of colors
	if v.stateR == Increasing {
		v.r += v.step
	} else if v.stateR == Decreasing {
		v.r -= v.step
	}
	if v.stateG == Increasing {
		v.g += v.step
	} else if v.stateG == Decreasing {
		v.g -= v.step
	}
	if v.stateB == Increasing {
		v.b += v.step
	} else if v.stateB == Decreasing {
		v.b -= v.step
	}

	n.Strips[v.curStrip].Pixels[v.curLED].R = v.r
	n.Strips[v.curStrip].Pixels[v.curLED].G = v.g
	n.Strips[v.curStrip].Pixels[v.curLED].B = v.b
}

func (v SpiralViz) RefreshRate() float64 {
	return .05
}
