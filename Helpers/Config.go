package Helpers

import (
	"time"
)

var (
	Valid          int32  = 0
	Invalid        int32  = 0
	Checked        int32  = 0
	PayloadsTested int32  = 0
	CPM            int32  = 0
	HighestCPM     int32  = 0
	Timeout        int32  = 5000
	Running        bool   = false
	OutputPath            = time.Now().Format("15;04;05")
	ThreadingType  string = "Lympia"
)

type ScanConfig struct {
	Threads      int
	Name         string
	PrintInvalid bool
}

type Runner struct {
	Line string
}

type RunnerResult struct {
	Line   string
	Status bool
	Error  error
}
