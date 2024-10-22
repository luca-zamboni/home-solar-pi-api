package main

import (
	"fmt"
	"home-solar-pi/pkg/api"
	"home-solar-pi/pkg/db"
	"home-solar-pi/pkg/device"
	"home-solar-pi/pkg/utils"
	"home-solar-pi/pkg/worker"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	INVERTER_HOST   string
	INVERTER_PORT   int
	INVERTER_API    string
	UPDATE_INTERVAL int

	HEATER_HOST            string
	HEATER_PORT            int
	HEATER_API             string
	HEATER_TOGGLE_INTERVAL int

	PG_USER string
	PG_PASS string
	PG_HOST string
	PG_PORT int
	PG_NAME string

	INVERTER_THRESHOLD int
)

func main() {

	envPath := ".env"
	if len(os.Args) == 2 {
		envPath = os.Args[1]
	}

	initVar(envPath)

	logger := log.Default()

	inverterDevice := device.NewInterver(device.DeviceConfig{
		Host:    INVERTER_HOST,
		Port:    INVERTER_PORT,
		ApiPath: INVERTER_API,
	})
	heaterDevice := device.NewHeater(device.DeviceConfig{
		Host:    HEATER_HOST,
		Port:    HEATER_PORT,
		ApiPath: HEATER_API,
	}, HEATER_TOGGLE_INTERVAL)

	devices := make(map[string]device.Device)

	devices["inverter"] = inverterDevice
	devices["heater"] = heaterDevice

	apiServer := api.NewApiServer(devices)

	dbService, err := db.New(db.PostresConf{
		User: PG_USER,
		Pass: PG_PASS,
		Host: PG_HOST,
		Port: PG_PORT,
		Name: PG_NAME,
	})

	if err != nil {
		panic("Error connecting to Postgres")
	}

	workerService := worker.NewHeaterInverterWorker(&inverterDevice, &heaterDevice, logger, dbService, INVERTER_THRESHOLD)

	updaterInterval := time.Second * time.Duration(UPDATE_INTERVAL)

	// worker works async
	go workerService.StartHeaterInverterCycle(updaterInterval)

	// listening for api
	apiServer.StartServer()

}

func initVar(envPath string) {

	godotenv.Load(envPath)

	utils.InitGlobals()

	UPDATE_INTERVAL = getVarint("UPDATE_INTERVAL")

	INVERTER_HOST = getVarString("INTERTER_HOST")
	INVERTER_PORT = getVarint("INTERTER_PORT")
	INVERTER_API = getVarString("INTERTER_API")

	HEATER_HOST = getVarString("HEATER_HOST")
	HEATER_PORT = getVarint("HEATER_PORT")
	HEATER_API = getVarString("HEATER_API")
	HEATER_TOGGLE_INTERVAL = getVarint("HEATER_TOGGLE")

	PG_USER = getVarString("PG_USER")
	PG_PASS = getVarString("PG_PASS")
	PG_HOST = getVarString("PG_HOST")
	PG_PORT = getVarint("PG_PORT")
	PG_NAME = getVarString("PG_NAME")

	INVERTER_THRESHOLD = getVarint("INVERTER_THRESHOLD")
}

func getVarString(varName string) string {
	return os.Getenv(varName)
}

func getVarint(varName string) int {
	varInt, err := strconv.Atoi(os.Getenv(varName))
	if err != nil {
		fmt.Printf("Not a number :%s %s\n", varName, os.Getenv(varName))
		panic(fmt.Sprintf("Not a number :%s %s", varName, os.Getenv(varName)))
	}

	return varInt
}
