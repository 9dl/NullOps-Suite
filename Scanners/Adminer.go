package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func scanAdminer(config *Helpers.Runner) *Helpers.RunnerResult {
	Response, Error := Helpers.SendRequest(Helpers.ExtractDomain(config.Line)+"adminer.php", "GET", "", Helpers.RequestOptions{})

	if Error == nil {
		if strings.Contains(string(Response.Body), "<a href='https://www.adminer.org/' target=\"_blank\" rel=\"noreferrer noopener\" id='h1'>Adminer</a>") {
			Version, _ := Helpers.ParseLR(string(Response.Body), "<span class=\"version\">", "</span>")
			VersionCleaned := strings.ReplaceAll(Version, ".", "")
			VersionInt, _ := strconv.Atoi(VersionCleaned)
			VersionVulnerable := false
			if VersionInt < 463 {
				VersionVulnerable = true
				err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/Adminer (Vulnerable).txt", []string{
					Helpers.ExtractDomain(config.Line) + "adminer.php"},
				)
				CLI_Handlers.LogError(err)
			}

			err := CLI_Handlers.AppendToFile(Helpers.OutputPath+fmt.Sprintf("/Adminer (%v).txt", Version), []string{
				Helpers.ExtractDomain(config.Line) + "adminer.php"},
			)
			CLI_Handlers.LogError(err)

			return &Helpers.RunnerResult{
				Line:   fmt.Sprintf("%v | Version: %v | Vulnerable?: %v", Helpers.ExtractDomain(config.Line), Version, VersionVulnerable),
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

func ScannerAdminer(config *Helpers.ScanConfig) {
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
		RunnerResult := scanAdminer(&ScanConfig)

		mu.Lock()
		Helpers.Checked++
		mu.Unlock()

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
