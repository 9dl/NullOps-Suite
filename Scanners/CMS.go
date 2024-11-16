package Scanners

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

func DetectCMS(responseHeaders map[string]string, responseBody string) string {
	var detectedCMS []string
	bodyPatterns := map[string]string{
		"Drupal":          `sites/all/themes/`,
		"Craft CMS":       `CraftSessionId`,
		"Joomla":          `X-Content-Encoded-By: Joomla`,
		"WordPress":       `wp-content/plugins/`,
		"vBulletin":       `vbulletin.js`,
		"Concrete5":       `concrete/js/`,
		"Contao":          `Contao Core Files`,
		"DataLife Engine": `DataLife Engine`,
		"Django":          `csrfmiddlewaretoken`,
		"Grav":            `grav-cms.js`,
		"PrestaShop":      `prestashop/js/`,
	}

	headerPatterns := map[string]string{
		"OpenCms":          `Server: OpenCms`,
		"PHP-Nuke":         `Powered by PHP-Nuke`,
		"SPIP":             `X-Spip-Cache`,
		"WebGUI":           `generator" content="WebGUI"`,
		"Laravel":          `Set-Cookie: laravel_session`,
		"DokuWiki":         `Set-Cookie: DokuWiki=`,
		"eSyndiCat":        `X-Directory-Script: eSyndiCat`,
		"eZ Publish":       `X-Powered-By: eZ Publish`,
		"GetSimple CMS":    `generator" content="GetSimple"`,
		"Kotisivukone":     `kotisivukone.js`,
		"Koala Framework":  `Koala Web Framework CMS`,
		"Kooboo CMS":       `X-KoobooCMS-Version`,
		"InstantCMS":       `Set-Cookie: InstantCMS`,
		"Liferay":          `Liferay-Portal`,
		"FlexCMP":          `X-Powered-By: FlexCMP`,
		"Sarka-SPIP":       `Sarka-SPIP`,
		"Green Valley CMS": `dsresource?objectid`,
		"Graffiti CMS":     `Set-Cookie: graffitibot`,
		"1C-Bitrix":        `X-Powered-CMS: Bitrix Site Manager`,
		"Cloudflare":       `cdnjs.cloudflare.com`,
		"Cloudfront":       `cloudfront.net`,
	}

	var responseHeadersStr string
	for key, value := range responseHeaders {
		responseHeadersStr += key + ": " + value + "\n"
	}

	// Detect CMS from response body
	for cms, pattern := range bodyPatterns {
		match, err := regexp.MatchString(pattern, responseBody)
		if err == nil && match {
			detectedCMS = append(detectedCMS, cms)
		}
	}

	// Detect CMS from response headers
	for cms, pattern := range headerPatterns {
		if strings.Contains(responseHeadersStr, pattern) {
			detectedCMS = append(detectedCMS, cms)
		}
	}

	if len(detectedCMS) == 0 {
		return "Unknown"
	}

	var detectedCMSStr string
	for _, cms := range detectedCMS {
		detectedCMSStr += cms + ","
	}

	return detectedCMSStr
}

func scanCMS(config *Helpers.Runner) *Helpers.RunnerResult {
	var detectedCMS = ""
	Response, Error := Helpers.SendRequest(config.Line, "GET", "", Helpers.RequestOptions{})

	if Error == nil {
		detectedCMS = DetectCMS(Response.Headers, string(Response.Body))
		if detectedCMS != "Unknown" {
			err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/"+strings.ReplaceAll(detectedCMS, ",", " ")+".txt", []string{config.Line})
			CLI_Handlers.LogError(err)
			return &Helpers.RunnerResult{
				Line:   Helpers.ExtractDomain(config.Line) + " | CMS: " + strings.ReplaceAll(detectedCMS, ",", " "),
				Status: true,
				Error:  nil,
			}
		}
	}

	return &Helpers.RunnerResult{
		Line:   config.Line,
		Status: false,
		Error:  fmt.Errorf(detectedCMS),
	}
}

func ScannerCMS(config *Helpers.ScanConfig) {
	Helpers.Valid, Helpers.Invalid, Helpers.Checked, Helpers.PayloadsTested, Helpers.CPM, Helpers.HighestCPM, Helpers.Running = 0, 0, 0, 0, 0, 0, true

	FilePath := CLI_Handlers.GetFilePath()
	lines, err := CLI_Handlers.ReadLines(FilePath)
	CLI_Handlers.LogError(err)

	startTime := time.Now()
	go func() {
		for {
			if Helpers.Running {
				elapsedTime := time.Since(startTime)
				Helpers.CPM = int32(int(Helpers.CalculateCPM(int(atomic.LoadInt32(&Helpers.Valid))+int(atomic.LoadInt32(&Helpers.Invalid)), elapsedTime)))
				Helpers.HighestCPM = int32(Helpers.BestCPM(int(Helpers.CPM), int(atomic.LoadInt32(&Helpers.HighestCPM))))

				Interface.StatsTitle(fmt.Sprintf("NullOps | CPM: %v | Highest CPM: %v", int(atomic.LoadInt32(&Helpers.CPM)), int(atomic.LoadInt32(&Helpers.HighestCPM))), int(atomic.LoadInt32(&Helpers.Valid)), int(atomic.LoadInt32(&Helpers.Invalid)), int(atomic.LoadInt32(&Helpers.Checked)), len(lines))
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
		RunnerResult := scanCMS(&ScanConfig)

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
