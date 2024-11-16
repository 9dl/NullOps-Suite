package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

var payloadPMA = []string{
	"phpmyadmin",
}

var payloadPMA2 = []string{
	"<bdo dir=\"ltr\" lang=\"en\">phpMyAdmin</bdo></h1><noscript>",
	"<form method=\"post\" action=\"index.php\" name=\"login_form\" class=\"disableAjax login hide js-show\">",
	"</form></div><div id=\"pma_footer\"></div></div></body></html>",
}

var payloadPMA3 = []string{
	"The requested URL was not found on this server.",
	"Not Found",
	"HTTPStatus::NotFound",
	"404",
	"Forbidden",
	"You don't have permission to access /",
	"You don't have permission to access",
	"Multiple Choices",
	"could not be found on this server",
	"404 Not Found",
	"Oops! An Error Occurred",
	"Fehler",
	"not found",
	"Fatal error",
	"Page Can't Be Found",
	"Page not found",
	"Error",
	"does not exist",
}

func scanPMA(config *Helpers.Runner) *Helpers.RunnerResult {
	for _, payload := range payloadPMA {

		mu.Lock()
		Helpers.PayloadsTested++
		mu.Unlock()

		Response, Error := Helpers.SendRequest(Helpers.ExtractDomain(config.Line)+payload, "GET", "", Helpers.RequestOptions{})

		if Error == nil {
			for _, payload2 := range payloadPMA2 {
				for _, payload3 := range payloadPMA3 {
					if strings.Contains(strings.ToLower(string(Response.Body)), strings.ToLower(payload2)) && !strings.Contains(strings.ToLower(string(Response.Body)), strings.ToLower(payload3)) {
						err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/PMA.txt", []string{"https://" + Helpers.ExtractDomain(config.Line) + "/" + payload})
						CLI_Handlers.LogError(err)
						return &Helpers.RunnerResult{
							Line:   Helpers.ExtractDomain(config.Line) + " | Panel: " + payload,
							Status: true,
							Error:  nil,
						}
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

func ScannerPMA(config *Helpers.ScanConfig) {
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
		RunnerResult := scanPMA(&ScanConfig)

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
