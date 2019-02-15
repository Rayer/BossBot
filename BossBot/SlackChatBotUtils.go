package BossBot

import (
	"ChatBot"
	"github.com/nlopes/slack"
	"regexp"
	"strings"
)

type SlackScenario interface {
	RenderSlackAttachments(string) []slack.Attachment
	ResponseSlackCallbacks()
}

type KeywordAction func(keyword string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) string

type Keyword struct {
	Keyword string
	Action  KeywordAction
}

type KeywordHandler struct {
	keywordList []Keyword
	scenario    ChatBot.Scenario
	state       ChatBot.ScenarioState
}

func NewKeywordHandler(scenario ChatBot.Scenario, state ChatBot.ScenarioState) *KeywordHandler {
	return &KeywordHandler{scenario: scenario, state: state}
}

func (kh *KeywordHandler) RegisterKeyword(keyword *Keyword) {
	if kh.keywordList == nil {
		kh.keywordList = []Keyword{}
	}
}

func (kh *KeywordHandler) GenerateAttachments(input string) slack.Attachment {
	//Find all surrended by "[]"
	r, _ := regexp.Compile(`\[([A-Za-z 0-9_]*)]`)
	keywords := r.FindAllString(input, -1)

	var ret slack.Attachment
	var actions []slack.AttachmentAction

	for _, keywordDefine := range kh.keywordList {
		//TODO: Maybe we should use map to avoid O(n^2)?
		for _, keyword := range keywords {
			if keywordDefine.Keyword == keyword {
				actions = append(actions, slack.AttachmentAction{
					Text: strings.Title(keyword),
					Type: "button",
				})
				break
			}
		}
	}
	ret.Actions = actions
	return ret
}

func (kh *KeywordHandler) ParseAction(input string) error {
	return nil
}
