package main

import (
	"bufio"
	"cloudin/cmd"
	"cloudin/config"
	"cloudin/utils"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

var (
	BuildVersion string = ""
	BuildTime    string = ""
)
var configFile = config.Config{}
var (
	configure = kingpin.Command("configure", "Configuration for Cloud In")
	mfa       = kingpin.Command("mfa", "aws cli login with mfa code")
	mfaCode   = mfa.Flag("code", "MFA Code").Required().String()
	docker    = kingpin.Command("ecr-login", "aws cli login to ecr repository")

	//cluster     = kingpin.Command("setup-cluster", "choose EKS cluster to load configuration")
	//clusterArg  = cluster.Arg("name", "name of eks cluster").Required().String()
	//regionFlag  = docker.Arg("region", "region on aws ecr repository").Required().String()
	//repoUrlFlag = docker.Arg("url", "url on aws of ecr repository").Required().String()
)
var configured bool
var dir string

func main() {
	printVersionInfo()
	configured, dir = utils.ScanAWSDirectory()
	awsCliExists := utils.CommandExists("aws")
	dockerCliExists := utils.CommandExists("docker")
	if !awsCliExists && !dockerCliExists {
		fmt.Printf("aws & docker cli not found, please install and try again")
		os.Exit(1)
	}

	kingpin.Version("0.0.1")
	commandsParsing()

}
func initConfiguration() {
	config.LoadConfig(&configFile)
}
func commandsParsing() {
	if !configured {
		fmt.Printf(".aws directory not found in %s\n", dir)
		fmt.Print("please run the next command: 'aws configure'\n")
		os.Exit(1)
	}

	switch kingpin.Parse() {
	// login mfa code
	case configure.FullCommand():
		readConfiguration()
	case mfa.FullCommand():
		initConfiguration()
		cmd.MFALogin(*mfaCode, configFile)
	case docker.FullCommand():
		initConfiguration()
		cmd.DockerLogin(configFile, false)
	}
}

func readConfiguration() {

	reader := bufio.NewReader(os.Stdin)
	configSet := make(map[string]string)
	for i := 0; i < len(config.ConfigurationParams); i++ {
		fmt.Printf("%s: ", config.ConfigurationParams[i])
		result, _ := reader.ReadString('\n')
		configSet[config.ConfigurationParams[i]] = strings.TrimSuffix(result, "\n")
	}
	config.SaveConfig(configSet)

}
func printVersionInfo() {
	fmt.Println(fmt.Sprintf("CLI Application CloudIn\n"))
	fmt.Println(fmt.Sprintf("BuildVersion: %s\n", BuildVersion))
	fmt.Println(fmt.Sprintf("BuildTime: %s", BuildTime))

}
