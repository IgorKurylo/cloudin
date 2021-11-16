package utils

import (
	"encoding/json"
	"errors"
	"fmt"
)

type AWSCredentials struct {
	Result Credentials `json:"Credentials"`
}
type Credentials struct {
	AccessKeyId     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
	Expiration      string `json:"Expiration"`
}

func ParseCredentials(credentialsJson string) (AWSCredentials, error) {
	credentials := AWSCredentials{}
	fmt.Printf(credentialsJson)
	err := json.Unmarshal([]byte(credentialsJson), &credentials)
	if err != nil {

		return credentials, errors.New("error on credentials parsing")
	}
	return credentials, nil

}
