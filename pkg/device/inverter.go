package device

import (
	"encoding/json"
	"home-solar-pi/pkg/utils"
	"net/http"
)

// Model
type InverterResponse struct {
	Head any
	Body struct {
		Data struct {
			PAC          Power
			DAY_ENERGY   Power
			YEAR_ENERGY  Power
			TOTAL_ENERGY Power
		}
	}
}

type Power struct {
	Unit   string
	Values map[string]int
}

// Device
type InverterDevice struct {
	BaseDevice
}

func NewInterver(config DeviceConfig) InverterDevice {
	return InverterDevice{
		BaseDevice: BaseDevice{
			config: &config,
		},
	}
}

func (s InverterDevice) ReadValue() (any, error) {

	if utils.Debug {
		return 27, nil
	}

	uri, err := s.GetDeviceUrl()
	if err != nil {
		return -1, err
	}

	resp, err := http.Get(uri)
	if err != nil {
		return -1, err
	}

	var power InverterResponse
	err = json.NewDecoder(resp.Body).Decode(&power)

	if err != nil {
		return -1, err
	}

	return power.Body.Data.PAC.Values["1"], err
}
