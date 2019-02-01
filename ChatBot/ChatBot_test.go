package ChatBot

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestEssentials(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := NewContextManager()
	uc := ctx.GetUserContext("rayer")
	sendMessage(uc, "submit report", t)
	sendMessage(uc, "create report", t)
	sendMessage(uc, "MCDS-12345 ggggg", t)
	sendMessage(uc, "MCDS-12346 rrggg", t)
	sendMessage(uc, "good for now", t)
	sendMessage(uc, "MCDS-12245 aaaaa", t)
	sendMessage(uc, "MCDS-12446 fffff", t)
	sendMessage(uc, "good for now", t)
	sendMessage(uc, "submit", t)
	t.Log(uc.RenderMessage())

}

func sendMessage(uc *UserContext, msg string, t *testing.T) {
	t.Log(uc.RenderMessage())
	t.Log("-> " + msg)
	t.Log(uc.HandleMessage(msg))
}
