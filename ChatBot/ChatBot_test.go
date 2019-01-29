package ChatBot

import "testing"

func TestEssentials(t *testing.T) {
	ctx := NewContextManager()
	uc := ctx.GetUserContext("rayer")
	t.Log(uc.RenderMessage())
	t.Log(uc.HandleMessage("Invoke first one"))
	t.Log(uc.RenderMessage())
	t.Log(uc.HandleMessage("Let's exit"))
	t.Log(uc.RenderMessage())

}
