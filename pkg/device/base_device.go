package device

import (
	"errors"
	"fmt"
	"log"
)

type BaseDeviceService struct {
	HOST   string
	PORT   int
	API    string
	logger *log.Logger
}

var (
	errVarNotInitialized = errors.New("variables not initialized")
)

type DeviceService interface {
	GetDeviceUrl()
}

func (ds *BaseDeviceService) GetDeviceUrl() (string, error) {
	if ds.HOST == "" || ds.PORT == 0 || ds.API == "" {
		return "", errVarNotInitialized
	}

	return fmt.Sprintf("http://%s:%d/%s", ds.HOST, ds.PORT, ds.API), nil
}
