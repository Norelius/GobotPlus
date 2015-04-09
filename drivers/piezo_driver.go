package drivers

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

type PiezoDriver struct {
	pin        string
	name       string
	connection gpio.DigitalWriter
	gobot.Commander
}

func NewPiezoDriver(a gpio.DigitalWriter, name string, pin string) *PiezoDriver {
	p := &PiezoDriver{
		name:       name,
		pin:        pin,
		connection: a,
		Commander:  gobot.NewCommander(),
	}

	p.AddCommand("Tone", func(params map[string]interface{}) interface{} {
		level := byte(params["level"].(float64))
		return p.Tone(level)
	})

	p.AddCommand("NoTone", func(params map[string]interface{}) interface{} {
		return p.NoTone()
	})

	p.AddCommand("ToneDuration", func(params map[string]interface{}) interface{} {
		level := byte(params["level"].(float64))
		duration_ms := int(params["duration_ms"].(float64))
		return p.ToneDuration(level, duration_ms)
	})

	return p
}

// Start implements the Driver interface.
func (p *PiezoDriver) Start() (errs []error) { return }

// Halt implements the Driver interface.
func (p *PiezoDriver) Halt() (errs []error) { return }

// Name returns the PiezoDrivers name
func (p *PiezoDriver) Name() string { return p.name }

// Pin returns the PiezoDrivers name
func (p *PiezoDriver) Pin() string { return p.pin }

// Connection returns the PiezoDrivers Connection
func (p *PiezoDriver) Connection() gobot.Connection {
	return p.connection.(gobot.Connection)
}

func (p *PiezoDriver) Tone(level byte) (err error) {
	if writer, ok := p.connection.(gpio.ServoWriter); ok {
		return writer.ServoWrite(p.Pin(), level)
	}
	return gpio.ErrPwmWriteUnsupported
}

func (p *PiezoDriver) NoTone() (err error) {
	return p.connection.DigitalWrite(p.Pin(), 0)
}

func (p *PiezoDriver) ToneDuration(level byte, duration_ms int) (err error) {
	if writer, ok := p.connection.(gpio.ServoWriter); ok {
		if err := writer.ServoWrite(p.Pin(), level); err != nil {
			return err
		}
		time.Sleep(time.Duration(duration_ms) * time.Millisecond)
		return p.connection.DigitalWrite(p.Pin(), 0)
	}
	return gpio.ErrPwmWriteUnsupported
}
