package main

import (
	"NullOps/CLI_Handlers"
	"NullOps/Dumpers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"NullOps/Scanners"
	"NullOps/Utilities"
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var configManager = Helpers.NewConfigurationManager()

func main() {
	configManager.LoadConfig()
	fmt.Print("\033[?25l")
	runtime.SetBlockProfileRate(1)
	go func() {
		for {
			if !Helpers.Running {
				Interface.Title("NullOps")
				time.Sleep(2 * time.Second)
			}
		}
	}()

	ExitKey := make(chan os.Signal, 1)
	signal.Notify(ExitKey, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-ExitKey
		if Helpers.Running {
			Helpers.ShowResults()
		}
		Interface.Option("KEY DETECTED", "CTRL + C -> EXITING")
		os.Exit(1)
	}()

	Interface.LogoASCII = []string{
		"███╗   ██╗██╗   ██╗██╗     ██╗      ██████╗ ██████╗ ███████╗",
		"████╗  ██║██║   ██║██║     ██║     ██╔═══██╗██╔══██╗██╔════╝",
		"██╔██╗ ██║██║   ██║██║     ██║     ██║   ██║██████╔╝███████╗",
		"██║╚██╗██║██║   ██║██║     ██║     ██║   ██║██╔═══╝ ╚════██║",
		"██║ ╚████║╚██████╔╝███████╗███████╗╚██████╔╝██║     ███████║",
		"╚═╝  ╚═══╝ ╚═════╝ ╚══════╝╚══════╝ ╚═════╝ ╚═╝     ╚══════╝",
	}

	NullOps()
}
func SetTerminalHeight(lines int) {
	if lines <= 0 {
		lines = 1
	}
	fmt.Printf("\033[1;%dH", lines)
}

func NullOps() {
	SetTerminalHeight(10)
	configManager.LoadConfig()
	err := os.Mkdir(Helpers.OutputPath, 0750)
	CLI_Handlers.LogError(err)

	var Option string

	Interface.Logo()
	Interface.Title("NullOps")

	Helpers.ThreadingType = configManager.GetThreadingType()
	Helpers.Timeout = int32(configManager.GetTimeout())

	Interface.WriteColoredCentered("[«] ================ Threading: ["+configManager.GetThreadingType()+"] ================ [»]", pterm.NewRGB(0, 255, 255))
	Interface.WriteColoredCentered("[«] ================ Scanners ================ ================ Utilities ================  [»]", pterm.NewRGB(0, 255, 255))
	Interface.Option2("01", "CMS Scanner")
	Interface.Option4("A", "\t  [IPv4] IP Checker")

	Interface.Option2("02", "Laravel Env Scanner")
	Interface.Option4("B", "   [url.com] Domain Checker")

	Interface.Option2("03", "LFI Scanner")
	Interface.Option4("C", "  [IPv4 -> Domain] Reverse DNS")

	Interface.Option2("04", "RCE Scanner")
	Interface.Option4("D", "  [Domain -> IPv4's] Reverse IP")

	Interface.Option2("05", "SQLi Scanner")
	Interface.Option4("E", "              Column Grabber")

	Interface.Option2("06", "XSS Scanner")
	Interface.Option4("F", "               Port Scanner")

	Interface.Option2("07", "Env Scanner")
	fmt.Println()

	Interface.Option2("08", "Admin Panels Scanner")
	fmt.Println()

	Interface.Option2("09", "phpMyAdmin Scanner")
	fmt.Println()

	Interface.Option2("10", "Adminer Scanner")
	fmt.Println()

	Interface.Option2("11", "cPanel Scanner")
	fmt.Println()

	Interface.Option2("12", "SFTP Config Scanner")
	fmt.Println()

	Interface.Option2("13", "Wordpress Scanner")
	fmt.Println()

	Interface.Option2("14", "Exposed Files Scanner")
	Interface.Option4("X", "                SQLi Dumper")

	Interface.Option2("15", "Git Config Scanner")
	Interface.Option4("Z", "        Env & Laravel Dumper")

	Interface.WriteColoredCentered("[«] ==================================================================================== [»]", pterm.NewRGB(0, 255, 255))
	Interface.WriteColoredCentered(fmt.Sprintf("[«] ========== Threads: [%v] | Timeout (ms): [%v] ========== [»]", configManager.GetThreads(), configManager.GetTimeout()), pterm.NewRGB(0, 255, 255))
	Interface.WriteColoredCentered("[«] ========== [github.com/9dl/NullOps-Suite] ========== [»]", pterm.NewRGB(0, 255, 255))
	
	Interface.Input2()
	_, err = fmt.Scanln(&Option)
	CLI_Handlers.LogError(err)

	if Option == "01" || Option == "1" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "CMS", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerCMS(&ScannerConfig)
	} else if Option == "02" || Option == "2" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "Laravel", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerLaravel(&ScannerConfig)
	} else if Option == "03" || Option == "3" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "LFI", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerLFI(&ScannerConfig)
	} else if Option == "04" || Option == "4" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "RCE", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerRCE(&ScannerConfig)
	} else if Option == "05" || Option == "5" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "SQLi", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerSQLi(&ScannerConfig)
	} else if Option == "06" || Option == "6" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "XSS", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerXSS(&ScannerConfig)
	} else if Option == "07" || Option == "7" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "Env", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerEnv(&ScannerConfig)
	} else if Option == "08" || Option == "8" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "Admin", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerAdmin(&ScannerConfig)
	} else if Option == "09" || Option == "9" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "PMA", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerPMA(&ScannerConfig)
	} else if Option == "10" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "Adminer", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerAdminer(&ScannerConfig)
	} else if Option == "11" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "cPanel", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerCPanel(&ScannerConfig)
	} else if Option == "12" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "SFTP", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerSFTP(&ScannerConfig)
	} else if Option == "13" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "Wordpress", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerWordpress(&ScannerConfig)
	} else if Option == "14" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "Exposed", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerVulnFiles(&ScannerConfig)
	} else if Option == "15" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "Git", PrintInvalid: configManager.GetPrintInvalid()}
		Scanners.ScannerGit(&ScannerConfig)
	} else if strings.ToLower(Option) == "a" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "IP Scanner", PrintInvalid: configManager.GetPrintInvalid()}
		Utilities.ScannerIP(&ScannerConfig)
	} else if strings.ToLower(Option) == "b" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "Domain Scanner", PrintInvalid: configManager.GetPrintInvalid()}
		Utilities.ScannerDomain(&ScannerConfig)
	} else if strings.ToLower(Option) == "c" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "ReverseDNS", PrintInvalid: configManager.GetPrintInvalid()}
		Utilities.ScannerReverseDNS(&ScannerConfig)
	} else if strings.ToLower(Option) == "tt" {
		ThreadingSys()
	} else if strings.ToLower(Option) == "d" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "ReverseIP", PrintInvalid: configManager.GetPrintInvalid()}
		Utilities.ScannerReverseIP(&ScannerConfig)
	} else if strings.ToLower(Option) == "z" {
		Dumpers.LaravelnEnvDumper()
	} else if strings.ToLower(Option) == "e" {
		Utilities.ColumnGrabber()
	} else if strings.ToLower(Option) == "f" {
		ScannerConfig := Helpers.ScanConfig{Threads: configManager.GetThreads(), Name: "PortsScanner", PrintInvalid: configManager.GetPrintInvalid()}
		Utilities.ScannerPorts(&ScannerConfig)
	} else if strings.ToLower(Option) == "x" {
		Dumpers.SQLiDumper()
	} else {
		NullOps()
	}
}

func ThreadingSys() {
	Interface.Clear()
	Interface.Logo()
	Interface.WriteColoredCentered("[«] ================ [Sentry] ================ [»]", pterm.NewRGB(0, 255, 255))
	Interface.WriteColoredCentered("[«] Keeps your system running smoothly by managing how many tasks it handles at once. [»]", pterm.NewRGB(0, 255, 255))

	fmt.Println()

	Interface.WriteColoredCentered("[«] ================ [Guardian] ================ [»]", pterm.NewRGB(0, 255, 255))
	Interface.WriteColoredCentered("[«] Efficiently handles tasks in groups for faster processing. [»]", pterm.NewRGB(0, 255, 255))

	fmt.Println()

	Interface.WriteColoredCentered("[«] ================ [Lympia] ================ [»]", pterm.NewRGB(0, 255, 255))
	Interface.WriteColoredCentered("[«] Processes lines one-by-one, ensuring each is unique before sending to workers. [»]", pterm.NewRGB(0, 255, 255))

	fmt.Println()

	Interface.WriteColoredCentered("[«] ================ [Vortex] ================ [»]", pterm.NewRGB(0, 255, 255))
	Interface.WriteColoredCentered("[«] Optimized for speed, Vortex parallelizes both task distribution and execution for maximum throughput. [»]", pterm.NewRGB(0, 255, 255))

	fmt.Println()

	Interface.WriteColoredCentered("[«] ================ [Threading Selection Mechanism] ================ [»]", pterm.NewRGB(0, 255, 255))
	Interface.WriteColoredCentered("[«] Automatic threading selection is based on CPU load and core count. [»]", pterm.NewRGB(0, 255, 255))

	fmt.Println()
	Interface.WriteColoredCentered("[«] press any key to go back [»]", pterm.NewRGB(0, 255, 255))

	_, err := fmt.Scanln()
	CLI_Handlers.LogError(err)
	NullOps()
}
