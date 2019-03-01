module BossBotLib

replace Utilities v0.0.0 => ../Utilities

replace ChatBot v0.0.0 => ../ChatBot

replace SlackChatBot v0.0.0 => ../SlackChatBot

require (
	ChatBot v0.0.0
	SlackChatBot v0.0.0
	Utilities v0.0.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/nlopes/slack v0.4.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.3.1
)
