package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"strings"
	"time"
)

func scanCPanel(config *Helpers.Runner) *Helpers.RunnerResult {
	Response, Error := Helpers.SendRequest(Helpers.ExtractDomain2(config.Line), "GET", "", Helpers.RequestOptions{})

	if Error == nil {
		if strings.Contains(string(Response.Body), "<div class=\"locale-cell\"><a href=\"?locale=i_cpanel_snowmen\">☃ cPanel Snowmen ☃ - i_cpanel_snowmen</a></div>") {

			err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/cPanel.txt", []string{
				Helpers.ExtractDomain2(config.Line),
			},
			)

			CLI_Handlers.LogError(err)

			return &Helpers.RunnerResult{
				Line:   Helpers.ExtractDomain2(config.Line),
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

func ScannerCPanel(config *Helpers.ScanConfig) {
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
		RunnerResult := scanCPanel(&ScanConfig)

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
