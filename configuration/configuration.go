package configuration

import (
  "crypto/tls"
  "fmt"
  "io/ioutil"
  "encoding/json"

  "github.com/layeh/gumble/gumble"
  "github.com/Sirupsen/logrus"
  "os"
)

type ButlerConfiguration struct {
  Log      struct {
             File  string `json:"file"`
             Level string `json:"level"`
           } `json:"log"`
  Bot      struct {
             Username string `json:"username"`
           } `json:"bot"`
  Server   struct {
             Host               string `json:"host"`
             Port               int `json:"port"`
             InsecureSkipVerify bool `json:"insecureSkipVerify"`
           } `json:"server"`
  Features struct {
           } `json:"features"`
}

var Logger *logrus.Logger

func (cfg *ButlerConfiguration)  GetLogger() *logrus.Logger {
  return Logger
}

// set 'Logger' to the desired logger.
// True for std logger False for logrus
func initLog(useStdLogger bool, config ButlerConfiguration) {
  logcfg := config.Log
  if !useStdLogger {
    Logger = logrus.New()
    //Logger.SetFormatter(&Logger.JSONFormatter{})

    //if not empty write to file (else STDout)
    if logcfg.File != "" {
      file, err := os.OpenFile(logcfg.File, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
      if err != nil {
        panic(err)
      }
      Logger.Out = file
    } else {
      Logger.Out = os.Stdout
    }
    // set log the warning severity
    switch logcfg.Level {
    case  "":
      Logger.Level = logrus.InfoLevel
      Logger.Info("No logging level set settting to Info")
    case "error":
      Logger.Level = logrus.ErrorLevel
    case "warning":
      Logger.Level = logrus.WarnLevel
    case "debug":
      Logger.Level = logrus.DebugLevel
    default:
      panic("unrecognized logger level")
    }
  }
}

// Loads config from file name
// returns config ButlerConfiguration struct + err
// also sets Config var
func LoadConfiguration() (ButlerConfiguration, error) {
  var configuration ButlerConfiguration

  file, readErr := ioutil.ReadFile("config.json")
  if readErr != nil {
    return configuration, readErr
  }

  fmt.Printf("%+v\n", configuration) //last print before we get the logger going
  json.Unmarshal(file, &configuration)

  initLog(false, configuration)

  return configuration, nil
}

func (cfg *ButlerConfiguration) ExplodeConfiguration() (tls.Config, gumble.Config) {

  tlsConfig := tls.Config{}
  tlsConfig.InsecureSkipVerify = true

  Logger.Info(tlsConfig)

  gumbleConfig := gumble.Config{}

  if len(cfg.Bot.Username) > 0 {
    gumbleConfig.Username = cfg.Bot.Username
  } else {
    gumbleConfig.Username = "gumble-test"
  }

  return tlsConfig, gumbleConfig
}

func (cfg *ButlerConfiguration) GetUri() string {
  uri := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
  return uri
}
