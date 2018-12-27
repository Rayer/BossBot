module BossBotLib

require github.com/pkg/errors v0.8.0

replace Utilities v0.0.0 => ../Utilities

require (
	Utilities v0.0.0
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/nlopes/slack v0.4.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.3.1
)
