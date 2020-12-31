package slackbi

import (
	"fmt"
	"github.com/gonyyi/aface"
	"github.com/orangenumber/slackbi/lib/module"
)

type SBI struct {
	config  *config
	logger  aface.Logger1a
	modules *module.Modules
	sys     func(MsgIncoming)
}

const (
	SBI_VERSION                = "0.0.6 (2020-12-31)"
	SLACK_ENDPOINT_MSG         = "https://slack.com/api/chat.postMessage"
	SLACK_ENDPOINT_FILE_UPLOAD = "https://slack.com/api/files.upload"
	SYS_COMMAND                = "sys" // todo: let this be customizable from the config?
)

func New(c *config, logger aface.Logger1a) (*SBI, error) {
	if c == nil {
		return nil, fmt.Errorf("config is empty")
	}
	b := &SBI{
		config: c,
		sys:    func(MsgIncoming) {}, // for default, do nothing unless specified.
	}

	if c.Logging && logger != nil {
		b.logger = logger
	} else {
		b.logger = &aface.LoggerDummy1a{} // to prevent null ptr error
	}

	b.logger.Infof(MF_SBI_SVersion.Format(SBI_VERSION))
	b.logger.Infof(MF_SBI_CREATING_SName_SVersion.Format(b.config.BotName, b.config.BotVersion))

	b.config.Validate()

	return b, nil
}

func (b *SBI) SetSysEvent(sysf func(MsgIncoming)) error {
	if sysf != nil {
		b.sys = sysf
		b.logger.Infof(M_SYSF_UPDATED.String())
		return nil
	}
	b.logger.Errorf(M_SYSF_INVALID.String())
	return M_SYSF_INVALID
}

func (b *SBI) Run() error {
	{
		m, err := module.NewModules(b.logger, b.config.Module.Dir, b.config.Module.Conf)
		if err != nil {
			return err
		}
		b.modules = m
	}
	b.logger.Infof(MF_HTTP_SERVING_SAddr_SPort_SPath.Format(b.config.Service.Address, b.config.Service.Port, b.config.Service.Path))
	return b.serve()
}
