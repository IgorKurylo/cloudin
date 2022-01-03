package unit_test

import (
	"cloudin/cmd"
	"cloudin/config"
	"testing"
)

var configFile = config.Config{}

func TestMFALogin(t *testing.T) {
	var code string
	code="11111"
	cmd.MFALogin(code, configFile)
}

