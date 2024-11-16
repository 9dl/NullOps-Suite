package Utilities

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func cleanDomainName(name string) string {
	name = strings.TrimPrefix(name, "dns.")
	name = strings.TrimSuffix(name, ".")
	return name
}

func scanReverseDNS(config *Helpers.Runner) *Helpers.RunnerResult {
	names, err := net.LookupAddr(config.Line)
	if err != nil {
		return &Helpers.RunnerResult{
			Line:   config.Line,
			Status: false,
			Error:  fmt.Errorf("Invalid IP"),
		}
	}

	if len(names) > 0 {
		domain := cleanDomainName(names[0])
		err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/ReverseDNS.txt", []string{domain})
		CLI_Handlers.LogError(err)
		return &Helpers.RunnerResult{
			Line:   domain,
			Status: true,
			Error:  nil,
		}
	}

	return &Helpers.RunnerResult{
		Line:   config.Line,
		Status: false,
		Error:  fmt.Errorf("No Domain"),
	}
}

func ScannerReverseDNS(config *Helpers.ScanConfig) {
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
		RunnerResult := scanReverseDNS(&ScanConfig)

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
