module BossBotApp

replace BossBot v0.0.0 => ./BossBot

replace Utilities v0.0.0 => ./Utilities

require (
	BossBot v0.0.0
	Utilities v0.0.0
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/Rayer/chatbot v0.3.0
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/lusis/go-slackbot v0.0.0-20180109053408-401027ccfef5 // indirect
	github.com/lusis/slack-test v0.0.0-20180109053238-3c758769bfa6 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.1
	google.golang.org/appengine v1.4.0 // indirect
)
