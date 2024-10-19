package worker

import (
	"home-solar-pi/pkg/db"
	"home-solar-pi/pkg/device"
	"home-solar-pi/pkg/utils"
	"log"
	"math/rand"
	"time"
)

type HeaterInverterWorker struct {
	inverterService *device.InterverService
	heaterService   *device.HeaterService
	logger          *log.Logger
	dbService       *db.DbService
	threshold       int
}

type FallBackIntervall int

var (
	INCREASE FallBackIntervall = 1
	NORMAL   FallBackIntervall = 0
)

func NewHeaterInverterWorker(inverterService *device.InterverService, heaterService *device.HeaterService,
	logger *log.Logger, dbService *db.DbService, threshold int) HeaterInverterWorker {
	return HeaterInverterWorker{
		inverterService: inverterService,
		heaterService:   heaterService,
		logger:          logger,
		dbService:       dbService,
		threshold:       threshold,
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

	statusOn, err := w.heaterService.GetStatus()

	if err != nil {
		w.logger.Printf("Error heater := %s\n", err.Error())
		return INCREASE
	}

	// shelly auto deactivates after HEATER_TOGGLE seconds
	if statusOn {
		return NORMAL
	}

	reading, err := w.getInverterReading()
	if err != nil {
		w.logger.Printf("Inverter error := %s\n", err.Error())
		return INCREASE
	}

	w.logger.Printf("Inverter Power %+v\n", reading)

	// Inserting in log table to do statistics
	w.dbService.InsertReading(reading)

	if reading > w.threshold {
		_, err := w.heaterService.PowerOn()
		if err != nil {
			w.logger.Printf("Error heater := %s\n", err.Error())
			return NORMAL
		}
		w.logger.Println("Heater activated")

		// Inserting in Heater logs when activated with which reading from inverter
		w.dbService.InsertHeaterlog(reading, true)
	}

	return NORMAL

}

func (w *HeaterInverterWorker) getInverterReading() (int, error) {
	var power *device.InverterResponse

	if utils.Inactive {
		return rand.Int()%300 + 400, nil
	}

	var err error
	power, err = w.inverterService.GetCurrentPower()

	if err != nil {
		w.logger.Printf("Error retrieving current power\n")
		return 0, err
	}

	return power.Body.Data.PAC.Values["1"], nil

}
