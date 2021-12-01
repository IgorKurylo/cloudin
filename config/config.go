package config

import (
	"cloudin/utils"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Cloud struct {
		STSDuration int32  `yaml:"token_duration"`
		AWSAccount  string `yaml:"aws_account"`
		AWSUser     string `yaml:"aws_user"`
		Profile     string `yaml:"profile"`
		StsProfile  string `yaml:"sts_profile"`
		Region      string `yaml:"aws_region"`
	} `yaml:"cloud"`
	Docker struct {
		RepoURL string `yaml:"repository_url"`
	} `yaml:"docker"`
}

func ReadConfig(config *Config, profile string) {
	configFile, err := os.Open(fmt.Sprintf("config_%s.yaml", profile))
	if err != nil {
		utils.ProcessError(err)
	}
	defer configFile.Close()
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		utils.ProcessError(err)
	}
	fmt.Printf("Configuration loaded %s", configFile.Name())

}
