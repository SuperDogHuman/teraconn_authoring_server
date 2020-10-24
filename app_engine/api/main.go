package main

import (
	"os"

	"github.com/super-dog-human/teraconnectgo"
)

func main() {
	if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		teraconnectgo.Main(appEnv)
	}
}
