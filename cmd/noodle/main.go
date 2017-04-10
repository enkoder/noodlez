package main

import (
	"fmt"
	"os"
	"time"

	"github.com/enkoder/noodle"
)

const (
	ButtonPin = "gpio17"
)

func main() {
	n, err := noodle.NewNoodle(ButtonPin)
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

	for {
		time.Sleep(time.Millisecond * 500)

		value, err := n.ButtonPressed()
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Error during DigitalRead: %v\n", err))
			os.Exit(1)
		}

		if value {
			n.Blue()
		} else {
			n.Red()
		}

		fmt.Printf("Button Value: %t\n", value)
	}
}
