/*
 * Created on Fri Jul 28 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package main

import (
	"fmt"
	"os"
	"time"

	jobcfg "github.com/jacklv111/aifs/cmd/job/dataset-zip/config"
	job "github.com/jacklv111/aifs/pkg/job/dataset-zip"
	"github.com/jacklv111/aifs/pkg/job/dataset-zip/config"
	"github.com/jacklv111/common-sdk/cli"
	aifsclient "github.com/jacklv111/common-sdk/client/aifs-client"
	utilerrors "github.com/jacklv111/common-sdk/errors"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/utils"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "decompression dataset-zip",
		Long: `decompression dataset-zip, parse data according to zip format. Create annotation template, raw data view ,annotation view automatically and upload data to data view`,

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
			if errs := jobcfg.JobConfig.Validate(); len(errs) > 0 {
				return utilerrors.NewAggregate(errs)
			}

			return run()
		},
	}

	cmdFlagSet := cmd.Flags()
	jobFlagSet := jobcfg.JobConfig.GetFlags()
	cmdFlagSet.AddFlagSet(jobFlagSet)

	return cmd
}

func run() error {
	if err := aifsclient.InitAifsClient(); err != nil {
		return err
	}
	client := aifsclient.GetAifsClient()
	err := job.UnzipDatasetZipView(client, config.DatasetZipConfig.GetDataViewId(), config.DatasetZipConfig.GetWorkDir())
	if err != nil {
		log.Errorf("decompression dataset zip view failed, err: %v", err)
	}
	return err
}

func main() {
	utils.StartProf()

	start := time.Now()
	if err := jobcfg.JobConfig.ReadFromFile(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	command := NewCommand()
	code := cli.Run(command)
	elapsed := time.Since(start)
	log.Infof("exist with code %d, time cost: %s", code, elapsed)
	os.Exit(code)
}
