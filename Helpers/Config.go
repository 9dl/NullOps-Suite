package Helpers

import (
	"time"
)

var (
	Valid          int    = 0
	Invalid        int    = 0
	Checked        int    = 0
	PayloadsTested int    = 0
	CPM            int    = 0
	HighestCPM     int    = 0
	Timeout        int    = 5000
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
