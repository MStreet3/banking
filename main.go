package main

import (
	"github.com/mstreet3/banking/app"
	"github.com/mstreet3/banking/logger"
)

func main() {
	logger.Info("Starting the application")
	app.Start()
}
