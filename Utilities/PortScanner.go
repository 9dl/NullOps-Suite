package Utilities

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

var startPort, endPort int

func scanPorts(config *Helpers.Runner) *Helpers.RunnerResult {

	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf("%s:%d", config.Line, port)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/OpenPorts.txt", []string{address})
			CLI_Handlers.LogError(err)
			return &Helpers.RunnerResult{
				Line:   address,
				Status: true,
				Error:  nil,
			}
			err = conn.Close()
			CLI_Handlers.LogError(err)
		}
	}
	return &Helpers.RunnerResult{
		Line:   config.Line,
		Status: false,
		Error:  fmt.Errorf("Invalid IP"),
	}
}

func ScannerPorts(config *Helpers.ScanConfig) {
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

	Interface.Option("?", "Start Port")
	Interface.Input()
	_, err = fmt.Scanln(&startPort)
	CLI_Handlers.LogError(err)

	Interface.Option("?", "End Port")
	Interface.Input()
	_, err = fmt.Scanln(&endPort)
	CLI_Handlers.LogError(err)

	Interface.Clear()

	Helpers.Threading(func(s string) {
		ScanConfig := Helpers.Runner{Line: s}
		RunnerResult := scanPorts(&ScanConfig)

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
