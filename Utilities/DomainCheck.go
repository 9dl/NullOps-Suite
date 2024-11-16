package Utilities

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

func scanDomain(config *Helpers.Runner) *Helpers.RunnerResult {
	_, Error := Helpers.SendRequest("https://"+strings.ReplaceAll(strings.ReplaceAll(config.Line, "https://", ""), "http://", ""), "GET", "", Helpers.RequestOptions{})
	atomic.AddInt32(&Helpers.PayloadsTested, 1)

	if Error == nil {
		err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"Domains.txt", []string{config.Line})
		CLI_Handlers.LogError(err)

		return &Helpers.RunnerResult{
			Line:   Helpers.ExtractDomain(config.Line),
			Status: true,
			Error:  nil,
		}
	}

	return &Helpers.RunnerResult{
		Line:   config.Line,
		Status: false,
		Error:  fmt.Errorf("Invalid"),
	}
}

func ScannerDomain(config *Helpers.ScanConfig) {
	Helpers.Valid, Helpers.Invalid, Helpers.Checked, Helpers.PayloadsTested = 0, 0, 0, 0

	FilePath := CLI_Handlers.GetFilePath()
	lines, err := CLI_Handlers.ReadLines(FilePath)
	CLI_Handlers.LogError(err)

	go func() {
		for {
			Interface.StatsTitle("NullOps", int(atomic.LoadInt32(&Helpers.Valid)), int(atomic.LoadInt32(&Helpers.Invalid)), int(atomic.LoadInt32(&Helpers.Checked)), len(lines))
			time.Sleep(2 * time.Second)
		}
	}()

	Helpers.Threading(func(s string) {
		ScanConfig := Helpers.Runner{Line: s}
		RunnerResult := scanDomain(&ScanConfig)

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
