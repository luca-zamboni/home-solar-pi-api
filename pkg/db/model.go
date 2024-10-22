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

type Action struct {
	ID        uint `gorm:"unique;primaryKey;autoIncrement"`
	CreatedAt time.Time
	Name      string
}

type HeaterAction struct {
	Reading   int
	CreatedAt time.Time
	Status    HeaterStatus
	ActionID  uint
	Action    Action
}
