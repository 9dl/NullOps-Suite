package Interface

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"os"
	"os/exec"
	"time"
)

func getUsage() (float64, float64, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, 0, err
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, err
	}

	cpuUsage := cpuPercent[0]
	memoryUsage := float64(memory.Used) / float64(memory.Total) * 100

	return cpuUsage, memoryUsage, nil
}

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func Title(title string) {
	fmt.Printf("\033]0;%s\007", title)
}

func StatsTitle(tool string, valid int, invalid int, checked int, remaining int) {
	go func() {
		cpuUsage, memoryUsage, err := getUsage()
		if err != nil {
			fmt.Println(err)
		}
		percentageRemaining := float64(remaining) / float64(checked) * 100

		fmt.Printf("\033]0;%s | Valid: %v | Invalid: %v | Checked: %v/%v [%.2f%%] | CPU [%.2f%%] & RAM [%.2f%%] \007", tool, valid, invalid, checked, remaining, percentageRemaining, cpuUsage, memoryUsage)
	}()
}

func DumperTitle(tool string, valid int, invalid int, empty int, checked int, remaining int) {
	go func() {
		cpuUsage, memoryUsage, err := getUsage()
		if err != nil {
			fmt.Println(err)
		}
		percentageRemaining := float64(remaining) / float64(checked) * 100

		fmt.Printf("\033]0;%s | Dumped: %v | Failed: %v | Empty: %v | Checked: %v/%v [%.2f%%] | CPU [%.2f%%] & RAM [%.2f%%] \007", tool, valid, invalid, empty, checked, remaining, percentageRemaining, cpuUsage, memoryUsage)
	}()
}
