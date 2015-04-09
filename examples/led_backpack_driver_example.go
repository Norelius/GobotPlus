/*
Portions adapted from George McBay.
https://bitbucket.org/corburn/i2c/src/1235f1776ee7/HT16K33/
*/

package main

import (
	"time"

	"github.com/Norelius/GobotPlus/drivers"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
)

var byteCommand = map[string][]byte{
	"full": []byte{
		byte(0xFF), // Row 0
		byte(0xFF), // Row 1
		byte(0xFF), // Row 2
		byte(0xFF), // Row 3
		byte(0xFF), // Row 4
		byte(0xFF), // Row 5
		byte(0xFF), // Row 6
		//byte(0xFF), // Row 7
	},
	"blank": []byte{
		byte(0x00), // Row 0
		byte(0x00), // Row 1
		byte(0x00), // Row 2
		byte(0x00), // Row 3
		byte(0x00), // Row 4
		byte(0x00), // Row 5
		byte(0x00), // Row 6
		//byte(0x00), // Row 7
	},
	"row": []byte{
		byte(0x01), // Row 0
		byte(0x01), // Row 1
		byte(0x01), // Row 2
		byte(0x01), // Row 3
		byte(0x01), // Row 4
		byte(0x01), // Row 5
		byte(0x01), // Row 6
		//byte(0x01), // Row 7
	},
	"smiley": []byte{
		0x1E, // Row 0
		0x42, // Row 1
		0xA5, // Row 2
		0x81, // Row 3
		0xA5, // Row 4
		0x99, // Row 5
		0x42, // Row 6
		//0x1E, // Row 7
	},
	"antiborksmiley": []byte{
		0x1E, // Row 0
		0x21, // Row 1
		0xD2, // Row 2
		0xC0, // Row 3
		0xD2, // Row 4
		0xCC, // Row 5
		0x21, // Row 6
		//0x1E, // Row 7
	},
}

var lol = 0

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/tty.usbmodem1431")
	backPack := drivers.NewBackPackDriver(firmataAdaptor, byte(0x71), "backpack")

	work := func() {
		gobot.Every(3*time.Second, func() {
			lol++
			display := "blank"
			if lol%3 == 0 {
				display = "smiley"
			}
			if lol%3 == 1 {
				display = "antiborksmiley"
			}
			b := byteCommand[display]
			backPack.Write(b)
		})
	}

	robot := gobot.NewRobot("ledMatrixBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{backPack},
		work,
	)

	gbot.AddRobot(robot)
	gbot.Start()
}
