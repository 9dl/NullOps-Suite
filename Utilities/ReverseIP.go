package Utilities

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"net"
	"strconv"
	"time"
)

func scanReverseIP(config *Helpers.Runner) *Helpers.RunnerResult {
	ips, err := net.LookupIP(Helpers.ExtractHost(config.Line))
	if err != nil {
		return &Helpers.RunnerResult{
			Line:   config.Line,
			Status: false,
			Error:  fmt.Errorf("Invalid Domain"),
		}
	}

	if len(ips) > 0 {
		ipStrings := make([]string, len(ips))
		for i, ip := range ips {
			ipStrings[i] = ip.String()
		}
		err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/ReverseIP.txt", ipStrings)
		CLI_Handlers.LogError(err)

		return &Helpers.RunnerResult{
			Line:   config.Line,
			Status: false,
			Error:  nil,
		}
	}

	return &Helpers.RunnerResult{
		Line:   config.Line,
		Status: false,
		Error:  fmt.Errorf("No IP's"),
	}
}

func ScannerReverseIP(config *Helpers.ScanConfig) {
	Helpers.Valid, Helpers.Invalid, Helpers.Checked, Helpers.PayloadsTested = 0, 0, 0, 0

	FilePath := CLI_Handlers.GetFilePath()
	lines, err := CLI_Handlers.ReadLines(FilePath)
	CLI_Handlers.LogError(err)
	go func() {
		for {
			tested := strconv.Itoa(Helpers.PayloadsTested)
			Interface.StatsTitle("NullOps | Payloads Tested: "+tested, Helpers.Valid, Helpers.Invalid, Helpers.Checked, len(lines))
			time.Sleep(2 * time.Second)
		}
	}()

	Helpers.Threading(func(s string) {
		ScanConfig := Helpers.Runner{Line: s}
		RunnerResult := scanReverseIP(&ScanConfig)

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
