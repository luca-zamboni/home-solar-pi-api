package device

import (
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/datatypes"
)

var (
	errVarNotInitialized    = errors.New("variables not initialized")
	errMethodNotImplemented = errors.New("method not implemented")
)

type DeviceConfig struct {
	Host string
	Port int
	Api  string
}

type Device struct {
	// config *DeviceConfig
	ID            uint `gorm:"unique;primaryKey;autoIncrement"`
	Name          string
	Type          DeviceType
	CurrentStatus DeviceStatus
	Actions       []DeviceActionLog
	Info          datatypes.JSON `sql:"type:jsonb"`
}

type DeviceActionLog struct {
	DeviceId     uint         `gorm:"primaryKey;autoIncrement:false"`
	StatusUpdate DeviceAction `gorm:"primaryKey;autoIncrement:false"`
}

type DeviceStatus string

const (
	ACTIVE   DeviceStatus = "ACTIVE"
	INACTIVE DeviceStatus = "INACTIVE"
	ON       DeviceStatus = "ON"
	OFF      DeviceStatus = "OFF"
)

type DeviceAction string

const (
	POWER_ON  DeviceAction = "POWER_ON"
	POWER_OFF DeviceAction = "POWER_OFF"
)

type DeviceImpl interface {
	GetDeviceUrl() (string, error)
	PowerOn() error
	PowerOff() error
	Status() (DeviceStatus, error)
	ReadValue() (any, error)
	GetConfig() (DeviceConfig, error)
}

func (config DeviceConfig) GetUrl() (string, error) {
	if config.Host == "" || config.Port == 0 || config.Api == "" {
		return "", errVarNotInitialized
	}

	return fmt.Sprintf("http://%s:%d/%s", config.Host, config.Port, config.Api), nil
}

func (device Device) GetConfig() (DeviceConfig, error) {
	var config DeviceConfig
	err := json.Unmarshal(device.Info, &config)
	return config, err
}

func (device Device) GetDeviceUrl() (string, error) {
	config, _ := device.GetConfig()
	return config.GetUrl()
}

func (device Device) PowerOn() error {
	return errMethodNotImplemented
}

func (device Device) PowerOff() error {
	return errMethodNotImplemented
}

func (device Device) Status() (DeviceStatus, error) {
	return INACTIVE, nil
}

func (device Device) ReadValue() (any, error) {
	return false, errMethodNotImplemented
}
