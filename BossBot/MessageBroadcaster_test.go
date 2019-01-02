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

/*
[debug] (2019-01-02T06:01:37Z) : Evaluating 4 items...
[fatal] (2019-01-02T06:01:37Z) : Broadcast fail!: Error parsing time: parsing time "2019 Jan 2 15:00:00 +0800 CST" as "2006 Jan 02 15:04:05 -0700 MST": cannot parse "2 15:00:00 +0800 CST" as "02"
*/
func TestParsingDateTime(t *testing.T) {
	source1 := "2019 Jan 2 15:00:00 +0800 CST"
	source2 := "2019 Jan 12 15:00:00 +0800 CST"
	layout := "2006 Jan _2 15:04:05 -0700 MST" // change from 2006 Jan 02 15:04:05 -0700 MST
	_, err := time.Parse(layout, source1)
	if err != nil {
		t.Fatalf("Error parsing : %s", err)
	}
	_, err = time.Parse(layout, source2)
}
