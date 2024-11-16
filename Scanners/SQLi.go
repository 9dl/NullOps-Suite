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
		mu.Lock()
		Helpers.PayloadsTested++
		mu.Unlock()
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
		RunnerResult := scanSQLi(&ScanConfig)

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
