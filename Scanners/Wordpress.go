package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"strings"
	"time"
)

var wp_vulnerablePaths = []string{
	"xmlrpc.php",
	//"wp-config.php",
	"wp-cron.php",
	"wp-includes",
	"wp-content",
	"wp-json",
	//"wp-admin",
	"wp-login.php",
	"wp-content/uploads",
	"wp-content/plugins",
	"wp-content/themes",
	"wp-includes/js/jquery/jquery.js",
	"wp-includes/js/jquery/jquery-migrate.min.js",
	"wp-includes/css/dashicons.min.css",
	"wp-includes/css/admin-bar.min.css",
}

func checkVulnerabilities(responseBody string, pathToCheck string) string {
	if contains(wp_vulnerablePaths, pathToCheck) {
		switch pathToCheck {
		case "xmlrpc.php":
			if strings.Contains(responseBody, "XML-RPC server accepts POST requests only") {
				return "XMLRPC"
			}
		case "wp-config.php":
			if strings.Contains(responseBody, "wp-config.php") {
				return "WPConfig"
			}
		case "wp-cron.php":
			if strings.Contains(responseBody, "wp-cron.php") {
				return "WPCron"
			}
		case "wp-includes":
			if strings.Contains(responseBody, "Index of /wp-includes") {
				return "WPIncludes"
			}
		case "wp-content":
			if strings.Contains(responseBody, "Index of /wp-content") {
				return "WPContent"
			}
		case "wp-json":
			if strings.Contains(responseBody, "rest_login_required") || strings.Contains(responseBody, "rest_cannot_access") {
				return "WPJSONDisabled"
			} /* else { // User Enumeration
				usersURL := website + "/wp-json/wp/v2/users"
				userData := `
						[
							{"id": 1, "name": "user1", "slug": "user1"},
							{"id": 2, "name": "user2", "slug": "user2"},
							{"id": 3, "name": "admin", "slug": "admin"}
						]
					`
				if strings.Contains(userData, "slug") {
					return "UserEnumerationSuccessful"
				}
				return "UserEnumerationFailed"
			}*/
		case "wp-admin":
			if strings.Contains(responseBody, "wp-admin") {
				return "WPAdmin"
			}
		case "wp-login.php":
			if strings.Contains(responseBody, "wp-login.php") {
				return "WPLogin"
			}
		case "wp-content/uploads":
			if strings.Contains(responseBody, "wp-content/uploads") {
				return "WPContentUploads"
			}
		case "wp-content/plugins":
			if strings.Contains(responseBody, "wp-content/plugins/") {
				return "VulnerablePlugins"
			}
		case "wp-content/themes":
			if strings.Contains(responseBody, "wp-content/themes/") {
				return "AccessibleThemes"
			}
		case "wp-includes/js/jquery/jquery.js":
			if strings.Contains(responseBody, "jquery.js") {
				return "VulnerableJQuery"
			}
		case "wp-includes/js/jquery/jquery-migrate.min.js":
			if strings.Contains(responseBody, "jquery-migrate.min.js") {
				return "VulnerableJQueryMigrate"
			}
		case "wp-includes/css/dashicons.min.css":
			if strings.Contains(responseBody, "dashicons.min.css") {
				return "DashiconsCSS"
			}
		case "wp-includes/css/admin-bar.min.css":
			if strings.Contains(responseBody, "admin-bar.min.css") {
				return "AdminBarCSS"
			}
		default:
			return "UnknownPath"
		}
	}
	return "Invalid"
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func scanWordpress(config *Helpers.Runner) *Helpers.RunnerResult {
	for _, payload := range wp_vulnerablePaths {
		Response, Error := Helpers.SendRequest(Helpers.ExtractDomain(config.Line)+payload, "GET", "", Helpers.RequestOptions{})

		Vuln := checkVulnerabilities(string(Response.Body), payload)
		if Error != nil {
			err := CLI_Handlers.AppendToFile(Helpers.OutputPath+fmt.Sprintf("/%s.txt", Vuln), []string{
				Helpers.ExtractDomain(config.Line) + payload,
			},
			)

			CLI_Handlers.LogError(err)

			return &Helpers.RunnerResult{
				Line:   Helpers.ExtractDomain(config.Line),
				Status: true,
				Error:  nil,
			}
		}
	}

	return &Helpers.RunnerResult{
		Line:   config.Line,
		Status: false,
		Error:  fmt.Errorf("Invalid"),
	}
}

func ScannerWordpress(config *Helpers.ScanConfig) {
	Helpers.Valid, Helpers.Invalid, Helpers.Checked, Helpers.PayloadsTested, Helpers.CPM, Helpers.HighestCPM, Helpers.Running = 0, 0, 0, 0, 0, 0, true

	FilePath := CLI_Handlers.GetFilePath()
	lines, err := CLI_Handlers.ReadLines(FilePath)
	CLI_Handlers.LogError(err)

	startTime := time.Now()
	go func() {
		for {
			if Helpers.Running {
				mu.Lock()
				elapsedTime := time.Since(startTime)
				valid := Helpers.Valid + Helpers.Invalid
				Helpers.CPM = int(Helpers.CalculateCPM(valid, elapsedTime))
				Helpers.HighestCPM = Helpers.BestCPM(Helpers.CPM, Helpers.HighestCPM)
				mu.Unlock()

				Interface.StatsTitle(fmt.Sprintf("NullOps | CPM: %v | Highest CPM: %v", Helpers.CPM, Helpers.HighestCPM), Helpers.Valid, Helpers.Invalid, Helpers.Checked, len(lines))
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
		RunnerResult := scanWordpress(&ScanConfig)

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
