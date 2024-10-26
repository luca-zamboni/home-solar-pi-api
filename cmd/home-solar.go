package main

import (
	"errors"
	"fmt"
	"home-solar-pi/pkg/api"
	"home-solar-pi/pkg/device"
	"home-solar-pi/pkg/rule"
	"home-solar-pi/pkg/utils"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	HOME_PATH   string
	DEVICE_PATH string
	RULES_PATH  string

	logger = log.Default()
)

const (
	DEVICE_SUBPATH = "devices"
	RULES_SUBPATH  = "rules"
)

func main() {

	envPath := ".env"
	if len(os.Args) == 2 {
		envPath = os.Args[1]
	}

	initVar(envPath)

	err := ensureDir(HOME_PATH)
	if err != nil {
		panic(err)
	}

	err = ensureDir(DEVICE_PATH)
	if err != nil {
		panic(err)
	}

	// if err != nil {
	// 	panic("Error connecting to Postgres")
	// }

	deviceManager, err := device.NewDeviceManager(DEVICE_PATH)

	if err != nil {
		panic(err)
	}

	all, _ := deviceManager.GetAllDevices()

	logger.Printf("Devices := %+v\n", all)

	apiServer := api.NewApiServer(deviceManager)
	apiServer = apiServer

	// // workerService := worker.NewHeaterInverterWorker(&inverterDevice, &heaterDevice, logger, dbService, INVERTER_THRESHOLD)

	// // updaterInterval := time.Second * time.Duration(UPDATE_INTERVAL)

	// // worker works async
	// // go workerService.StartHeaterInverterCycle(updaterInterval)

	ruleManager := rule.NewRuleManager(RULES_PATH, *deviceManager)

	logger.Printf("Rules := %+v", ruleManager.GetAllRules())

	panicChan := make(chan error)

	go ruleManager.StartRuleServer(panicChan)

	panicMessage := <-panicChan

	panic(panicMessage.Error())

	// // listening for api
	// apiServer.StartServer()

}

func initVar(envPath string) {

	err := godotenv.Load(envPath)
	if err != nil {
		panic(err)
	}

	utils.InitGlobals()

	HOME_PATH = getVarString("HOME_PATH")
	DEVICE_PATH = path.Join(HOME_PATH, DEVICE_SUBPATH)
	RULES_PATH = path.Join(HOME_PATH, RULES_SUBPATH)
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

func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}
