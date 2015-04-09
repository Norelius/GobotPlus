package main

import (
	"fmt"

	"github.com/Norelius/GobotPlus/drivers"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

const (
	METERS_PER_SECOND      = 340.29
	METERS_PER_MILLISECOND = METERS_PER_SECOND / 1000
	METERS_PER_MICROSEOND  = METERS_PER_SECOND / 1000000
	METERS_PER_NANOSECOND  = METERS_PER_SECOND / 1000000000
	CM_PER_NANOSECOND      = METERS_PER_NANOSECOND / 100
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/tty.usbmodem1411")

	sonicDriver, err := drivers.NewSonicDriver(firmataAdaptor, "distance", 8, "6", "2")
	if err != nil {
		fmt.Println(err)
		return
	}
	sonicPinger := sonicDriver.Pinger()

	red := gpio.NewLedDriver(firmataAdaptor, "led_red", "8")
	green := gpio.NewLedDriver(firmataAdaptor, "led_green", "5")

	work := func() {

		gobot.On(sonicDriver.Event("distance"), func(data interface{}) {
			distance := data.(int64)
			fmt.Printf("%v cm\n", data.(int64))
			if distance > 100 && distance < 200 {
				red.Off()
				green.On()
			} else {
				red.On()
				green.Off()
			}
		})

		gobot.On(sonicDriver.Event("outofrange"), func(data interface{}) {
			red.On()
			green.Off()
		})
		gobot.On(sonicDriver.Event("error"), func(data interface{}) {
			fmt.Println("ERROR!!!")
		})
	}

	robot := gobot.NewRobot("sonicBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{red, green, sonicDriver, sonicPinger},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
