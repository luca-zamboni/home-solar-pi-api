package main

import (
	"fmt"
	"home-solar-pi/pkg/api"
	"home-solar-pi/pkg/device"
	"home-solar-pi/pkg/utils"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	UPDATE_INTERVAL int

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

	dbService, err := device.New(device.PostresConf{
		User: PG_USER,
		Pass: PG_PASS,
		Host: PG_HOST,
		Port: PG_PORT,
		Name: PG_NAME,
	})

	if err != nil {
		panic("Error connecting to Postgres")
	}

	deviceManager := device.NewDeviceManager(dbService, utils.GetLogger())

	apiServer := api.NewApiServer(&deviceManager)

	// workerService := worker.NewHeaterInverterWorker(&inverterDevice, &heaterDevice, logger, dbService, INVERTER_THRESHOLD)

	// updaterInterval := time.Second * time.Duration(UPDATE_INTERVAL)

	// worker works async
	// go workerService.StartHeaterInverterCycle(updaterInterval)

	// listening for api
	apiServer.StartServer()

}

func initVar(envPath string) {

	godotenv.Load(envPath)

	utils.InitGlobals()

	UPDATE_INTERVAL = getVarint("UPDATE_INTERVAL")

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
