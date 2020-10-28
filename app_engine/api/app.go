package main

import (
	"os"

	"github.com/super-dog-human/teraconnectgo/interface/handler"
)

func main() {
	if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		handler.Main(appEnv)
	}
}
