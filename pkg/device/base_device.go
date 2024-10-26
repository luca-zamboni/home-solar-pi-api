package device

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

var (
	errVarNotInitialized    = errors.New("variables not initialized")
	errMethodNotImplemented = errors.New("method not implemented")
	errNotFound             = errors.New("not found")
)

type DriverType string

const (
	InverterType DriverType = "Inverter"
	HeaterType   DriverType = "Heater"
)

type DeviceConfig struct {
	Host string
	Port int
	Api  string
}

type Device struct {
	Name   string     `yaml:"Name"`
	Driver DriverType `yaml:"Driver"`
	State  string     `yaml:"State"`
	Info   any        `yaml:"Info"`
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

type DeviceDriver interface {
	GetDeviceName() string
	GetDriverName() DriverType
	GetDeviceUrl() (string, error)
	PowerOn() error
	PowerOff() error
	Status() (DeviceStatus, error)
	ReadValue() (any, error)
	GetConfig() (DeviceConfig, error)
}

func (device Device) GetDeviceName() string {
	return device.Name
}

func (device Device) GetDriverName() DriverType {
	return device.Driver
}

func (config DeviceConfig) GetUrl() (string, error) {
	if config.Host == "" || config.Port == 0 || config.Api == "" {
		return "", errVarNotInitialized
	}

	return fmt.Sprintf("%s:%d/%s", config.Host, config.Port, config.Api), nil
}

func (device Device) GetConfig() (DeviceConfig, error) {
	var config DeviceConfig
	err := mapstructure.Decode(device.Info, &config)
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
