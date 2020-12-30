package slackbi

const (
	DEFAULT_BOT_NAME                 = "sbiTest"
	DEFAULT_BOT_VERSION              = "0.0.1 (20xx-xx-xx)"
	DEFAULT_LOGGING                  = true
	DEFAULT_SERVICE_PORT             = ":8080"
	DEFAULT_SERVICE_PATH             = "/bot"
	DEFAULT_SERVICE_ACCEPT_CHALLENGE = true
	DEFAULT_MODULE_DIR               = "./modules"
	DEFAULT_MODULE_CONF              = "config.json"
)

func (c *config) Default() {
	c.BotVersion = DEFAULT_BOT_VERSION
	c.BotName = DEFAULT_BOT_NAME
	c.Logging = DEFAULT_LOGGING
	c.Service.Port = DEFAULT_SERVICE_PORT
	c.Service.Path = DEFAULT_SERVICE_PATH
	c.Service.AcceptChallege = DEFAULT_SERVICE_ACCEPT_CHALLENGE
	c.Module.Dir = DEFAULT_MODULE_DIR
	c.Module.Conf = DEFAULT_MODULE_CONF
}

func (c *config) Validate() {
	val := func(name, val, defaultVal string) string {
		if val == "" {
			c.logger.Warnf(MF_CONF_OVERRIDE_SName_SVal_SNewVal.Format(name, val, defaultVal))
			return defaultVal
		}
		return val
	}

	c.BotVersion = val("BotVersion", c.BotVersion, DEFAULT_BOT_VERSION)
	c.BotName = val("BotName", c.BotName, DEFAULT_BOT_NAME)
	c.Service.Port = val("Service.Port", c.Service.Port, DEFAULT_SERVICE_PORT)
	c.Service.Path = val("Service.Path", c.Service.Path, DEFAULT_SERVICE_PATH)
	c.Module.Dir = val("Module.Dir", c.Module.Dir, DEFAULT_MODULE_DIR)
	c.Module.Conf = val("Module.Conf", c.Module.Conf, DEFAULT_MODULE_CONF)
}
