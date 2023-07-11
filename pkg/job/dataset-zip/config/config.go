/*
 * Created on Mon Jul 10 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package config

import (
	"strings"

	"github.com/spf13/pflag"
)

type datasetZipConfig struct {
	dataViewId string
	workDir    string
}

func (cfg *datasetZipConfig) GetDataViewId() string {
	splitList := strings.Split(cfg.dataViewId, "@")
	if len(splitList) > 1 {
		return splitList[0]
	}
	return cfg.dataViewId
}

func (cfg *datasetZipConfig) GetWorkDir() string {
	return cfg.workDir
}

func (cfg *datasetZipConfig) ReadFromFile() error {
	return nil
}

func (cfg *datasetZipConfig) AddFlags(flagSet *pflag.FlagSet) {
	flagSet.StringVar(&cfg.dataViewId, "aifs.input.dataset_zip", "", "Value to indicate the dataset zip view id, it must be uuid")
	flagSet.StringVar(&cfg.workDir, "work_dir", "", "Value to indicate the work directory, the data will be pulled to work directory, the data generated in the decompression will also be put in work directory")
}

func (cfg *datasetZipConfig) Validate() []error {
	return []error{}
}

var DatasetZipConfig *datasetZipConfig

func init() {
	DatasetZipConfig = &datasetZipConfig{}
}
