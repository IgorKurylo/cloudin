package main

import (
	"bufio"
	"cloudin/cmd"
	"cloudin/config"
	"cloudin/utils"
	"fmt"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	BuildVersion string = ""
	BuildTime    string = ""
)
var configFile = config.Config{}
var (
	configure        = kingpin.Command("configure", "Configuration for Cloud In")
	mfa              = kingpin.Command("mfa", "aws cli login with mfa code")
	mfaCode          = mfa.Flag("code", "MFA Code").Required().String()
	docker           = kingpin.Command("ecr", "login to ecr repository")
	cluster          = kingpin.Command("k8s", "choose EKS cluster to load configuration")
	updateKubeConfig = cluster.Flag("update", "set T/F (true/false) for update kube config or export as environment variable").Required().String()
)
var configured bool
var dir string

func main() {
	printVersionInfo()
	configured, dir = utils.ScanAWSDirectory()
	awsCliExists := utils.CommandExists("aws")
	if !awsCliExists {
		fmt.Printf("aws cli not found please install and try again")
		os.Exit(1)
	}

	kingpin.Version(BuildVersion)
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
		dockerCliExists := utils.CommandExists("docker")
		if dockerCliExists {
			initConfiguration()
			cmd.DockerLogin(configFile, false)
		} else {
			fmt.Printf("docker cli not found please install and try again")
		}
	case cluster.FullCommand():
		kubectlCliExists := utils.CommandExists("kubectl")
		if kubectlCliExists {
			initConfiguration()
			updateConfig := convertFlag(*updateKubeConfig)
			cmd.K8SLogin(configFile, updateConfig)

		} else {
			fmt.Printf("kubectl cli not found please install and try again")
		}
	}
}

func convertFlag(s string) bool {
	flag := strings.ToLower(s)
	if flag == "t" || flag == "true" {
		return true
	} else {
		return false
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
