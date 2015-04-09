package drivers

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
)

const (
	REGISTER_DISPLAY_SETUP = 0x80
	REGISTER_SYSTEM_SETUP  = 0x20
	REGISTER_DIMMING       = 0xE0
	BLINKRATE_OFF          = 0x00
	BLINKRATE_2HZ          = 0x01
	BLINKRATE_1HZ          = 0x02
	BLINKRATE_HALFHZ       = 0x03
)

var _ gobot.Driver = (*LedBackpackDriver)(nil)

// LedBackpackDriver respresents a digital 8x8 led matrix.
type LedBackpackDriver struct {
	name        string
	i2c_address byte
	connection  i2c.I2c
	gobot.Commander
}

// NewLedBackpackDriver return a new LedBackpackDriver given a I2c connection,
// address and name.
//
// Adds the following API Commands:
//	"Brightness" - See LedBackpackDriver.Brightness
//	"SetBlinkFrequency" - See LedBackpackDriver.SetBlinkFrequency
//	"Write" - See LedBackpackDriver.Write
func NewLedBackpackDriver(a i2c.I2c, i2c_address byte, name string) *LedBackpackDriver {
	bp := &LedBackpackDriver{
		name:        name,
		i2c_address: i2c_address,
		connection:  a,
		Commander:   gobot.NewCommander(),
	}

	bp.AddCommand("Write", func(params map[string]interface{}) interface{} {
		byteCommands := params["byteCommands"].([]byte)
		return bp.Write(byteCommands)
	})

	return bp
}

// Name returns the LedBackpackDrivers name
func (bp *LedBackpackDriver) Name() string { return bp.name }

// Pin returns the LedBackpackDrivers pin
func (bp *LedBackpackDriver) Connection() gobot.Connection {
	return bp.connection.(gobot.Connection)
}

// Halt returns true if device is halted successfully
func (bp *LedBackpackDriver) Halt() (errs []error) { return }

// Start writes start bytes
func (bp *LedBackpackDriver) Start() (errs []error) {
	if err := bp.connection.I2cStart(bp.i2c_address); err != nil {
		return []error{err}
	}
	if err := bp.connection.I2cWrite([]byte{REGISTER_SYSTEM_SETUP | 0x01, 0x00}); err != nil {
		return []error{err}
	}
	if err := bp.connection.I2cWrite([]byte{REGISTER_DISPLAY_SETUP | 0x01 |
		(BLINKRATE_OFF << 1), 0x00}); err != nil {
		return []error{err}
	}
	return
}

// Brightness sets the ledbackpack to the specified level of brightness
func (bp *LedBackpackDriver) Brightness(brightness byte) (err error) {
	if brightness < 0 || brightness > 15 {
		return fmt.Errorf("Brightness not within acceptable range.")
	}
	if err := bp.connection.I2cWrite([]byte{REGISTER_DIMMING | brightness, 0x00}); err != nil {
		return err
	}
	return
}

// SetBlinkFrequency for the LedBackpack.
func (bp *LedBackpackDriver) SetBlinkFrequency(freq byte) (err error) {
	if err = bp.connection.I2cWrite([]byte{REGISTER_DISPLAY_SETUP | 0x01 |
		(freq << 1), 0x00}); err != nil {
		return
	}
	return
}

// Write takes a bytearray and writes it to the leds. Each byte corresponds to
// a row, each bit to a pixel. 1 sets the led to on and 0 to off.
//
// Currently only gives control of the first 7 rows.
func (bp *LedBackpackDriver) Write(byteCommand []byte) (err error) {
	fmt.Printf("%v\n", byteCommand)
	/*if len(byteCommand) != 7 {
		return err
	}*/
	// Interjecting 0's because the drivers also support RGB colours, and this
	// was developed for a monochrom display.
	toSend := make([]byte, len(byteCommand)*2)
	for index, data := range byteCommand {
		toSend[index*2+1] = data
	}
	if err = bp.connection.I2cWrite(toSend); err != nil {
		return
	}
	return
}
