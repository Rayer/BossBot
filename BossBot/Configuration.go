package BossBot

import (
	"ChatBot"
	"Utilities"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/spf13/viper"
)

type Configuration struct {
	LogLevel         uint32
	SqlHost          string
	SqlPort          uint32
	SqlAcc           string
	SqlPass          string
	PIDFilePath      string
	SlackAppToken    string
	SlackBotToken    string
	SlackVerifyToken string
	ServiceContext   ServiceContext
}

type ServiceContext struct {
	SlackClient *slack.Client
	//SlackRTM    *slack.RTM
	DBObject      *Utilities.DBObject
	ChatBotClient *ChatBot.ContextManager
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
	conf.SlackAppToken = viper.GetString("SlackAppToken")
	conf.SlackBotToken = viper.GetString("SlackBotToken")
	conf.SlackVerifyToken = viper.GetString("SlackVerifyToken")

	conf.ServiceContext.SlackClient = slack.New(conf.SlackAppToken)
	conf.ServiceContext.DBObject, err = Utilities.CreateDBObject(conf.SqlHost, conf.SqlAcc, conf.SqlPass)
	if err != nil {
		panic(fmt.Errorf("error creating DB object! Please check sql credential and address! (%s)", err))
	}
	conf.ServiceContext.ChatBotClient = ChatBot.NewContextManager()
	//rtm := conf.ServiceContext.SlackClient.NewRTM()
	//go rtm.ManageConnection()

	return conf, nil
}
