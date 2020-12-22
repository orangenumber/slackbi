package slackbi

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/aninterface"
	"strings"
)

type SlackBot struct {
	config config
	logger aninterface.Logger1a
}

const (
	VERSION = "1.0.0 (2020-12-21)"

	DEFAULT_BOT_NAME           = "TestBot"
	DEFAULT_BOT_VERSION        = "0.0.1"
	DEFAULT_SERVICE_PORT       = 8080
	DEFAULT_SERVICE_PATH       = "/"
	DEFAULT_SERVICE_MODULE_DIR = "./modules"
)

func New(c config) *SlackBot {
	b := &SlackBot{
		config: c,
	}
	if c.Logging.Enable {
		b.logger = alog.New(c.Logging.output, "sbi ", alog.FDefault|alog.FDateYYYYMMDD|alog.FPrefix|alog.FLevelTrace)
	} else {
		b.logger = &aninterface.DummyLogger1a{} // to prevent null ptr error
	}

	b.logger.Infof("SlackBotInterface %s", VERSION)
	b.logger.Infof("Creating %s: %s", b.config.BotName, b.config.BotVersion)

	// check service
	if b.config.Service.Port == 0 {
		b.logger.Warnf("Invalid service port (:%d) -> using default port (:%d) instead", b.config.Service.Port, DEFAULT_SERVICE_PORT)
		b.config.Service.Port = DEFAULT_SERVICE_PORT
	}
	if !strings.HasPrefix(b.config.Service.Path, "/") {
		newPath := "/" + b.config.Service.Path
		b.logger.Warnf("Invalid service path (%s) -> using (%s) instead", b.config.Service.Path, newPath)
		b.config.Service.Path = newPath
	}

	return b
}

func (b *SlackBot) Run() error {
	b.logger.Infof("Serving HTTP %s:%d%s", b.config.Service.Address, b.config.Service.Port, b.config.Service.Path)
	b.serve()
	return nil
}
