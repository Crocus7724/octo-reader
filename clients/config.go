package clients

import (
  "bytes"
  "encoding/json"
  "errors"
  "flag"
  "io/ioutil"
  "strings"
)

const (
  ConfigFileName = "config.json"
)

var (
  Debug bool
)

type GithubConfig struct {
  Name     string `json:"name"`
  Password string `json:"password"`
  Token    string `json:"token"`
  Host     string `json:"host"`
}

func (g *GithubConfig) GetApiHost() (string, error) {
  if !strings.HasPrefix(g.Host, "https://") {
    return "", errors.New("ホストは\"https://\"プレフィックスをつけて下さい")
  }

  var b bytes.Buffer

  for i, s := range g.Host {
    b.WriteRune(s)

    if i == 7 {
      b.WriteString("api.")
    }
  }

  if !strings.HasSuffix(g.Host, "/") {
    b.WriteString("/")
  }

  return b.String(), nil
}

func CreateConfig() GithubConfig {
  config, _ := parseConfigFile()

  return *parseArgs(&config)
}

func parseConfigFile() (GithubConfig, error) {

  file, err := ioutil.ReadFile(ConfigFileName)

  if err != nil {
    return GithubConfig{}, err
  }

  var config GithubConfig
  err = json.Unmarshal(file, &config)

  if err != nil {
    return GithubConfig{}, err
  }

  return config, err
}

func parseArgs(config *GithubConfig) *GithubConfig {
  name := flag.String("name", "", "Github Login Name")
  password := flag.String("password", "", "Github Login Password")
  token := flag.String("token", "", "Github Personal Access Token")
  host := flag.String("host", "https://github.com/", "Github Host")
  debug := flag.Bool("debug", false, "Debug mode")

  flag.Parse()

  if *name != "" {
    config.Name = *name
  }

  if *password != "" {
    config.Password = *password
  }

  if *token != "" {
    config.Token = *token
  }

  if *host != "https://github.com/" {
    config.Host = *host
  }

  if *debug {
    Debug = *debug
  }

  return config
}
