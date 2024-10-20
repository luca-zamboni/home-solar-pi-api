package api

import (
	"home-solar-pi/pkg/device"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApiService struct {
	InverterService *device.InterverService
}

func (s *ApiService) StartServer() {
	r := gin.Default()
	apiRouter := r.Group("/api")
	apiRouter.GET("/inverter", s.getCurrentPower)

	err := r.Run(":5000")
	panic(err.Error())
}

func (api *ApiService) getCurrentPower(c *gin.Context) {

	power, err := api.InverterService.GetCurrentPower()

	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading power")
	}

	if power == nil {
		c.String(http.StatusInternalServerError, "Error reading power")
	}

	v := strconv.Itoa(power.Body.Data.PAC.Values["1"])
	c.String(http.StatusOK, v)
}
