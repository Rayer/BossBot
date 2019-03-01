module BossBotApp

replace BossBot v0.0.0 => ./BossBot

replace Utilities v0.0.0 => ./Utilities

replace ChatBot v0.0.0 => ./ChatBot

replace SlackChatBot v0.0.0 => ./SlackChatBot

require (
	BossBot v0.0.0
	ChatBot v0.0.0
	SlackChatBot v0.0.0
	Utilities v0.0.0
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/lusis/go-slackbot v0.0.0-20180109053408-401027ccfef5 // indirect
	github.com/lusis/slack-test v0.0.0-20180109053238-3c758769bfa6 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.2.0
	google.golang.org/appengine v1.4.0 // indirect
)
