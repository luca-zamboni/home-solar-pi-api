package device

import (
	"log"
)

type DeviceManager struct {
	dbService *DbService
	logger    *log.Logger
}

type DeviceType string

const (
	InverterType DeviceType = "Inverter"
	HeaterType   DeviceType = "Heater"
)

func NewDeviceManager(dbService *DbService, logger *log.Logger) DeviceManager {
	return DeviceManager{
		dbService: dbService,
		logger:    logger,
	}
}

func (m DeviceManager) GetAllDevices() ([]Device, error) {
	return m.dbService.GetAllDevices()
}

func (m DeviceManager) GetDeviceImpl(id DeviceType) (DeviceImpl, error) {
	device, err := m.dbService.GetDeviceByType(id)

	if err != nil {
		return nil, err
	}

	if device.Type == HeaterType {
		return NewHeater(device), nil
	}
	if device.Type == InverterType {
		return NewInterver(device), nil
	}

	return nil, nil
}
