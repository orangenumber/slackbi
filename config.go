package slackbi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type config struct {
	SlackToken string `json:"slack_token"`
	BotVersion string `json:"bot_version"`
	BotName    string `json:"bot_name"`
	Logging    struct {
		Enable   bool `json:"enable"`
		output   io.Writer
		file     *os.File
		Filename string `json:"filename,omitempty"`
		ToFile   bool   `json:"to_file,omitempty"`
		ToStdErr bool   `json:"to_std_err,omitempty"`
		ToStdOut bool   `json:"to_std_out,omitempty"`
	} `json:"logging"`
	Service struct {
		Address        string `json:"address"`
		Port           int    `json:"port"`
		Path           string `json:"path"`
		ModuleDir      string `json:"module_directory"`
		AcceptChallege bool   `json:"accept_challenge"` // this is one time URL verification from slack. (https://api.slack.com/events/url_verification)
	} `json:"service"`
}

func (c *config) LogToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	c.Logging.file = f
	c.Logging.output = f
	c.Logging.ToFile = true
	c.Logging.Enable = true
	c.Logging.ToStdErr = false
	c.Logging.ToStdOut = false
	return nil
}

func (c *config) LogToStdOut() {
	c.Logging.Enable = true
	c.Logging.file = nil
	c.Logging.output = os.Stdout
	c.Logging.ToFile = false
	c.Logging.ToStdErr = false
	c.Logging.ToStdOut = true
}

func (c *config) LogToStdErr() {
	c.Logging.Enable = true
	c.Logging.ToFile = false
	c.Logging.ToStdErr = true
	c.Logging.ToStdOut = false
	c.Logging.file = nil
	c.Logging.output = os.Stderr
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
	c := config{
		BotName:    DEFAULT_BOT_NAME,
		BotVersion: fmt.Sprintf("%s (%s)", DEFAULT_BOT_VERSION, time.Now().Format("2006-01-02 15:04:05")),
	}
	c.Service.Port = DEFAULT_SERVICE_PORT

	c.Service.ModuleDir = DEFAULT_SERVICE_MODULE_DIR
	c.Service.Path = DEFAULT_SERVICE_PATH
	c.LogToStdErr()
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
		if c.Logging.Enable {
			if c.Logging.ToFile {
				if err := c.LogToFile(c.Logging.Filename); err != nil {
					return nil, err
				}
			} else if c.Logging.ToStdOut {
				c.LogToStdOut()
			} else if c.Logging.ToStdErr {
				c.LogToStdErr()
			}
		}
		return &c, nil
	}
}
