package archiver

import (
	"io/ioutil"
	"os"

	"github.com/slack-go/slack"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

// Config - config for archiver
type Config struct {
	App        appInfo
	Exceptions exceptionsInfo `yaml:"exceptions"`
	Log        zerolog.Logger
	Slack      *slack.Client
}

type appInfo struct {
	UserToken string `yaml:"usertoken"`
	Stack     string `yaml:"stack"`
}

type exceptionsInfo struct {
	Channels  []string `yaml:"channels"`
	BlackList []string `yaml:"blacklist"`
}

// ReadConfigs - reads in the yml file in the secrets folder
func ReadConfigs() Config {

	var c Config

	// Get Logger
	c.Log = returnLogger()

	bytes, err := ioutil.ReadFile(secretsPath)
	if err != nil {
		c.Log.Error().Msgf("Failed to read the app config file %s - %s", secretsPath, err)
		os.Exit(0)
	}

	err = yaml.UnmarshalStrict(bytes, &c.App)
	if err != nil {
		c.Log.Error().Msgf("Failed to parse the app config file - %s", err)
		os.Exit(0)
	}

	bytes, err = ioutil.ReadFile(exceptionsPath)
	if err != nil {
		c.Log.Error().Msgf("Failed to read the exceptions config file %s - %s", exceptionsPath, err)
		os.Exit(0)
	}

	err = yaml.UnmarshalStrict(bytes, &c.Exceptions)
	if err != nil {
		c.Log.Error().Msgf("Failed to parse the exceptions config file - %s", err)
		os.Exit(0)
	}

	// Create Slack client
	c.Slack = slack.New(c.App.UserToken)

	return c
}

// ReturnLogger - returns the logger for the config to use
func returnLogger() zerolog.Logger {

	if os.Getenv("LOGLEVEL") == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
