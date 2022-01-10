package utils

import "fmt"

func HandleErrorStdErr(stdError string) {
	fmt.Println(fmt.Sprintf("Error result message: %s", stdError))
}
