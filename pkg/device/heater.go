package device

import (
	"encoding/json"
	"fmt"
	"home-solar-pi/pkg/utils"
	"io"
	"log"
	"net/http"
)

type HeaterService struct {
	BaseDeviceService
	HEATER_TOGGLE int
}

func NewHeaterService(HOST string, PORT int, API string, HEATER_TOGGLE int, logger *log.Logger) HeaterService {
	return HeaterService{
		BaseDeviceService{
			HOST:   HOST,
			PORT:   PORT,
			API:    API,
			logger: logger,
		},
		HEATER_TOGGLE,
	}
}

func (s *HeaterService) PowerOn() (any, error) {
	return s.changePower(true)
}

func (s *HeaterService) PowerOff() (any, error) {
	return s.changePower(false)
}

func (s *HeaterService) changePower(on bool) (any, error) {

	if utils.Inactive {
		s.logger.Println("Fake Heater POWER ON")
		return "", nil
	}

	deviceUrl, _ := s.GetDeviceUrl()

	uri := fmt.Sprintf("%s/Switch.Set?id=0&on=%t", deviceUrl, on)
	if on {
		uri = fmt.Sprintf("%s&toggle_after=%d", uri, s.HEATER_TOGGLE)
	}
	resp, err := http.Get(uri)

	if err != nil {
		return "", err
	}

	// var x any = interface{}
	// err = json.NewDecoder(resp.Body).Decode(x)
	bytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (s *HeaterService) GetStatus() (bool, error) {
	deviceUrl, _ := s.GetDeviceUrl()

	uri := fmt.Sprintf("%s/Shelly.GetStatus", deviceUrl)
	resp, err := http.Get(uri)

	if err != nil {
		return false, err
	}

	var response HeaterStatusResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}

	return response.Switch0.Output, nil

}

type HeaterStatusResponse struct {
	Ble struct {
	} `json:"ble"`
	Cloud struct {
		Connected bool `json:"connected"`
	} `json:"cloud"`
	Mqtt struct {
		Connected bool `json:"connected"`
	} `json:"mqtt"`
	PlugsUI struct {
	} `json:"plugs_ui"`
	Switch0 struct {
		ID      int     `json:"id"`
		Source  string  `json:"source"`
		Output  bool    `json:"output"`
		Apower  float64 `json:"apower"`
		Voltage float64 `json:"voltage"`
		Current float64 `json:"current"`
		Aenergy struct {
			Total    float64   `json:"total"`
			ByMinute []float64 `json:"by_minute"`
			MinuteTs int       `json:"minute_ts"`
		} `json:"aenergy"`
		Temperature struct {
			TC float64 `json:"tC"`
			TF float64 `json:"tF"`
		} `json:"temperature"`
	} `json:"switch:0"`
	Sys struct {
		Mac              string `json:"mac"`
		RestartRequired  bool   `json:"restart_required"`
		Time             string `json:"time"`
		Unixtime         int    `json:"unixtime"`
		Uptime           int    `json:"uptime"`
		RAMSize          int    `json:"ram_size"`
		RAMFree          int    `json:"ram_free"`
		FsSize           int    `json:"fs_size"`
		FsFree           int    `json:"fs_free"`
		CfgRev           int    `json:"cfg_rev"`
		KvsRev           int    `json:"kvs_rev"`
		ScheduleRev      int    `json:"schedule_rev"`
		WebhookRev       int    `json:"webhook_rev"`
		AvailableUpdates struct {
			Stable struct {
				Version string `json:"version"`
			} `json:"stable"`
		} `json:"available_updates"`
	} `json:"sys"`
	Wifi struct {
		StaIP  string `json:"sta_ip"`
		Status string `json:"status"`
		Ssid   string `json:"ssid"`
		Rssi   int    `json:"rssi"`
	} `json:"wifi"`
	Ws struct {
		Connected bool `json:"connected"`
	} `json:"ws"`
}
