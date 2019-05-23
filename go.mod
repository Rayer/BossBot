module BossBotApp

replace BossBot v0.0.0 => ./BossBot

replace Utilities v0.0.0 => ./Utilities

require (
	BossBot v0.0.0
	Utilities v0.0.0
	github.com/Rayer/chatbot v0.3.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.1
)
