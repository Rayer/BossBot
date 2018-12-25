package BossBot

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configuration struct {
	LogLevel 	uint32
	SqlHost 	string
	SqlPort 	uint32
	SqlAcc 		string
	SqlPass 	string
	PIDFilePath	string
}

func CreateConfigurationFromFile() (*Configuration, error) {
	viper.SetConfigName("BossBot")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Can't read configuration file! "))
	}

	viper.SetDefault("SqlPort", 3306)

	conf := &Configuration{}
	conf.LogLevel = uint32(viper.GetInt32("LogLevel"))
	conf.PIDFilePath = viper.GetString("PIDFilePath")
	conf.SqlHost = viper.GetString("SqlHost")
	conf.SqlPort = uint32(viper.GetInt32("SqlPort"))
	conf.SqlAcc = viper.GetString("SqlAcc")
	conf.SqlPass = viper.GetString("SqlPass")

	return conf, nil
}