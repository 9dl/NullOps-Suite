package CLI_Handlers

import (
	Interface "NullOps/Interface"
	"bufio"
	"fmt"
	"github.com/sqweek/dialog"
	"os"
	"strings"
	"sync"
)

var (
	fileMutex sync.Mutex
)

func GetFilePath() string {
	Interface.Clear()
	FilePath, err := dialog.File().Title("Select Input File").Load()
	if err != nil {
		fmt.Println(err)
	}

	return FilePath
}

func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func AppendToFile(filename string, lines []string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}