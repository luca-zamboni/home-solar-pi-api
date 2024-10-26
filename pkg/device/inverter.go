package device

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

// Device
type InverterDevice struct {
	Device
}

func NewInterver(baseDevice Device) InverterDevice {
	return InverterDevice{
		Device: baseDevice,
	}
}

func (s InverterDevice) ReadValue() (any, error) {

	if s.State == "Debug" {
		return rand.Int31()%200 + 500, nil
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
