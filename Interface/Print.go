package Interface

import (
	"fmt"
	"github.com/pterm/pterm"
	"regexp"
	"strings"
)

func WriteColoredCentered(message string, colorAttribute pterm.RGB) {
	screenWidth := getTerminalWidth()
	cleanMessage := strings.ReplaceAll(message, "[", "")
	cleanMessage = strings.ReplaceAll(cleanMessage, "]", "")
	messageLength := len(cleanMessage)

	regexPattern := regexp.MustCompile(`(\[[^\]]*\])`)
	matches := regexPattern.FindAllStringIndex(message, -1)

	padding := (screenWidth - messageLength) / 2
	fmt.Print(strings.Repeat(" ", padding))

	previousIndex := 0
	for _, match := range matches {
		start, end := match[0], match[1]

		// Print the non-colored text
		pterm.Print(message[previousIndex:start])

		// Print the colored text
		coloredText := message[start+1 : end-1]
		pterm.Print(colorAttribute.Sprint(coloredText))

		previousIndex = end
	}
	// Print any remaining non-colored text
	pterm.Print(message[previousIndex:])

	fmt.Println()
}

func Gradient(Text string) {
	from := pterm.NewRGB(0, 255, 255)
	to := pterm.NewRGB(63, 78, 77)

	str := Text
	strs := strings.Split(str, "")
	var fadeInfo string
	for i := 0; i < len(str); i++ {
		fadeInfo += from.Fade(0, float32(len(str)), float32(i), to).Sprint(strs[i])
	}

	pterm.Println(pterm.Gray("[") + pterm.Blue("~") + pterm.Gray("] ") + fadeInfo)
}

func Write(Text string) {
	pterm.Println(pterm.Gray("[") + pterm.Cyan(">") + pterm.Gray("] ") + pterm.White(Text))
}

func Valid(Text string) {
	pterm.Println(pterm.Gray("[") + pterm.Green("+") + pterm.Gray("] ") + pterm.White(Text))
}

func Invalid(Text string) {
	pterm.Println(pterm.Gray("[") + pterm.Red("-") + pterm.Gray("] ") + pterm.White(Text))
}

func Info(Text string) {
	pterm.Println(pterm.Gray("[") + pterm.Magenta("?") + pterm.Gray("] ") + pterm.White(Text))
}

func Option(Number string, Text string) {
	pterm.Println(pterm.Gray("[") + pterm.Cyan(Number) + pterm.Gray("] ") + pterm.White(Text))
}

func Option2(Number string, Text string) {
	Width := 53
	availableSpace := Width - len(Number) - 19

	if len(Text) > availableSpace {
		Text = Text[:availableSpace]
	} else {
		Text = Text + strings.Repeat(" ", availableSpace-len(Text))
	}

	formattedText := pterm.White(fmt.Sprintf("	  	  %v", Text)) + pterm.Gray("  <= [") + pterm.Cyan(Number) + pterm.Gray("]")
	pterm.Print(formattedText)
}

func Option4(Number string, Text string) {
	width := 35 // Set the desired width here
	padding := width - len(Number) - len(Text)

	leftPadding := padding / 8
	rightPadding := padding - leftPadding

	formattedText := pterm.Gray("[") + pterm.Cyan(Number) + pterm.Gray("] =>") + pterm.Gray(strings.Repeat(" ", rightPadding)) + pterm.White(Text) + pterm.Gray(strings.Repeat("", leftPadding))

	pterm.Println(formattedText)
}

func Input() {
	pterm.Print(pterm.Gray("[") + pterm.Cyan(">") + pterm.Gray("] "))
}

func Input2() {
	pterm.Print(pterm.Gray("\t\t  [") + pterm.Cyan(">") + pterm.Gray("] "))
}
