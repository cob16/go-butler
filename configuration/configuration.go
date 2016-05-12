package configuration

import (
  "crypto/tls"
  "fmt"
  "io/ioutil"
  "encoding/json"

  "github.com/layeh/gumble/gumble"
)

type ButlerConfiguration struct {
  Bot struct {
    Username string `json:"username"`
  } `json:"bot"`
  Server struct {
    Host string `json:"host"`
    Port int `json:"port"`
    InsecureSkipVerify bool `json:"insecureSkipVerify"`
  } `json:"server"`
  Features struct {
  } `json:"features"`
}

func LoadConfiguration() (ButlerConfiguration, error) {
  var configuration ButlerConfiguration

  file, readErr := ioutil.ReadFile("config.json")
  if readErr != nil {
    return configuration, readErr
  }

  fmt.Println(string(file))

  json.Unmarshal(file, &configuration)

  fmt.Printf("%+v\n", configuration)

  return configuration, nil
}

func (cfg *ButlerConfiguration) ExplodeConfiguration() (tls.Config, gumble.Config) {

  tlsConfig := tls.Config{}
  tlsConfig.InsecureSkipVerify = true

  fmt.Printf("%+v\n", tlsConfig)

  gumbleConfig := gumble.Config{}
  gumbleConfig.Username = "gumble-test"

  return tlsConfig, gumbleConfig
}

func (cfg *ButlerConfiguration) GetUri() string {
  uri := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
  fmt.Println(cfg.Server.Host)
  fmt.Println(cfg.Server.Port)
  return uri
}
