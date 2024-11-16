package Interface

import (
	"fmt"
	"github.com/pterm/pterm"
	"golang.org/x/term"
	"os"
)

func getTerminalWidth() int {
	screenWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		screenWidth = 80
	}
	return screenWidth
}

var LogoASCII = []string{
	"NullOps: Logo isn't set!",
}

func Logo() {
	Clear()
	fmt.Println()
	fmt.Println()
	var R uint8 = 0
	var G uint8 = 255
	var B uint8 = 255

	for i := 0; i < len(LogoASCII); i++ {
		pterm.DefaultCenter.Print(pterm.NewRGB(R, G, B).Sprint(LogoASCII[i]))
		G -= 29
		B -= 29
	}
	fmt.Println()
}
