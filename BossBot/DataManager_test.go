package BossBot

import (
	"log"
	"testing"
)

func TestDBReflect(t *testing.T) {
	conf, err := CreateConfigurationFromFile()
	if err != nil {
		log.Fatal(err)
	}
	dm, err := GetDataManager(conf)
	if err != nil {
		log.Fatal(err)
	}
	ret, err := dm.GetBroadcastList()
	t.Logf("%+v", ret)
}
