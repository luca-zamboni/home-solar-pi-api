package device

import (
	"encoding/json"
	"net/http"
)

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

type InterverService struct {
	BaseDeviceService
}

func NewInterverService(HOST string, PORT int, API string) InterverService {
	return InterverService{
		BaseDeviceService{
			HOST: HOST,
			PORT: PORT,
			API:  API,
		},
	}
}

func (s *InterverService) GetCurrentPower() (*InverterResponse, error) {

	uri, err := s.GetDeviceUrl()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	var power InverterResponse
	err = json.NewDecoder(resp.Body).Decode(&power)

	if err != nil {
		return nil, err
	}

	return &power, err

}
