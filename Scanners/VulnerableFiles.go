package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"strings"
	"sync/atomic"
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
	var siteBackup = []string{
		fmt.Sprintf("%v.rar", Helpers.ExtractHost(config.Line)),
		fmt.Sprintf("%v.zip", Helpers.ExtractHost(config.Line)),
		fmt.Sprintf("%v.7z", Helpers.ExtractHost(config.Line)),
		fmt.Sprintf("%v.tar", Helpers.ExtractHost(config.Line)),

		fmt.Sprintf("%v.rar", strings.ToLower(Helpers.ExtractHost(config.Line))),
		fmt.Sprintf("%v.zip", strings.ToLower(Helpers.ExtractHost(config.Line))),
		fmt.Sprintf("%v.7z", strings.ToLower(Helpers.ExtractHost(config.Line))),
		fmt.Sprintf("%v.tar", strings.ToLower(Helpers.ExtractHost(config.Line))),

		fmt.Sprintf("%v.rar", strings.ToUpper(Helpers.ExtractHost(config.Line))),
		fmt.Sprintf("%v.zip", strings.ToUpper(Helpers.ExtractHost(config.Line))),
		fmt.Sprintf("%v.7z", strings.ToUpper(Helpers.ExtractHost(config.Line))),
		fmt.Sprintf("%v.tar", strings.ToUpper(Helpers.ExtractHost(config.Line))),
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

	for _, payload := range siteBackup {
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
				elapsedTime := time.Since(startTime)
				Helpers.CPM = int32(int(Helpers.CalculateCPM(int(atomic.LoadInt32(&Helpers.Valid))+int(atomic.LoadInt32(&Helpers.Invalid)), elapsedTime)))
				Helpers.HighestCPM = int32(Helpers.BestCPM(int(Helpers.CPM), int(atomic.LoadInt32(&Helpers.HighestCPM))))

				Interface.StatsTitle(fmt.Sprintf("NullOps | CPM: %v | Highest CPM: %v", int(atomic.LoadInt32(&Helpers.CPM)), int(atomic.LoadInt32(&Helpers.HighestCPM))), int(atomic.LoadInt32(&Helpers.Valid)), int(atomic.LoadInt32(&Helpers.Invalid)), int(atomic.LoadInt32(&Helpers.Checked)), len(lines))
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
			atomic.AddInt32(&Helpers.Valid, 1)
			Interface.Option(config.Name, fmt.Sprintf("%v | Status: %v", RunnerResult.Line, RunnerResult.Status))
		} else {
			atomic.AddInt32(&Helpers.Invalid, 1)
			if config.PrintInvalid {
				Interface.Option(config.Name, fmt.Sprintf("%v | Status: %v | Reason: %v", Helpers.ExtractDomain(RunnerResult.Line), RunnerResult.Status, RunnerResult.Error))
			}
		}
		atomic.AddInt32(&Helpers.Checked, 1)
	}, config.Threads, lines)
}
