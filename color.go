package noodle

import (
	"math/rand"
)

type Pixel struct {
	R uint8
	G uint8
	B uint8
}

func (p Pixel) Equals(other Pixel) bool {
	return p.R == other.R && p.G == other.G && p.B == other.B
}

var (
	Off         = Pixel{R: 0, G: 0, B: 0}
	Blue        = Pixel{R: 0, G: 0, B: MaxBrightness}
	Red         = Pixel{R: MaxBrightness, G: 0, B: 0}
	Cyan        = Pixel{R: 0, G: MaxBrightness, B: MaxBrightness}
	Green       = Pixel{R: 0, G: MaxBrightness, B: 0}
	Magenta     = Pixel{R: MaxBrightness, G: 0, B: MaxBrightness}
	Pink        = Pixel{R: 138, G: 0, B: 211}
	Yellow      = Pixel{R: MaxBrightness, G: MaxBrightness, B: 0}
	Orange      = Pixel{R: MaxBrightness, G: 130, B: 0}
	NamedColors = []Pixel{Blue, Red, Cyan, Green, Magenta, Pink, Yellow, Orange}
)

func RandNamedColor() Pixel {
	return NamedColors[rand.Intn(len(NamedColors))]
}
