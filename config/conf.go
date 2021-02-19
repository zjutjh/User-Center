package config

import (
  "gopkg.in/yaml.v3"
  "io/ioutil"
  "log"
)

type Session struct {
  Name         string    `yaml:"name"`
  SecretKey    string    `yaml:"secret_key"`
}

type Client struct {
  ID          string     `yaml:"id"`
  Secret      string     `yaml:"secret"`
  Name        string     `yaml:"name"`
  Domain      string     `yaml:"domain"`
  Scope       []string   `yaml:"scope"`
}

type Db struct {
  Type        string     `yaml:"type"`
  Host        string     `yaml:"host"`
  Port        int        `yaml:"port"`
  User        string     `yaml:"user"`
  Password    string     `yaml:"password"`
  DbName      string     `yaml:"dbname"`
}

type Config struct {
  Port        string     `yaml:"port"`
  Session     Session    `yaml:"session"`
  Db struct {
    Auth      Db         `yaml:"auth"`
    Resource  Db         `yaml:"resource"`
  }                      `yaml:"db"`
  Client      []Client   `yaml:"client"`
}

var config *Config

func init() {
  configPath := "config.yaml"
  configBytes, err := ioutil.ReadFile(configPath)
  if err != nil {
    log.Fatalln(err)
  }
  err = yaml.Unmarshal(configBytes, &config)
  if err != nil {
    log.Fatalln(err)
  }
}

func GetConfig() *Config {
  return config
}
