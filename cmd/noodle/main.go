package main

import (
	"flag"
	"os"

	"github.com/enkoder/noodlez"
)

const (
	ButtonPin = "gpio17"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug output")
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

	n.VizLoop(*debug)
}
