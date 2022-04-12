package cmd

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/midoks/simdht/internal/app"
	"github.com/midoks/simdht/internal/app/router"
	"github.com/midoks/simdht/internal/conf"
)

var Service = cli.Command{
	Name:        "service",
	Usage:       "This command starts all services",
	Description: `Start DHT services`,
	Action:      runAllService,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}

func runAllService(c *cli.Context) error {

	err := router.Init("")
	fmt.Println("runAllService:", err)
	if err != nil {
		return err
	}
	app.Start(conf.Web.HttpPort)
	return nil
}
