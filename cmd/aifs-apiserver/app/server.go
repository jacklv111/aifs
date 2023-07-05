/*
 * Created on Wed Jul 05 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package app

import (
	"fmt"
	"sync"

	"github.com/jacklv111/aifs/app/apigin"
	"github.com/jacklv111/aifs/cmd/aifs-apiserver/app/config"
	"github.com/jacklv111/common-sdk/cli"
	"github.com/jacklv111/common-sdk/database"
	utilerrors "github.com/jacklv111/common-sdk/errors"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/s3"
	"github.com/spf13/cobra"
)

var waitGroup sync.WaitGroup

// NewApiServerCommand creates a *cobra.Command object with default parameters
func NewApiServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "aifs-apiserver",
		Long: `Run aifs apiserver to provide rest operations`,

		// stop printing usage when the command errors
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			flagSet := cmd.Flags()

			// Activate logging as soon as possible, after that
			// show flags with the final logging configuration.
			if errs := log.ValidateAndApply(log.LogConfig); len(errs) > 0 {
				return utilerrors.NewAggregate(errs)
			}
			cli.PrintFlags(flagSet)

			// validate configs
			if errs := config.ServerConfig.Validate(); len(errs) > 0 {
				return utilerrors.NewAggregate(errs)
			}

			return run()
		},
	}

	cmdFlagSet := cmd.Flags()
	serverFlagSet := config.ServerConfig.GetFlags()
	cmdFlagSet.AddFlagSet(serverFlagSet)

	return cmd
}

// run starts servers, init components with configs.
func run() error {
	if err := database.InitDb(); err != nil {
		return err
	}

	if err := s3.InitS3(); err != nil {
		return err
	}

	// add scheduler

	go startRestServer()

	waitGroup.Add(1)
	waitGroup.Wait()
	return nil
}

func startRestServer() {
	defer waitGroup.Done()
	router := apigin.NewRouter()
	err := router.Run(fmt.Sprintf(":%d", config.ServerConfig.GetRestPort()))
	if err != nil {
		log.Errorf("start rest server error %s", err)
	}
}
