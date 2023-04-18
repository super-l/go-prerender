package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type SystemConfig struct {
	System struct {
		LogLevel  string `yaml:"logLevel"`
		BaseUrl   string `yaml:"baseUrl"`
		H5        bool   `yaml:"h5"`
		CacheTime int64  `yaml:"cacheTime"`
	} `yaml:"system"`
	Browser struct {
		Timeout int `yaml:"timeout"`
		Show    int `yaml:"show"`
	} `yaml:"browser"`
	HttpServ struct {
		Ip   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"httpServ"`
}

var sysConfig SystemConfig

func init() {
	conf, err := ReadYamlConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
		return
	}
	sysConfig = conf
	return
}

func GetConfig() SystemConfig {
	return sysConfig
}

func ReadYamlConfig(path string) (config SystemConfig, err error) {
	conf := SystemConfig{}
	if f, err := os.Open(path); err != nil {
		return config, err
	} else {
		yaml.NewDecoder(f).Decode(&conf)
	}
	return conf, nil
}

func WriteConfig(path string, config SystemConfig) error {
	data, err := yaml.Marshal(config)
	err = ioutil.WriteFile(path, data, 0777)
	return err
}
