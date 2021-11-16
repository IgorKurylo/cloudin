package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func ScanAWSDirectory() (bool, string) {
	userHomeDir, err := os.UserHomeDir()
	var files []string
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Home dir is %s\n", userHomeDir)
	fileInfo, err := ioutil.ReadDir(userHomeDir)
	if err != nil {
		return false, userHomeDir
	}
	for _, file := range fileInfo {
		files = append(files, file.Name())

	}
	for i := 0; i < len(files); i++ {
		if files[i] == ".aws" {
			return true, userHomeDir
		}
	}
	return false, userHomeDir

}
func ProcessError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
func BuildConfigureCredentialsCmd(key string, value string, profile string) ([]string) {
	return []string{"configure", "set", key, value, "--profile", profile}

}
