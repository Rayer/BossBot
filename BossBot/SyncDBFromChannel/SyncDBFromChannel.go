package main

import (
	"BossBotLib"
	"fmt"
)

func main() {
	conf, err := BossBot.CreateConfigurationFromFile()

	//Not really necessary, it will panic itself lol
	if err != nil {
		panic(err)
	}

	app := conf.ServiceContext.SlackClient

	//Get channel list, and get channel named "general"
	channels, err := app.GetChannels(true)
	if err != nil {
		panic(err)
	}

	db := conf.ServiceContext.DBObject.GetDB()
	for _, channel := range channels {
		//fmt.Printf("Handling channel : %+v\n", channel)
		if channel.Name == "general" {
			fmt.Printf("Handling channel : %+v\n", channel)
			members := channel.Members
			for _, member := range members {
				userProfile, err := app.GetUserProfile(member, false)
				if err != nil {
					panic(err)
				}
				fmt.Printf("We have %s / %s / %s\n", userProfile.RealName, userProfile.DisplayName, member)
				fmt.Printf("Raw data : %+v \n", userProfile)
				_, err = db.Exec("update mcds_tw_members set slack_uid = ?, slack_id = ? where slack_name = ?", member, userProfile.RealName, userProfile.DisplayName)
				if err != nil {
					panic(err)
				}
			}
		}
	}

}
