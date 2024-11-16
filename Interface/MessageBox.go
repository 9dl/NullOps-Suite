package Interface

import (
	"fmt"
	"github.com/gen2brain/dlgs"
)

func ErrorMSG(title, message string) {
	_, err := dlgs.Error(title, message)
	if err != nil {
		fmt.Println("Error displaying message box:", err)
	}
}

func InfoMSG(title, message string) {
	_, err := dlgs.Info(title, message)
	if err != nil {
		fmt.Println("Error displaying message box:", err)
	}
}

func WarningMSG(title, message string) {
	_, err := dlgs.Warning(title, message)
	if err != nil {
		fmt.Println("Error displaying message box:", err)
	}
}
