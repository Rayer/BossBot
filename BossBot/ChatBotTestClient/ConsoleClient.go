package main

import (
	"BossBotLib"
	"bufio"
	"fmt"
	"github.com/Rayer/chatbot"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	logrus.SetLevel(logrus.WarnLevel)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter ID [BotSpec] : ")
	id, err := reader.ReadString('\n')
	//Trim \n
	id = strings.Replace(id, "\n", "", -1)

	if err != nil {
		panic(err.Error())
	}
	if id == "" {
		id = "BotSpec"
	}
	fmt.Println("Welcome " + id + ", start invoking session...")
	ctx := ChatBot.NewContextManager()
	utx := ctx.CreateUserContext(id, func() ChatBot.Scenario {
		return &BossBot.RootScenario{}
	})

	fmt.Println(utx.RenderMessage())
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "exitloop" {
			break
		}
		fmt.Println(utx.HandleMessage(text))
		fmt.Println(utx.RenderMessage())
	}
}
