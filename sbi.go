package slackbi

import (
	"fmt"
	aface "github.com/gonyyi/aFace/logger"
	"github.com/gonyyi/aninterface"
)

type SBI struct {
	config *config
	logger aninterface.Logger1a
}

const (
	SBI_VERSION                = "1.0.0 (2020-12-21)"
	SLACK_ENDPOINT_MSG         = "https://slack.com/api/chat.postMessage"
	SLACK_ENDPOINT_FILE_UPLOAD = "https://slack.com/api/files.upload"
)

func New(c *config, logger aface.Logger1a) (*SBI, error) {
	if c == nil {
		return nil, fmt.Errorf("config is empty")
	}
	b := &SBI{
		config: c,
	}

	if c.Logging && logger != nil {
		b.logger = logger
	} else {
		b.logger = &aface.DummyLogger1a{} // to prevent null ptr error
	}

	b.logger.Infof("SlackBotInterface %s", SBI_VERSION)
	b.logger.Infof("Creating %s: %s", b.config.BotName, b.config.BotVersion)

	b.config.Validate()

	return b, nil
}

func (b *SBI) Run() error {
	b.logger.Infof("Serving HTTP <%s%s%s>", b.config.Service.Address, b.config.Service.Port, b.config.Service.Path)
	return b.serve()
}
