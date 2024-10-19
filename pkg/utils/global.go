package utils

import "os"

var (
	Inactive bool
)

func InitGlobals() {
	Inactive = os.Getenv("INACTIVE") == "true"
}