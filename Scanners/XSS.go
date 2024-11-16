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

var payloadXSS = []string{
	"</script><script>alert(1)</script>",
	"<script>alert(123);</script>",
	"<ScRipT>alert(\"XSS\");</ScRipT>",
	"<script>alert(123)</script>",
	"<script>alert( XSS )</script> ",
	"<script>alert( XSS );</script>",
}

func scanXSS(config *Helpers.Runner) *Helpers.RunnerResult {
	for _, payload := range payloadXSS {
		mu.Lock()
		Helpers.PayloadsTested++
		mu.Unlock()
		Response, Error := Helpers.SendRequest(config.Line, "GET", payload, Helpers.RequestOptions{})

		if Error == nil {
			if strings.ContainsAny(string(Response.Body), payload) {
				err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/XSS.txt", []string{config.Line + payload})
				CLI_Handlers.LogError(err)
				return &Helpers.RunnerResult{
					Line:   Helpers.ExtractDomain(config.Line) + " | Payload: " + payload,
					Status: true,
					Error:  nil,
				}
			}
		}
	}

	return &Helpers.RunnerResult{
		Line:   config.Line,
		Status: false,
		Error:  fmt.Errorf("Invalid"),
	}
}

func ScannerXSS(config *Helpers.ScanConfig) {
	Helpers.Valid, Helpers.Invalid, Helpers.Checked, Helpers.PayloadsTested, Helpers.CPM, Helpers.HighestCPM, Helpers.Running = 0, 0, 0, 0, 0, 0, true

	FilePath := CLI_Handlers.GetFilePath()
	lines, err := CLI_Handlers.ReadLines(FilePath)
	CLI_Handlers.LogError(err)

	startTime := time.Now()
	go func() {
		for {
			if Helpers.Running {
				mu.Lock()
				tested := strconv.Itoa(Helpers.PayloadsTested)
				elapsedTime := time.Since(startTime)
				valid := Helpers.Valid + Helpers.Invalid
				Helpers.CPM = int(Helpers.CalculateCPM(valid, elapsedTime))
				Helpers.HighestCPM = Helpers.BestCPM(Helpers.CPM, Helpers.HighestCPM)
				mu.Unlock()

				Interface.StatsTitle(fmt.Sprintf("NullOps | Payloads Tested: %v | CPM: %v | Highest CPM: %v", tested, Helpers.CPM, Helpers.HighestCPM), Helpers.Valid, Helpers.Invalid, Helpers.Checked, len(lines))
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
		RunnerResult := scanXSS(&ScanConfig)

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
