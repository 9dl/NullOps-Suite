package Helpers

import (
	"NullOps/CLI_Handlers"
	"NullOps/Interface"
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ExtractDomain(inputURL string) string {
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "http://" + inputURL
	}

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return inputURL
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host + "/"

	return baseURL
}

func ExtractDomain2(inputURL string) string {
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "http://" + inputURL
	}

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return inputURL
	}

	baseURL := parsedURL.Scheme + "://cpanel." + parsedURL.Host + "/"

	return baseURL
}

func ExtractHost(inputURL string) string {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return inputURL
	}

	return parsedURL.Host
}

func BytesToString(input interface{}) string {
	switch input := input.(type) {
	case []byte:
		return string(input)
	case string:
		return input
	default:
		return fmt.Sprintf("Unsupported input type: %T", input)
	}
}

func CalculateCPM(valid int, elapsedTime time.Duration) float64 {
	elapsedMinutes := elapsedTime.Minutes()
	cpm := float64(valid) / elapsedMinutes
	return cpm
}

func BestCPM(now int, before int) int {

	if now > before {
		before = now
	}
	return before
}

func ShowResults() {
	fmt.Println()
	Interface.Write("FINAL RESULTS")
	Interface.Option("Valid", strconv.Itoa(int(Valid)))
	Interface.Option("Invalid", strconv.Itoa(int(Invalid)))
	Interface.Option("Checked", strconv.Itoa(int(Valid+Invalid)))
	Interface.Option("CPM", strconv.Itoa(int(CPM)))
	Interface.Option("Highest CPM", strconv.Itoa(int(HighestCPM)))
}

func ParseLR(str, before, after string) (string, error) {
	idx := strings.Index(str, before)
	if idx == -1 {
		return "", fmt.Errorf("substring '%s' not found in input", before)
	}

	start := idx + len(before)
	end := strings.Index(str[start:], after)
	if end == -1 {
		return "", fmt.Errorf("substring '%s' not found after '%s'", after, before)
	}

	return str[start : start+end], nil
}

func CleanTableNames(tables []string) []string {
	cleanedTables := make([]string, 0)
	keywords := []string{
		"email", "emails", "mail", "mails", "user", "users", "username", "usernames",
		"pseudo", "user", "users", "member", "members", "customer",
		"customers", "login", "signin", "password", "passwords", "pass", "pwd", "pw", "pws", "passwort",
	}

	for _, tableName := range tables {
		isKeyword := false
		for _, keyword := range keywords {
			if tableName == keyword {
				isKeyword = true
				break
			}
		}

		if isKeyword {
			cleanedTables = append(cleanedTables, tableName)
		}
	}

	return cleanedTables
}

func IsPythonInstalled() bool {
	_, err := exec.LookPath("python")
	if err == nil {
		return true
	}

	_, err = exec.LookPath("python3")
	return err == nil
}

func DirectoryExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func DownloadAndExtractFile(url, folderPath string) error {
	if err := os.MkdirAll(SanitizeFile(folderPath), 0750); err != nil {
		return err
	}

	resp, err := http.Get(SanitizeURL(url))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	filePath := filepath.Join(folderPath, "sqlmap.zip")
	file, err := os.Create(SanitizeFile(filePath))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	zipFile, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	for _, file := range zipFile.File {
		targetPath := filepath.Join(SanitizeFile(folderPath), SanitizeFile(file.Name))
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(targetPath, 0750)
			CLI_Handlers.LogError(err)
			continue
		}

		targetFile, err := os.Create(SanitizeFile(targetPath))
		if err != nil {
			return err
		}
		defer targetFile.Close()

		sourceFile, err := file.Open()
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		_, err = io.CopyN(targetFile, sourceFile, 30000*1024*1024) // 30 GB
		if err != nil {
			return err
		}
	}

	return nil
}
