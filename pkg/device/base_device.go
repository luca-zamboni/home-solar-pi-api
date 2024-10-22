package device

import (
	"errors"
	"fmt"
)

var (
	errVarNotInitialized    = errors.New("variables not initialized")
	errMethodNotImplemented = errors.New("method not implemented")
)

type DeviceConfig struct {
	Host    string
	Port    int
	ApiPath string
}

type BaseDevice struct {
	config *DeviceConfig
}

type DeviceStatus string

const (
	POWER_ON  DeviceStatus = "POWER_ON"
	POWER_OFF DeviceStatus = "POWER_OFF"
	INACTIVE  DeviceStatus = "INACTIVE"
)

type Device interface {
	GetDeviceUrl() (string, error)
	PowerOn() error
	PowerOff() error
	Status() (DeviceStatus, error)
	ReadValue() (any, error)
}

func (config DeviceConfig) GetUrl() (string, error) {
	if config.Host == "" || config.Port == 0 || config.ApiPath == "" {
		return "", errVarNotInitialized
	}

	return fmt.Sprintf("http://%s:%d/%s", config.Host, config.Port, config.ApiPath), nil
}

func (device BaseDevice) GetDeviceUrl() (string, error) {
	return device.config.GetUrl()
}

func (device BaseDevice) PowerOn() error {
	return errMethodNotImplemented
}

func (device BaseDevice) PowerOff() error {
	return errMethodNotImplemented
}

func (device BaseDevice) Status() (DeviceStatus, error) {
	return INACTIVE, nil
}

func (device BaseDevice) ReadValue() (any, error) {
	return false, errMethodNotImplemented
}
