package db

import "time"

type InverterReading struct {
	Reading   int
	CreatedAt time.Time
}

type HeaterStatus string

var (
	POWER_ON  HeaterStatus = "POWER_ON"
	POWER_OFF HeaterStatus = "POWER_OFF"
)

type HeaterLogs struct {
	Reading   int
	CreatedAt time.Time
	Status    HeaterStatus
}
