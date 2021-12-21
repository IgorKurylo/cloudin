package cmd

import (
	"cloudin/config"
	"cloudin/utils"
	"fmt"
)

var (
	dockerUsername = "AWS"
)

func DockerLogin(config config.Config, passwordOutput bool) {
	var stdOutResult string
	var stdErrResult string
	awsRegion := config.Cloud.Region
	ecrRepositoryUrl := config.Docker.RepoURL
	flags := []string{"ecr", "get-login-password", "--region", awsRegion}
	if passwordOutput {
		err := Executor("aws", &stdOutResult, &stdErrResult, flags...)
		if err != nil {
			utils.ProcessError(err)
		}
		if stdOutResult != "" {
			fmt.Println(fmt.Sprintf("ecr password %s", stdOutResult))
		} else {
			fmt.Println(stdErrResult)
		}
	} else {
		flags = append(flags, "|", "docker", "login", "--username", dockerUsername, "--password-stdin", ecrRepositoryUrl)
		err := Executor("aws", &stdOutResult, &stdErrResult, flags...)
		if err != nil {
			utils.ProcessError(err)
		}
		if stdOutResult != "" {
			fmt.Println(stdOutResult)
		} else {
			fmt.Println(stdErrResult)
		}
	}
}
