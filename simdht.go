package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/midoks/simdht/internal/cmd"
	"github.com/midoks/simdht/internal/conf"
)

const Version = "0.0.1"
const AppName = "simdht"

func main() {

	app := cli.NewApp()
	app.Name = conf.App.Name
	app.Version = conf.App.Version
	app.Usage = "A simple DHT service"
	app.Commands = []cli.Command{
		cmd.Service,
	}

	if err := app.Run(os.Args); err != nil {
		log.Println("Failed to start application: %s", err)
	}
}
