package config

import (
	"cloudin/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	exportProfile       = "sts"
	tokenDuration       = 129600
	directoryConfigName = ".cloudin"
	fileConfigName      = "config.yaml"
	repoUrlTemplate     = "%s.dkr.ecr.%s.amazonaws.com"
	ConfigurationParams = []string{"aws_account", "aws_user", "profile", "region"}
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

func buildConfig(configSet map[string]string, config *Config) {
	config.Cloud.StsProfile = exportProfile
	config.Cloud.AWSAccount = configSet[ConfigurationParams[0]]
	config.Cloud.AWSUser = configSet[ConfigurationParams[1]]
	config.Cloud.Profile = configSet[ConfigurationParams[2]]
	config.Cloud.Region = configSet[ConfigurationParams[3]]
	config.Cloud.STSDuration = int32(tokenDuration)
	config.Docker.RepoURL = fmt.Sprintf(repoUrlTemplate, config.Cloud.AWSAccount, config.Cloud.Region)
}
func LoadConfig(config *Config) {
	configPathDir := configPath()
	configFile, err := os.Open(configPathDir)
	if err != nil {
		fmt.Printf("Config file not exists, please run configure command")
		utils.ProcessError(err)
	}
	defer configFile.Close()
	load(configFile, config)
	fmt.Printf("Configuration loaded %s\n", configFile.Name())
}
func SaveConfig(configSet map[string]string) {
	config := &Config{}
	configPathDir := configPath()
	_, err := os.Stat(configPathDir)

	if err != nil {
		err := os.Mkdir(configPathDir, 0755)
		if err != nil {
			log.Fatalf("Error occurred on configuration %v", err)
		}
		buildConfig(configSet, config)
	} else {
		buildConfig(configSet, config)

	}
	saveToFile(config, fmt.Sprintf("%s/%s", configPathDir, fileConfigName))
}

func configPath() string {
	homeDir, _ := os.UserHomeDir()
	configPathDir := fmt.Sprintf("%s/%s/%s", homeDir, directoryConfigName, fileConfigName)
	return filepath.FromSlash(configPathDir)
}
func saveToFile(config *Config, pathFile string) {
	configData, err := yaml.Marshal(config)
	err = ioutil.WriteFile(pathFile, configData, 0644)
	if err != nil {
		utils.ProcessError(err)
	}
	fmt.Printf("Configuration saved %s", pathFile)
}
func load(configFilePath *os.File, config *Config) {
	decoder := yaml.NewDecoder(configFilePath)
	err := decoder.Decode(&config)
	if err != nil {
		utils.ProcessError(err)
	}
}
