package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"smart.login.aws/cmd"
	"smart.login.aws/config"
	"smart.login.aws/utils"
)

var configFile = config.Config{}
var (
	mfa         = kingpin.Command("mfa", "aws cli login with mfa code")
	mfaCode     = mfa.Flag("code", "MFA Code").Required().String()
	cluster     = kingpin.Command("setup-cluster", "choose EKS cluster to load configuration")
	clusterArg  = cluster.Arg("name", "name of eks cluster").Required().String()
	docker      = kingpin.Command("ecr-login", "login with docker cli to ecr repository")
	regionFlag  = docker.Flag("region", "region on aws ecr repository").Required().String()
	repoUrlFlag = docker.Flag("url", "url on aws of ecr repository").Required().String()
)

func main() {
	configured, dir := utils.ScanAWSDirectory()
	commandExists := utils.CommandExists("aws")
	if !commandExists {
		fmt.Printf("aws cli not found")
		os.Exit(1)
	}
	if configured {
		initConfiguration()
	} else {
		fmt.Printf(".aws directory not found in %s", dir)
		fmt.Print("please run the next command: 'aws configure' ") //TODO: ask from user set the configuration creds and config this.
		os.Exit(1)
	}
	kingpin.Version("0.0.1")
	commandsParsing()

}
func initConfiguration() {

	config.ReadConfig(&configFile)

}
func commandsParsing() {
	switch kingpin.Parse() {
	// login mfa code
	case "mfa":
		cmd.MFALogin(*mfaCode, configFile)
	}
}
