package main

import (
	"flag"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/node-real/greenfield-bundle-service/bundler"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	flagConfigPath = "config-path"
)

func initFlags() {
	flag.String(flagConfigPath, "", "config path")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(fmt.Sprintf("bind flags error, err=%s", err))
	}
}

func printUsage() {
	fmt.Print("usage: ./bundler --config-path config_file_path\n")
}

func main() {
	initFlags()

	configFilePath := viper.GetString(flagConfigPath)
	if configFilePath == "" {
		printUsage()
		return
	}
	config := util.ParseServerConfigFromFile(configFilePath)

	util.InitLogger(config.LogConfig)

	db, err := database.ConnectDBWithConfig(config.DBConfig)
	if err != nil {
		util.Logger.Errorf("connect database error, err=%s", err.Error())
		return
	}

	bundler, err := bundler.NewBundler(config, db)
	if err != nil {
		util.Logger.Errorf("new bundler error, err=%s", err.Error())
		return
	}
	bundler.Run()
}
