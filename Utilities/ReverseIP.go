package Utilities

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
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
			tested := strconv.Itoa(int(atomic.LoadInt32(&Helpers.PayloadsTested)))
			Interface.StatsTitle("NullOps | Payloads Tested: "+tested, int(atomic.LoadInt32(&Helpers.Valid)), int(atomic.LoadInt32(&Helpers.Invalid)), int(atomic.LoadInt32(&Helpers.Checked)), len(lines))
			time.Sleep(2 * time.Second)
		}
	}()

	Helpers.Threading(func(s string) {
		ScanConfig := Helpers.Runner{Line: s}
		RunnerResult := scanReverseIP(&ScanConfig)

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
