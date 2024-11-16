package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

func scanLaravel(config *Helpers.Runner) *Helpers.RunnerResult {
	response, err := Helpers.SendRequest(Helpers.ExtractDomain(config.Line)+".env", "GET", "", Helpers.RequestOptions{LoggingEnabled: false})
	if err != nil {
		return &Helpers.RunnerResult{
			Line:   config.Line,
			Status: false,
			Error:  fmt.Errorf("Failed to send request: %v", err),
		}
	}

	if response.StatusCode != http.StatusOK {
		return &Helpers.RunnerResult{
			Line:   config.Line,
			Status: false,
			Error:  fmt.Errorf("Received non-OK status code: %d", response.StatusCode),
		}
	}

	if err == nil {
		responeStr := Helpers.BytesToString(response.Body)

		re_laravel := regexp.MustCompile(`(?i)APP_NAME=Laravel`)
		re_normal := regexp.MustCompile(`(?i)APP_NAME=`)
		if re_laravel.MatchString(responeStr) {
			_, capture := captureEnvironmentInfo(responeStr, Helpers.ExtractDomain(config.Line))
			return &Helpers.RunnerResult{
				Line:   Helpers.ExtractDomain(config.Line) + fmt.Sprintf(" | Capture: [%v]", capture),
				Status: true,
				Error:  nil,
			}
		} else if re_normal.MatchString(responeStr) {
			return &Helpers.RunnerResult{
				Line:   Helpers.ExtractDomain(config.Line) + fmt.Sprintf(" | Not Laravel."),
				Status: true,
				Error:  nil,
			}
		}
	}

	return &Helpers.RunnerResult{
		Line:   config.Line,
		Status: false,
		Error:  fmt.Errorf("No relevant data found"),
	}
}

func ScannerLaravel(config *Helpers.ScanConfig) {
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
		Helpers.ShowResults()
	}()

	Helpers.Threading(func(s string) {
		ScanConfig := Helpers.Runner{Line: s}
		RunnerResult := scanLaravel(&ScanConfig)

		mu.Lock()
		Helpers.Checked++
		mu.Unlock()
		if RunnerResult.Error == nil {
			mu.Lock()
			Helpers.Valid++
			mu.Unlock()
			Interface.Option(config.Name, RunnerResult.Line)
		} else {
			mu.Lock()
			Helpers.Invalid++
			mu.Unlock()
			if config.PrintInvalid {
				Interface.Option(config.Name, fmt.Sprintf("%v | Status: %v | Reason: %v", Helpers.ExtractDomain(RunnerResult.Line), RunnerResult.Status, RunnerResult.Error))
			}
		}
	}, config.Threads, lines)
}
