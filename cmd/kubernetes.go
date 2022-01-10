package cmd

import (
	"cloudin/config"
	"cloudin/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var (
	kubeDirectory  = ".kube"
	ClusterSize    = 0
	ClusterEKSList []string
)

func K8SLogin(config config.Config, clusterName string) {

	var stdOutResult string
	var stdErrResult string
	awsRegion := config.Cloud.Region
	homeDir, err := os.UserHomeDir()
	if err != nil {
		utils.ProcessError(err)
	}
	kubeConfigPath := fmt.Sprintf("%s/%s/%s", homeDir, kubeDirectory, clusterName)
	flags := []string{"eks", "--region", awsRegion, "update-kubeconfig", "--name", clusterName, "--kubeconfig", kubeConfigPath, "--profile", config.Cloud.StsProfile}
	err = Executor("aws", &stdOutResult, &stdErrResult, flags...)
	if err != nil {
		utils.HandleErrorStdErr(stdErrResult)
	}
	if stdOutResult != "" {
		fmt.Printf("eks update kubeconfig output %s\n", stdOutResult)
	} else {
		fmt.Println(stdErrResult)
	}
}
func clustersList(region string, profile string) []string {

	var stdOutResult string
	var stdErrResult string
	var listResult map[string][]string
	var clusters []string
	flags := []string{"eks", "list-clusters", "--region", region, "--profile", profile}
	err := Executor("aws", &stdOutResult, &stdErrResult, flags...)
	if err != nil {
		utils.HandleErrorStdErr(stdErrResult)
		utils.ProcessError(err)
	}
	if stdOutResult != "" {
		err := json.Unmarshal([]byte(stdOutResult), &listResult)
		if err != nil {
			utils.ProcessError(err)
		}
		clusters = listResult["clusters"]
	}
	return clusters
}
func EKSClusters(config config.Config) {
	ClusterEKSList = clustersList(config.Cloud.Region, config.Cloud.StsProfile)
	ClusterSize = len(ClusterEKSList)
	for index, name := range ClusterEKSList {
		fmt.Printf("%d. %s\n", index+1, name)
	}
}
func GetEKSClusterName(index int) (string, error) {
	if index > 0 {
		return ClusterEKSList[index-1], nil
	}
	return "", errors.New("index out of range")
}
