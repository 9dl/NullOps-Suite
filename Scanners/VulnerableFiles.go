package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"strings"
	"time"
)

var filetypes_vuln = []string{
	".zip",
	".rar",
	".7z",
	".tar",
}

var files_vulnerablePaths = []string{
	"backup.zip",
	"backup.rar",
	"backup.7z",
	"backup.tar",

	"pass.zip",
	"pass.rar",
	"pass.7z",
	"pass.tar",
	"pass.db",
	"pass.sql",

	"users.7z",
	"users.tar",
	"users.zip",
	"users.rar",
	"users.db",
	"users.sql",
}

func scanVulnFiles(config *Helpers.Runner) *Helpers.RunnerResult {
	host := Helpers.ExtractHost(config.Line)

	var siteBackups []string
	for _, ext := range filetypes_vuln {
		siteBackups = append(siteBackups,
			fmt.Sprintf("%v%v", host, ext),
			fmt.Sprintf("%v%v", strings.ToLower(host), ext),
			fmt.Sprintf("%v%v", strings.ToUpper(host), ext),
		)
	}

	for _, payload := range files_vulnerablePaths {
		_, Error := Helpers.SendRequest(Helpers.ExtractDomain(config.Line)+payload, "GET", "", Helpers.RequestOptions{})

		if Error != nil {
			err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/Files.txt", []string{
				Helpers.ExtractDomain(config.Line) + payload,
			},
			)

			CLI_Handlers.LogError(err)

			return &Helpers.RunnerResult{
				Line:   Helpers.ExtractDomain(config.Line),
				Status: true,
				Error:  nil,
			}
		}
	}

	for _, payload := range siteBackups {
		_, Error := Helpers.SendRequest(Helpers.ExtractDomain(config.Line)+payload, "GET", "", Helpers.RequestOptions{})

		if Error != nil {
			err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/Files.txt", []string{
				Helpers.ExtractDomain(config.Line) + payload,
			},
			)

			CLI_Handlers.LogError(err)

			return &Helpers.RunnerResult{
				Line:   Helpers.ExtractDomain(config.Line),
				Status: true,
				Error:  nil,
			}
		}
	}

	return &Helpers.RunnerResult{
		Line:   config.Line,
		Status: false,
		Error:  fmt.Errorf("Invalid"),
	}
}

func ScannerVulnFiles(config *Helpers.ScanConfig) {
	Helpers.Valid, Helpers.Invalid, Helpers.Checked, Helpers.PayloadsTested, Helpers.CPM, Helpers.HighestCPM, Helpers.Running = 0, 0, 0, 0, 0, 0, true

	FilePath := CLI_Handlers.GetFilePath()
	lines, err := CLI_Handlers.ReadLines(FilePath)
	CLI_Handlers.LogError(err)

	startTime := time.Now()
	go func() {
		for {
			if Helpers.Running {
				mu.Lock()
				elapsedTime := time.Since(startTime)
				valid := Helpers.Valid + Helpers.Invalid
				Helpers.CPM = int(Helpers.CalculateCPM(valid, elapsedTime))
				Helpers.HighestCPM = Helpers.BestCPM(Helpers.CPM, Helpers.HighestCPM)
				mu.Unlock()

				Interface.StatsTitle(fmt.Sprintf("NullOps | CPM: %v | Highest CPM: %v", Helpers.CPM, Helpers.HighestCPM), Helpers.Valid, Helpers.Invalid, Helpers.Checked, len(lines))
				time.Sleep(1 * time.Second)
			} else {
				return
			}
		}
	}()

	defer func() {
		Helpers.Running = false
	}()

	Helpers.Threading(func(s string) {
		ScanConfig := Helpers.Runner{Line: s}
		RunnerResult := scanVulnFiles(&ScanConfig)

		if RunnerResult.Error == nil {
			mu.Lock()
			Helpers.Valid++
			mu.Unlock()
			Interface.Option(config.Name, fmt.Sprintf("%v | Status: %v", RunnerResult.Line, RunnerResult.Status))
		} else {
			mu.Lock()
			Helpers.Invalid++
			mu.Unlock()
			if config.PrintInvalid {
				Interface.Option(config.Name, fmt.Sprintf("%v | Status: %v | Reason: %v", Helpers.ExtractDomain(RunnerResult.Line), RunnerResult.Status, RunnerResult.Error))
			}
		}
		mu.Lock()
		Helpers.Checked++
		mu.Unlock()
	}, config.Threads, lines)
}
