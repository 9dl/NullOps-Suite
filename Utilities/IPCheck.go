package Utilities

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"net"
	"time"
)

func scanIP(config *Helpers.Runner) *Helpers.RunnerResult {
	conn, err := net.DialTimeout("ip4:icmp", config.Line, 5*time.Second)
	if err != nil {
		return &Helpers.RunnerResult{
			Line:   config.Line,
			Status: false,
			Error:  fmt.Errorf("Invalid IP"),
		}
	}
	defer conn.Close()

	err = CLI_Handlers.AppendToFile(Helpers.OutputPath+"/IPCheck.txt", []string{config.Line})
	CLI_Handlers.LogError(err)

	return &Helpers.RunnerResult{
		Line:   config.Line,
		Status: true,
		Error:  nil,
	}
}

func ScannerIP(config *Helpers.ScanConfig) {
	Helpers.Valid, Helpers.Invalid, Helpers.Checked, Helpers.PayloadsTested = 0, 0, 0, 0

	FilePath := CLI_Handlers.GetFilePath()
	lines, err := CLI_Handlers.ReadLines(FilePath)
	CLI_Handlers.LogError(err)

	go func() {
		for {
			Interface.StatsTitle("NullOps", Helpers.Valid, Helpers.Invalid, Helpers.Checked, len(lines))
			time.Sleep(2 * time.Second)
		}
	}()

	Helpers.Threading(func(s string) {
		ScanConfig := Helpers.Runner{Line: s}
		RunnerResult := scanIP(&ScanConfig)

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
