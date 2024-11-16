package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"
)

var payloadSQLi = []string{
	"\"",
	"'",
}

var databasePatterns = map[string][]string{
	"MySql":              {"SQL syntax.*MySQL", "Warning.*mysql_.*", "valid MySQL result", "MySqlClient\\."},
	"PostgreSQL":         {"PostgreSQL.*ERROR", "Warning.*\\Wpg_.*", "valid PostgreSQL result", "Npgsql\\."},
	"MicrosoftSQLServer": {"Driver.* SQL[\\-\\_\\ ]*Server", "OLE DB.* SQL Server", "(\\W|\\A)SQL Server.*Driver", "Warning.*mssql_.*", "(\\W|\\A)SQL Server.*[0-9a-fA-F]{8}", "(?s)Exception.*\\WSystem\\.Data\\.SqlClient\\.", "(?s)Exception.*\\WRoadhouse\\.Cms\\."},
	"MicrosoftAccess":    {"Microsoft Access Driver", "JET Database Engine", "Access Database Engine"},
	"Oracle":             {"\\bORA-[0-9][0-9][0-9][0-9]", "Oracle error", "Oracle.*Driver", "Warning.*\\Woci_.*", "Warning.*\\Wora_.*"},
	"IBMDB2":             {"CLI Driver.*DB2", "DB2 SQL error", "\\bdb2_\\w+\\("},
	"SQLite":             {"SQLite/JDBCDriver", "SQLite.Exception", "System.Data.SQLite.SQLiteException", "Warning.*sqlite_.*", "Warning.*SQLite3::", "\\[SQLITE_ERROR\\]"},
	"Sybase":             {"(?i)Warning.*sybase.*", "Sybase message", "Sybase.*Server message.*"},
}

func scanSQLi(config *Helpers.Runner) *Helpers.RunnerResult {
	for _, payload := range payloadSQLi {
		atomic.AddInt32(&Helpers.PayloadsTested, 1)
		Response, Error := Helpers.SendRequest(config.Line, "GET", payload, Helpers.RequestOptions{})

		if Error == nil {
			for db, patterns := range databasePatterns {
				for _, regg := range patterns {
					if regexp.MustCompile(regg).MatchString(string(Response.Body)) {
						err := CLI_Handlers.AppendToFile(fmt.Sprintf(Helpers.OutputPath+"/%v.txt", db), []string{config.Line})
						CLI_Handlers.LogError(err)
						return &Helpers.RunnerResult{
							Line:   Helpers.ExtractDomain(config.Line) + " | DB: " + db,
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

func ScannerSQLi(config *Helpers.ScanConfig) {
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
		RunnerResult := scanSQLi(&ScanConfig)

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
