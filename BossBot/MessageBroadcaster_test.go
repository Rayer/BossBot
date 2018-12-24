package BossBot

import (
	"testing"
)

func TestProcessing(t *testing.T) {
	err := Processing()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartBroadcaster(t *testing.T) {
	StartBroadcaster()
}

