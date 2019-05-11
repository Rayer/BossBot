module BossBotLib

replace Utilities v0.0.0 => ../Utilities

replace github.com/Rayer/chatbot v0.0.0 => ../ChatBot

replace SlackChatBot v0.0.0 => ../SlackChatBot

require (
	SlackChatBot v0.0.0
	Utilities v0.0.0
	github.com/Rayer/chatbot v0.1.1
	github.com/go-sql-driver/mysql v1.4.1
	github.com/nlopes/slack v0.4.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/viper v1.3.1
)
