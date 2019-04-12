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

	log.SetReportCaller(true)
	log.SetLevel(log.Level(conf.LogLevel))
	log.SetOutput(os.Stdout)
	log.SetFormatter(&CustomFormatter{})
	log.Printf("Start BossBot with configuration : %+v\n", *conf)
	Utilities.ExecuteCode(conf.PIDFilePath, func() {

		go func() {
			BossBot.StartBroadcaster(*conf)
		}()
		BossBot.RespServer(*conf)
		select {}
	})

}
