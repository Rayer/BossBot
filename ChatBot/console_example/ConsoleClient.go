package main

import (
	"ChatBot"
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logrus.SetLevel(logrus.WarnLevel)
	ctx := ChatBot.NewContextManager()
	utx := ctx.CreateUserContext("BotSpec", func() ChatBot.Scenario {
		return &RootScenario{}
	})

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(utx.RenderMessage())
	for {
		text, _ := reader.ReadString('\n')
		if text == "exitloop" {
			break
		}
		fmt.Println(utx.HandleMessage(text))
		fmt.Println(utx.RenderMessage())
	}
}
