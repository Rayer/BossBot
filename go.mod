module BossBotApp

replace BossBot v0.0.0 => ./BossBot

replace Utilities v0.0.0 => ./Utilities

require BossBot v0.0.0

require (
	github.com/sirupsen/logrus v1.2.0
	google.golang.org/appengine v1.4.0 // indirect
)
