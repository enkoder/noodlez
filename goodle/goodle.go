package goodle

import (
	"fmt"
	"github.com/currantlabs/gatt"
	"github.com/currantlabs/gatt/examples/option"
)

func InitBluetooth(h gatt.WriteHandlerFunc) error {
	d, err := gatt.NewDevice(option.DefaultServerOptions...)
	if err != nil {
		return fmt.Errorf("Failed to open device, err: %s", err)
	}

	// Register optional handlers.
	d.Handle(
		gatt.CentralConnected(func(c gatt.Central) { fmt.Println("Connect: ", c.ID()) }),
		gatt.CentralDisconnected(func(c gatt.Central) { fmt.Println("Disconnect: ", c.ID()) }),
	)

	// A mandatory handler for monitoring device state.
	onStateChanged := func(d gatt.Device, s gatt.State) {
		fmt.Printf("State: %s\n", s)
		switch s {
		case gatt.StatePoweredOn:
			// Setup GAP and GATT services for Linux implementation.
			// OS X doesn't export the access of these services.
			d.AddService(NewGapService("Goodle")) // no effect on OS X
			d.AddService(NewGattService())        // no effect on OS X

			goodleService := NewGoodleService(h)
			// A simple count service for demo.
			d.AddService(goodleService)

			// Advertise device name and service's UUIDs.
			fmt.Printf("Advertising: %s\n", 'x')
			d.AdvertiseNameAndServices("Goodle", []gatt.UUID{goodleService.UUID()})
		default:
		}
	}

	return d.Init(onStateChanged)
}
