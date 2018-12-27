package BossBot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

type CustomFormatter struct{}

func (CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	ret := fmt.Sprintf("[%s]\t (%s) : %s\n", entry.Level.String(), entry.Time.Format(time.RFC3339), entry.Message)
	return []byte(ret), nil
}

func TestProcessing(t *testing.T) {
	//log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&CustomFormatter{})
	conf, _ := CreateConfigurationFromFile()

	log.Debugf("Start testing TestProcessing")

	err := Processing(*conf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartBroadcaster(t *testing.T) {
	conf, _ := CreateConfigurationFromFile()
	StartBroadcaster(*conf)
}
