package main // import "BossBotApp"

import (
	"BossBot"
	"Utilities"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type CustomFormatter struct{}

func (CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	ret := fmt.Sprintf("[%s] (%s) : %s\n", entry.Level.String(), entry.Time.Format(time.RFC3339), entry.Message)
	return []byte(ret), nil
}

func main() {
	conf, _ := BossBot.CreateConfigurationFromFile()

	//Setup Logger
	log.SetLevel(log.Level(conf.LogLevel))
	log.SetReportCaller(true)
	if conf.LogFilePath == "" {
		log.SetOutput(os.Stdout)
		log.SetFormatter(&CustomFormatter{})
	} else {
		f, err := os.OpenFile(conf.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Warn("Error opening log file : " + conf.LogFilePath + " and error is : " + err.Error())
			log.Warn("Log to logfile failed, now log to stdout instead")
		} else {
			log.Infof("Now logging to : %s", f.Name())
			log.SetOutput(f)
			log.SetFormatter(&log.JSONFormatter{})
		}

	}

	log.Printf("Start BossBot with configuration : %+v\n", *conf)
	Utilities.ExecuteCode(conf.PIDFilePath, func() {

		go func() {
			BossBot.StartBroadcaster(*conf)
		}()
		BossBot.RespServer(*conf)
		select {}
	})

}
