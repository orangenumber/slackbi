package slackbi

import (
	"encoding/json"
	"github.com/gonyyi/aface"
	"io/ioutil"
	"sync"
)

type SBI struct {
	config         *config
	configFilename string
	logger         aface.Logger1a
	modules        Modules
	mu             sync.RWMutex
}

const (
	SBI_VERSION                = "0.0.6 (2020-12-31)"
	SLACK_ENDPOINT_MSG         = "https://slack.com/api/chat.postMessage"
	SLACK_ENDPOINT_FILE_UPLOAD = "https://slack.com/api/files.upload"
)

func New(configFile string, logger aface.Logger1a) (*SBI, error) {
	if configFile == "" {
		return nil, m_config_empty
	}
	b := &SBI{}

	// READ CONFIGURATION
	if cb, err := ioutil.ReadFile(configFile); err != nil {
		return nil, mf_config_read_err_sErr.Errorf(err.Error())
	} else {
		c := config{}
		if err := json.Unmarshal(cb, &c); err != nil {
			return nil, err
		}
		b.configFilename = configFile
		b.config = &c
	}

	// VALIDATE CONFIG
	b.config.Validate()

	// IF NO LOGGER IS SET, USE A DUMMY
	if logger != nil {
		b.logger = logger
	} else {
		b.logger = &aface.LoggerDummy1a{} // to prevent null ptr error
	}

	// Add a parent SBI pointer to modules.
	b.modules.p = b
	if err := b.modules.Load(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *SBI) Run() error {
	// Log version information
	b.logger.Infof(mf_sbi_sVersion.Format(SBI_VERSION))
	b.logger.Infof(mf_app_sName_sVersion.Format(b.config.BotName, b.config.BotVersion))

	// Log starting info
	b.logger.Infof(
		mf_http_serving_sAddr_sPort_sPath.Format(
			b.config.Service.Address,
			b.config.Service.Port,
			b.config.Service.Path))
	return b.serve()
}
