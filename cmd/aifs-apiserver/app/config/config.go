/*
 * Created on Wed Jul 05 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package config

import (
	"github.com/jacklv111/common-sdk/config"
	"github.com/jacklv111/common-sdk/database"
	"github.com/jacklv111/common-sdk/database/gorm"
	"github.com/jacklv111/common-sdk/env"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/s3"
	"github.com/spf13/pflag"
)

// ServerConfig aifs api server config.
type serverConfig struct {
	config.Configs
	restPort int
}

var ServerConfig serverConfig

func init() {
	ServerConfig = serverConfig{}

	ServerConfig.AddConfig(env.EnvConfig)
	ServerConfig.AddConfig(log.LogConfig)
	ServerConfig.AddConfig(gorm.GormConfig)
	ServerConfig.AddConfig(database.DbConfig)
	ServerConfig.AddConfig(s3.S3Config)
}

// restPort getter
func (serverCfg serverConfig) GetRestPort() int {
	return serverCfg.restPort
}

func (serverCfg serverConfig) GetFlags() (flagSet *pflag.FlagSet) {
	flagSet = serverCfg.Configs.GetFlags()
	flagSet.IntVar(&ServerConfig.restPort, "rest-port", 8080, "rest server port.")
	return
}
