package main

import (
	"os"

	"github.com/SuperDogHuman/teraconnectgo"
)

func main() {
	if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		teraconnectgo.Main(appEnv)
	}
}
