package main

import (
	"os"

	"github.com/enkoder/noodlez"
)

const (
	ButtonPin = "gpio17"
)

func main() {
	n, err := noodlez.NewNoodle(ButtonPin)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	err = n.Off()
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	n.VizLoop()
}
