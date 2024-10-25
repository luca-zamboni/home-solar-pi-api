package api

import (
	"fmt"
	"home-solar-pi/pkg/device"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiService struct {
	deviceManager *device.DeviceManager
}

func NewApiServer(deviceManager *device.DeviceManager) ApiService {
	return ApiService{
		deviceManager: deviceManager,
	}
}

func (s *ApiService) StartServer() {
	r := gin.Default()
	apiRouter := r.Group("/api")

	deviceGroup := apiRouter.Group("/device")

	deviceGroup.GET("/all", s.getAllDevices)
	deviceGroup.GET("/:device/value", s.getDeviceValue)
	deviceGroup.PUT("/:device/on", s.setDeviceOn)
	deviceGroup.PUT("/:device/off", s.setDeviceOff)
	deviceGroup.GET("/:device/status", s.getDeviceStatus)

	err := r.Run(":5000")
	panic(err.Error())
}

func (api *ApiService) getAllDevices(c *gin.Context) {
	devices, err := api.deviceManager.GetAllDevices()

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, devices)
}

func (api *ApiService) getDeviceValue(c *gin.Context) {

	device, err := api.deviceManager.GetDeviceImpl(device.DeviceType(c.Param("device")))

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	power, err := device.ReadValue()

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if power == nil {
		c.String(http.StatusInternalServerError, "Power nil")
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("%d", power.(int)))
}

func (api *ApiService) setDeviceOn(c *gin.Context) {

	device, err := api.deviceManager.GetDeviceImpl(device.DeviceType(c.Param("device")))

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = device.PowerOn()

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "ON")
}

func (api *ApiService) setDeviceOff(c *gin.Context) {

	device, err := api.deviceManager.GetDeviceImpl(device.DeviceType(c.Param("device")))

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	device.PowerOff()

	c.String(http.StatusOK, "OFF")
}

func (api *ApiService) getDeviceStatus(c *gin.Context) {

	device, err := api.deviceManager.GetDeviceImpl(device.DeviceType(c.Param("device")))

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
