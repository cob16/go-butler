package configuration

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"layeh.com/gumble/gumble"
	"os"
)

type ButlerConfiguration struct {
	Log struct {
		File  string `json:"file"`
		Level string `json:"level"`
	} `json:"log"`
	Bot struct {
		Username                 string   `json:"username"`
		RecursiveChannelMessages bool     `json:"recursiveChannelMessages"`
		DefaultChannel           string   `json:"defaultChannel"`
		AccessTokens             []string `json:"accessTokens"`
	} `json:"bot"`
	Server struct {
		Host               string `json:"host"`
		Port               int    `json:"port"`
		InsecureSkipVerify bool   `json:"insecureSkipVerify"`
	} `json:"server"`
	Greeter struct {
		WelcomeUsers             bool `json:"welcomeUsers"`
		PassConnectOnChannelJoin bool `json:"passConnectOnChannelJoin"`
	} `json:"greeter"`
}

var Logger *logrus.Logger

func (cfg *ButlerConfiguration) GetLogger() *logrus.Logger {
	return Logger
}

func initLog(File string, Level string) {
	Logger = logrus.New()
	//Logger.SetFormatter(&Logger.JSONFormatter{})

	//if not empty write to file (else STDout)
	if File != "" {
		file, err := os.OpenFile(File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		Logger.Out = file
	} else {
		Logger.Out = os.Stdout
	}
	// set log the warning severity
	switch Level {
	case "", "info":
		Logger.Info("No logging level set settting to Info")
		Logger.Level = logrus.InfoLevel
	case "error":
		Logger.Level = logrus.ErrorLevel
	case "warning":
		Logger.Level = logrus.WarnLevel
	case "debug":
		Logger.Level = logrus.DebugLevel
	default:
		panic("Unrecognized logger level")
	}
}

// Loads config from file name
// returns config ButlerConfiguration struct + err
// also sets Config var
func LoadConfiguration(configurationPath string) (ButlerConfiguration, error) {
	if configurationPath == "" {
		configurationPath = "config.json"
	}
	var configuration ButlerConfiguration

	file, readErr := ioutil.ReadFile(configurationPath)
	if readErr != nil {
		return configuration, readErr
	}

	json.Unmarshal(file, &configuration)

	initLog(configuration.Log.File, configuration.Log.Level)
	Logger.Debug("Logger started")

	return configuration, nil
}

func (cfg *ButlerConfiguration) GetGumbleConfig() (tls.Config, gumble.Config) {

	tlsConfig := tls.Config{}
	tlsConfig.InsecureSkipVerify = cfg.Server.InsecureSkipVerify

	gumbleConfig := gumble.Config{}

	if len(cfg.Bot.Username) > 0 {
		gumbleConfig.Username = cfg.Bot.Username
	} else {
		Logger.Warn("No user set falling back to user GoButler")
		gumbleConfig.Username = "GoButler"
	}
	gumbleConfig.Tokens = cfg.Bot.AccessTokens

	return tlsConfig, gumbleConfig
}

func (cfg *ButlerConfiguration) GetUri() string {
	uri := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	return uri
}
