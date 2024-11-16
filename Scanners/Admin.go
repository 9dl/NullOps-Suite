package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var payloadAdmin = []string{
	"admin.php",
	"admin/login.php",
	"administrator.php",
	"administrator",
	"adminpanel.php",
	"adminlogin.php",
	"admincp.php",
	"adminer.php",
}

var payloadAdmin2 = []string{
	"login",
	"admin",
	"password",
	"user",
	"panel",
	"controlpanel",
	"administrator",
}

var payloadAdmin3 = []string{
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
}

func scanAdmin(config *Helpers.Runner) *Helpers.RunnerResult {
	for _, payload := range payloadAdmin {
		atomic.AddInt32(&Helpers.PayloadsTested, 1)
		Response, Error := Helpers.SendRequest(Helpers.ExtractDomain(config.Line)+payload, "GET", "", Helpers.RequestOptions{})

		if Error == nil {
			for _, payload2 := range payloadAdmin2 {
				for _, payload3 := range payloadAdmin3 {
					if strings.Contains(strings.ToLower(string(Response.Body)), strings.ToLower(payload2)) && !strings.Contains(strings.ToLower(string(Response.Body)), strings.ToLower(payload3)) {
						err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/Admin.txt", []string{"https://" + Helpers.ExtractDomain(config.Line) + "/" + payload})
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

func ScannerAdmin(config *Helpers.ScanConfig) {
	Helpers.Valid, Helpers.Invalid, Helpers.Checked, Helpers.PayloadsTested, Helpers.CPM, Helpers.HighestCPM, Helpers.Running = 0, 0, 0, 0, 0, 0, true

	FilePath := CLI_Handlers.GetFilePath()
	lines, err := CLI_Handlers.ReadLines(FilePath)
	CLI_Handlers.LogError(err)

	startTime := time.Now()
	go func() {
		for {
			if Helpers.Running {
				tested := strconv.Itoa(int(atomic.LoadInt32(&Helpers.PayloadsTested)))
				elapsedTime := time.Since(startTime)
				Helpers.CPM = int32(int(Helpers.CalculateCPM(int(atomic.LoadInt32(&Helpers.Valid))+int(atomic.LoadInt32(&Helpers.Invalid)), elapsedTime)))
				Helpers.HighestCPM = int32(Helpers.BestCPM(int(Helpers.CPM), int(atomic.LoadInt32(&Helpers.HighestCPM))))

				Interface.StatsTitle(fmt.Sprintf("NullOps | Payloads Tested: %v | CPM: %v | Highest CPM: %v", tested, int(atomic.LoadInt32(&Helpers.CPM)), int(atomic.LoadInt32(&Helpers.HighestCPM))), int(atomic.LoadInt32(&Helpers.Valid)), int(atomic.LoadInt32(&Helpers.Invalid)), int(atomic.LoadInt32(&Helpers.Checked)), len(lines))
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
		RunnerResult := scanAdmin(&ScanConfig)

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