/*
 * Created on Wed Jul 05 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/jacklv111/aifs/cmd/aifs-apiserver/app"
	"github.com/jacklv111/aifs/cmd/aifs-apiserver/app/config"
	"github.com/jacklv111/common-sdk/cli"
	"github.com/jacklv111/common-sdk/utils"
)

func main() {
	utils.StartProf()
	// pre run
	// read config from file
	if err := config.ServerConfig.ReadFromFile(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	command := app.NewApiServerCommand()
	code := cli.Run(command)
	os.Exit(code)
}
