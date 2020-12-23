package slackbi

import (
	"encoding/json"
	logface "github.com/gonyyi/aface/logger"
	"io/ioutil"
)

type config struct {
	SlackToken string `json:"slack_token"`
	BotVersion string `json:"bot_version"`
	BotName    string `json:"bot_name"`
	Logging    bool   `json:"logging"`
	logger     logface.Logger1a
	Service    struct {
		Address        string `json:"address"`
		Port           string `json:"port"`
		Path           string `json:"path"`
		AcceptChallege bool   `json:"accept_challenge"` // this is one time URL verification from slack. (https://api.slack.com/events/url_verification)
	} `json:"service"`
	Module struct {
		Dir  string `json:"directory"`
		Conf string `json:"config_filename"`
	} `json:"module"`
}

func (c *config) Save(filename string) error {
	j, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, j, 0755); err != nil {
		return err
	}

	return nil
}

func NewConfig() *config {
	c := config{}
	c.Default()
	return &c
}

func ReadConfig(filename string) (*config, error) {
	if b, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		c := config{}
		if err := json.Unmarshal(b, &c); err != nil {
			return nil, err
		}
		return &c, nil
	}
}
