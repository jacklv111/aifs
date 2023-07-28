/*
 * Created on Fri Jul 28 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package config

import (
	datasetZipCfg "github.com/jacklv111/aifs/pkg/job/dataset-zip/config"
	aifsclient "github.com/jacklv111/common-sdk/client/aifs-client"
	basicCfg "github.com/jacklv111/common-sdk/config"
	"github.com/jacklv111/common-sdk/log"
)

type jobConfig struct {
	basicCfg.Configs
}

var JobConfig jobConfig

func init() {
	JobConfig = jobConfig{}

	JobConfig.AddConfig(log.LogConfig)
	JobConfig.AddConfig(aifsclient.AifsConfig)
	JobConfig.AddConfig(datasetZipCfg.DatasetZipConfig)
}
