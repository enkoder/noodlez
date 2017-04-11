package noodle

import (
	"fmt"
)

const (
	MaxBrightness = 200
	Increasing    = 255
	Decreasing    = 1
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
