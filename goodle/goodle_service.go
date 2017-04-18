package goodle

import (
	"fmt"
	"github.com/currantlabs/gatt"
	"time"
)

const GOODLE_SERVICE_UUID string = "6dfdee53-2fb8-4842-9981-9078176f5649"

const GOODLE_COLOR_CHARACTERISTIC string = "7394c995-0ae1-404a-9192-b8a50dfc288a"
const GOODLE_VIZ_CHARACTERISTIC string = "f546514f-758b-4466-88f7-9e1291a0672e"
const GOODLE_REFRESH_RATE_CHARACTERISTIC string = "0b59be6a-90f1-41a5-baea-49cc2df5a30c"

func NewGoodleService(h gatt.WriteHandlerFunc) *gatt.Service {
	n := 0
	s := gatt.NewService(gatt.MustParseUUID(GOODLE_SERVICE_UUID))
	s.AddCharacteristic(gatt.MustParseUUID(GOODLE_COLOR_CHARACTERISTIC)).HandleReadFunc(
		func(rsp gatt.ResponseWriter, req *gatt.ReadRequest) {
			fmt.Fprintf(rsp, "count: %d", n)
			n++
		})

	s.AddCharacteristic(gatt.MustParseUUID(GOODLE_VIZ_CHARACTERISTIC)).HandleWriteFunc(h)

	s.AddCharacteristic(gatt.MustParseUUID(GOODLE_REFRESH_RATE_CHARACTERISTIC)).HandleNotifyFunc(
		func(r gatt.Request, n gatt.Notifier) {
			cnt := 0
			for !n.Done() {
				fmt.Fprintf(n, "Count: %d", cnt)
				cnt++
				time.Sleep(time.Second)
			}
		})

	return s
}
