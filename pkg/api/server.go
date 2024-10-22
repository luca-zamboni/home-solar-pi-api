package api

import (
	"fmt"
	"home-solar-pi/pkg/device"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceNotFound struct {
	id string
}

func (d *DeviceNotFound) Error() string { return fmt.Sprintf("Device not found %s", d.id) }

func NewApiServer(devices map[string]device.Device) ApiService {
	return ApiService{
		devices: devices,
	}
}

type ApiService struct {
	devices map[string]device.Device
}

func (s *ApiService) StartServer() {
	r := gin.Default()
	apiRouter := r.Group("/api")

	deviceGroup := apiRouter.Group("/device")

	deviceGroup.GET("/:device/value", s.getDeviceValue)
	deviceGroup.PUT("/:device/on", s.setDeviceOn)
	deviceGroup.PUT("/:device/off", s.setDeviceOff)
	deviceGroup.GET("/:device/status", s.getDeviceStatus)

	err := r.Run(":5000")
	panic(err.Error())
}

func (s *ApiService) getDeviceById(id string) (device.Device, error) {
	device, ok := s.devices[id]

	if ok {
		return device, nil
	}

	return nil, &DeviceNotFound{id: id}

}

func (api *ApiService) getDeviceValue(c *gin.Context) {

	device, err := api.getDeviceById(c.Param("device"))

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	power, err := (device).ReadValue()

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if power == nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("%d", power.(int)))
}

func (api *ApiService) setDeviceOn(c *gin.Context) {

	device, err := api.getDeviceById(c.Param("device"))

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	device.PowerOn()

	c.String(http.StatusOK, "ON")
}

func (api *ApiService) setDeviceOff(c *gin.Context) {

	device, err := api.getDeviceById(c.Param("device"))

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	device.PowerOff()

	c.String(http.StatusOK, "OFF")
}

func (api *ApiService) getDeviceStatus(c *gin.Context) {

	device, err := api.getDeviceById(c.Param("device"))

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	status, err := device.Status()

	if err != nil {
		c.String(http.StatusBadRequest, "asd"+err.Error())
		return
	}

	c.String(http.StatusOK, string(status))
}
