package BossBot

import (
	"Utilities"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/spf13/viper"
)

type Configuration struct {
	LogLevel      uint32
	SqlHost       string
	SqlPort       uint32
	SqlAcc        string
	SqlPass       string
	PIDFilePath   string
	SlackBotToken string
	Context       ServiceContext
}

type ServiceContext struct {
	SlackClient *slack.Client
	DBObject    *Utilities.DBObject
}

func CreateConfigurationFromFile() (*Configuration, error) {
	viper.SetConfigName("BossBot")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Can't read configuration file (%s) ", err))
	}

	viper.SetDefault("SqlPort", 3306)
	viper.SetDefault("PIDFilePath", ".")

	//Stupid.... can we do better?
	conf := &Configuration{}
	conf.LogLevel = uint32(viper.GetInt32("LogLevel"))
	conf.PIDFilePath = viper.GetString("PIDFilePath")
	conf.SqlHost = viper.GetString("SqlHost")
	conf.SqlPort = uint32(viper.GetInt32("SqlPort"))
	conf.SqlAcc = viper.GetString("SqlAcc")
	conf.SqlPass = viper.GetString("SqlPass")
	conf.SlackBotToken = viper.GetString("SlackBotToken")

	conf.Context.SlackClient = slack.New(conf.SlackBotToken)
	conf.Context.DBObject, err = Utilities.CreateDBObject(conf.SqlHost, conf.SqlAcc, conf.SqlPass)
	if err != nil {
		panic(fmt.Errorf("error creating DB object! Please check sql credential and address! (%s)", err))
	}

	return conf, nil
}
