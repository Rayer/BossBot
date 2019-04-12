package BossBot

import (
	"ChatBot"
	"Utilities"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/spf13/viper"
	"reflect"
)

type Configuration struct {
	LogLevel          int64
	LogFilePath       string
	SqlHost           string
	SqlPort           int64
	SqlAcc            string
	SqlPass           string
	PIDFilePath       string
	SlackAppToken     string
	SlackBotToken     string
	SlackVerifyToken  string
	ServiceContext    ServiceContext
	ChatBotResetTimer int64
}

type ServiceContext struct {
	conf        *Configuration
	SlackClient *slack.Client
	//SlackRTM    *slack.RTM
	DBObject      *Utilities.DBObject
	ChatBotClient *ChatBot.ContextManager
}

var globalConfig *Configuration

func GetConfiguration() *Configuration {
	if globalConfig == nil {
		_, err := CreateConfigurationFromFile()
		if err != nil {
			panic("Fail to load configuration from file!")
		}
	}

	return globalConfig
}

//TODO: Error will always be NOT nil, need to improve it?
func CreateConfigurationFromFile() (*Configuration, error) {
	viper.SetConfigName("BossBot")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Can't read configuration file (%s) ", err))
	}

	viper.SetDefault("SqlPort", 3306)
	viper.SetDefault("PIDFilePath", ".")
	viper.SetDefault("ChatBotResetTimer", 300)

	conf := &Configuration{}
	v := reflect.ValueOf(conf).Elem()

	for i, n := 0, v.NumField(); i < n; i++ {
		fieldName := v.Type().Field(i).Name
		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(viper.GetString(fieldName))
		case reflect.Int32:
			v.Field(i).SetInt(int64(viper.GetInt32(fieldName)))
		case reflect.Int64:
			v.Field(i).SetInt(viper.GetInt64(fieldName))
		case reflect.Struct:
		default:
			continue
		}
	}

	conf.ServiceContext.SlackClient = slack.New(conf.SlackAppToken)
	conf.ServiceContext.DBObject, err = Utilities.CreateDBObject(conf.SqlHost, conf.SqlAcc, conf.SqlPass)
	if err != nil {
		panic(fmt.Errorf("error creating DB object! Please check sql credential and address! (%s)", err))
	}
	conf.ServiceContext.ChatBotClient = ChatBot.NewContextManager()
	//rtm := conf.ServiceContext.SlackClient.NewRTM()
	//go rtm.ManageConnection()

	globalConfig = conf

	return conf, nil
}
