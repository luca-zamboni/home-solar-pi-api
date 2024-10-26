package rule

import (
	"home-solar-pi/pkg/device"
	"log"
	"time"
)

type HeaterInverterWorker struct {
	inverterDevice *device.InverterDevice
	heaterDevice   *device.HeaterDevice
	logger         *log.Logger
	dbService      *device.DbService
	threshold      int
}

type FallBackIntervall int

var (
	INCREASE FallBackIntervall = 1
	NORMAL   FallBackIntervall = 0
)

func NewHeaterInverterWorker(inverterDevice *device.InverterDevice, heaterDevice *device.HeaterDevice,
	logger *log.Logger, dbService *device.DbService, threshold int) HeaterInverterWorker {
	return HeaterInverterWorker{
		inverterDevice: inverterDevice,
		heaterDevice:   heaterDevice,
		logger:         logger,
		dbService:      dbService,
		threshold:      threshold,
	}
}

// starts the heater inverte cycle
// Intervall is the normal interval of updates.
// if any error occurs the interval is doubled until interval * 16
// if no error the interval is set to normal
// Steps
// 1. Checks the heater status. if already on -> skip cycle
// 2. Retrieving Inverter power produced
// 3. if reading > threashold
// then 3.1 start the heater
func (w *HeaterInverterWorker) StartHeaterInverterCycle(interval time.Duration) {

	intervalSleep := interval
	for {

		fallbackIntervall := w.doWork()

		if fallbackIntervall == NORMAL {
			intervalSleep = interval
		} else {
			intervalSleep = min(interval*16, intervalSleep*2)
		}

		time.Sleep(intervalSleep)
	}

}

func (w *HeaterInverterWorker) doWork() FallBackIntervall {

	status, err := w.heaterDevice.Status()

	if err != nil {
		w.logger.Printf("Error heater := %s\n", err.Error())
		return INCREASE
	}

	// shelly auto deactivates after HEATER_TOGGLE seconds
	if status == device.INACTIVE || status == device.ON {
		return NORMAL
	}

	reading, err := w.getInverterReading()
	if err != nil {
		w.logger.Printf("Inverter error := %s\n", err.Error())
		return INCREASE
	}

	w.logger.Printf("Inverter Power %+v\n", reading)

	// Inserting in log table to do statistics
	// w.dbService.InsertReading(reading)

	if reading > w.threshold {
		err := w.heaterDevice.PowerOn()
		if err != nil {
			w.logger.Printf("Error heater := %s\n", err.Error())
			return NORMAL
		}
		w.logger.Println("Heater activated")

		// Inserting in Heater logs when activated with which reading from inverter
		// w.dbService.InsertHeaterAction(reading, true)
	}

	return NORMAL

}

func (w *HeaterInverterWorker) getInverterReading() (int, error) {

	var err error
	power, err := w.inverterDevice.ReadValue()

	if err != nil {
		w.logger.Printf("Error retrieving current power\n")
		return 0, err
	}

	return power.(int), nil

}
