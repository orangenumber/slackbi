package slackbi

import (
	"encoding/json"
	"github.com/gonyyi/aface"
	"io/ioutil"
)

type config struct {
	SlackToken string `json:"slack_token"`
	BotVersion string `json:"bot_version"`
	BotName    string `json:"bot_name"`
	Logging    bool   `json:"logging"`
	logger     aface.Logger1a
	Service    struct {
		SysCommand     string `json:"sys_command"`
		Address        string `json:"address"`
		Port           string `json:"port"`
		Path           string `json:"path"`
		AcceptChallege bool   `json:"accept_challenge"` // this is one time URL verification from slack. (https://api.slack.com/events/url_verification)
	} `json:"service"`
	Module struct {
		Dir  string `json:"directory"`
		Conf string `json:"config_filename"`
	} `json:"module"`
	Admin map[string]string `json:"admin"` // map[ID] = note
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
