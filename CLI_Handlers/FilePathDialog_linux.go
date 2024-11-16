//go:build linux
// +build linux

package CLI_Handlers

import (
	"bufio"
	"fmt"
	"os"
)

func GetFilePath() string {
	fmt.Println("Please enter the file path:")
	reader := bufio.NewReader(os.Stdin)
	FilePath, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	return FilePath[:len(FilePath)-1] // Trim the newline character
}
