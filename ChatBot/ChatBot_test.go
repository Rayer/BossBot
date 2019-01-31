package ChatBot

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestEssentials(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := NewContextManager()
	uc := ctx.GetUserContext("rayer")
	t.Log(uc.RenderMessage())
	t.Log(uc.HandleMessage("submit report"))
	t.Log(uc.RenderMessage())
	t.Log(uc.HandleMessage("Let's exit"))
	t.Log(uc.RenderMessage())

}
