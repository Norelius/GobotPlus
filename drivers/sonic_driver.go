package drivers

import (
	"errors"
	"math"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

const (
	METERS_PER_SECOND = 340.29
	Distance          = "distance"
	OutOfRange        = "outofrange"
	Error             = "error"
)

var _ gobot.Driver = (*SonicDriver)(nil)

// LedBackpackDriver respresents a sonic distance sensor.
type SonicDriver struct {
	pinger     SonicPinger
	smoothing  int
	pin        string
	name       string
	halt       chan bool
	interval   time.Duration
	connection gpio.DigitalReader
	gobot.Eventer
}

type SonicPinger struct {
	pin        string
	name       string
	halt       chan bool
	interval   time.Duration
	connection gpio.DigitalWriter
	gobot.Eventer
}

// NewSonicDriver return a new SonicDriver given connection, name, smoothing
// value, trigger pin to start the measuring with and echo pin to read the distance from.
//
func NewSonicDriver(a gobot.Connection, name string, smoothing int, trigger_pin string, echo_pin string) (*SonicDriver, error) {
	reader, err := a.(gpio.DigitalReader)
	if !err {
		return nil, errors.New("DigitalRead is not supported by this platform")
	}

	writer, err := a.(gpio.DigitalWriter)
	if !err {
		return nil, errors.New("DigitalWrite is not supported by this platform")
	}

	b := &SonicDriver{
		name:       name,
		connection: reader,
		smoothing:  smoothing,
		pinger: SonicPinger{
			name:       name + "_pinger",
			connection: writer,
			pin:        trigger_pin,
			Eventer:    gobot.NewEventer(),
			interval:   10 * time.Microsecond,
			halt:       make(chan bool),
		},
		pin:      echo_pin,
		Eventer:  gobot.NewEventer(),
		interval: 10 * time.Microsecond,
		halt:     make(chan bool),
	}

	b.AddEvent(Distance)
	b.AddEvent(OutOfRange)
	b.AddEvent(Error)

	return b, nil
}

func (a *SonicDriver) Pinger() *SonicPinger {
	return &a.pinger
}

func (a *SonicDriver) Start() (errs []error) {
	high := time.Now()
	distances := make([]int64, a.smoothing)
	pointer := 0
	active := false
	state := 0
	go func() {
		for {

			newValue, err := a.connection.DigitalRead(a.Pin())
			if err != nil {
				gobot.Publish(a.Event(Error), err)
			} else if newValue != state && newValue != -1 {
				state = newValue
				if state == 1 {
					high = time.Now()
				} else {
					duration := time.Now().Sub(high).Nanoseconds()
					if duration > 27000000 {
						gobot.Publish(a.Event(OutOfRange), nil)
					} else {
						distance := duration / 58000
						distances[pointer] = distance
						active = active || (pointer == 4)
						pointer = (pointer + 1) % a.smoothing
						if active {
							var avg float64 = 0
							for i := 0; i < len(distances); i++ {
								avg += float64(distances[i] * distances[i])
							}
							avg = math.Sqrt(avg / float64(a.smoothing))
							gobot.Publish(a.Event(Distance), int64(avg))
						}
					}
				}
			}

			select {
			case <-time.After(60 * time.Microsecond):
			case <-a.halt:
				return
			}
		}
	}()
	return
}

func (a *SonicPinger) Start() (errs []error) {
	go func() {
		for {
			a.connection.DigitalWrite(a.Pin(), byte(0))
			time.Sleep(50 * time.Microsecond)
			a.connection.DigitalWrite(a.Pin(), byte(1))
			time.Sleep(100 * time.Microsecond)
			a.connection.DigitalWrite(a.Pin(), byte(0))
			select {
			case <-time.After(60 * time.Millisecond):
			case <-a.halt:
				return
			}
		}
	}()
	return
}

func (a *SonicDriver) Halt() (errs []error) {
	a.halt <- true
	return
}

func (a *SonicPinger) Halt() (errs []error) {
	a.halt <- true
	return
}

func (a *SonicDriver) Pin() string { return a.pin }

func (a *SonicPinger) Pin() string { return a.pin }

func (a *SonicDriver) Name() string { return a.name }

func (a *SonicPinger) Name() string { return a.name }

func (a *SonicDriver) Connection() gobot.Connection { return a.connection.(gobot.Connection) }

func (a *SonicPinger) Connection() gobot.Connection { return a.connection.(gobot.Connection) }
