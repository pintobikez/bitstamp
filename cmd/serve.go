package main

import (
	"fmt"
	"github.com/pintobikez/bitstamp/api"
	"github.com/pintobikez/bitstamp/engine"
	"github.com/robfig/cron"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Start Http Server
func Serve(c *cli.Context) error {

	cfg := new(api.Config)

	if err := LoadConfigFile("config.yaml", cfg); err != nil {
		fmt.Println("Error in loading config file")
		return nil
	}

	a := api.New(cfg)

	// run crontab mode
	if c.Bool("robot") {
		// launch a cron to check everyday for posted items
		cr := cron.New()
		//cr.AddFunc("* 0 */6 * * *", func() { cj.CheckUpdatedReverses("C") })     // checks for Colect updates
		//cr.AddFunc("* 10 */6 * * *", func() { cj.CheckUpdatedReverses("A") })    // checks for Postage updates
		//cr.AddFunc("* */20 * * * *", func() { cj.ReprocessRequestsWithError() }) // checks for Requests with Error and reprocesses them
		cr.Start()
		defer cr.Stop()
	}

	if c.Bool("balance") {
		// check balance
		x, err := a.AccountBalance()
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("%+v\n", x)
	}

	if c.Bool("check") {
		// check open orders
		x, err := a.OpenOrders()
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("%+v\n", x)
	}

	return nil
}

// Loads a Yaml file and returns in a struct mode
func LoadConfigFile(filename string, c interface{}) error {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(file, c); err != nil {
		return err
	}

	return nil
}
