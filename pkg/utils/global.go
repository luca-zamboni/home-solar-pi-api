package utils

import (
	"log"
	"os"
)

var (
	Debug bool
)

func InitGlobals() {
	Debug = os.Getenv("DEBUG") == "true"
}

func GetLogger() *log.Logger {
	return log.Default()
}
