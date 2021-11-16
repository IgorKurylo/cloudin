package unit_test

import (
	"smart.login.aws/cmd"
	"smart.login.aws/config"
	"testing"
)

var configFile = config.Config{}

func TestMFALogin(t *testing.T) {
	initConfiguration()
	var code string
	code="11111"
	cmd.MFALogin(code, configFile)

}
func initConfiguration() {
	config.ReadConfig(&configFile)

}
