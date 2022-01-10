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
	configure = kingpin.Command("configure", "Configuration for Cloud In")
	mfa       = kingpin.Command("mfa", "aws cli login with mfa code")
	mfaCode   = mfa.Flag("code", "MFA Code").Required().String()
	docker    = kingpin.Command("ecr", "login to ecr repository")
	cluster   = kingpin.Command("k8s", "login to EKS cluster")
	//clusterName = cluster.Flag("cluster-name", "eks cluster name").Required().String()
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
	reader := bufio.NewReader(os.Stdin)
	switch kingpin.Parse() {
	// login mfa code
	case configure.FullCommand():
		readConfiguration(*reader)
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
		var index = 0
		kubectlCliExists := utils.CommandExists("kubectl")
		if kubectlCliExists {
			initConfiguration()
			cmd.EKSClusters(configFile)
			if cmd.ClusterSize > 0 {
				fmt.Println("Choose cluster from next list")
				result, _ := fmt.Scanf("%d", &index)
				if result <= cmd.ClusterSize {
					clusterName, err := cmd.GetEKSClusterName(result)
					if err != nil {
						utils.ProcessError(err)
					}
					cmd.K8SLogin(configFile, clusterName)
				}
			} else {
				fmt.Println("No eks cluster found")
			}
		} else {
			fmt.Printf("kubectl cli not found please install and try again")
		}
	}
}

func readConfiguration(reader bufio.Reader) {

	configSet := make(map[string]string)
	for i := 0; i < len(config.ConfigurationParams); i++ {
		fmt.Printf("%s: ", config.ConfigurationParams[i])
		result, _ := reader.ReadString('\n')
		configSet[config.ConfigurationParams[i]] = strings.TrimSuffix(result, "\n")
	}
	config.SaveConfig(configSet)

}
func printVersionInfo() {
	fmt.Println(fmt.Sprintf("CLI Application CloudIn"))
	fmt.Println(fmt.Sprintf("BuildVersion: %s", BuildVersion))
	fmt.Println(fmt.Sprintf("BuildTime: %s", BuildTime))
}
