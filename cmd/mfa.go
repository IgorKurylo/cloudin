package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"smart.login.aws/config"
	_ "smart.login.aws/config"
	"smart.login.aws/utils"
)

var (
	arnTpl       = "arn:aws:iam::%s:mfa/%s"
	AccessKey    = "aws_access_key_id"
	SecretKey    = "aws_secret_access_key"
	SessionToken = "aws_session_token"
	FileName     = "stsinput.json"
)

type STSProfile struct {
	SecretKey    string
	AccessKey    string
	SessionToken string
	Profile      string
	Region       string
	Size         struct {
		accessKey  int
		secretKey  int
		sessionKey int
	}
}
type StsInput struct {
	DurationSeconds int32  `json:"DurationSeconds"`
	SerialNumber    string `json:"SerialNumber"`
	TokenCode       string `json:"TokenCode"`
}

var commandsMap map[string][]string

func MFALogin(mfaCode string, config config.Config) {
	account := config.Cloud.AWSAccount
	user := config.Cloud.AWSUser

	arn := fmt.Sprintf(arnTpl, account, user)
	awsProfile := config.Cloud.Profile
	var stdOutResult string
	var stdErrResult string
	flags := []string{"sts", "get-session-token", "--cli-input-json", "--profile", awsProfile, fmt.Sprintf("file://%s", FileName)}
	input := StsInput{
		DurationSeconds: config.Cloud.STSDuration,
		SerialNumber:    arn,
		TokenCode:       mfaCode,
	}
	createTmpInputFile(input, FileName)
	err := Executor("aws", &stdOutResult, &stdErrResult, flags...)
	if err != nil {
		utils.ProcessError(err)
	}
	if stdOutResult != "" {
		credentials, err := utils.ParseCredentials(stdOutResult)
		if err != nil {
			utils.ProcessError(err)
		}
		exportingVariables(credentials, config.Cloud.StsProfile)
	}
}
func createTmpInputFile(input StsInput, fileName string) {
	file, _ := json.MarshalIndent(input, "", " ")
	_ = ioutil.WriteFile(fileName, file, 0644)
}
func exportingVariables(credentials utils.AWSCredentials, stsProfile string) {
	commandsMap = map[string][]string{
		AccessKey:    utils.BuildConfigureCredentialsCmd(AccessKey, credentials.Result.AccessKeyId, stsProfile),
		SecretKey:    utils.BuildConfigureCredentialsCmd(SecretKey, credentials.Result.SecretAccessKey, stsProfile),
		SessionToken: utils.BuildConfigureCredentialsCmd(SessionToken, credentials.Result.SessionToken, stsProfile),
	}

	var stdOutResult string
	var stdErrResult string
	for _, value := range commandsMap {
		err := Executor("aws", &stdOutResult, &stdErrResult, value...)
		if err != nil {
			utils.ProcessError(err)
		}
	}
}
