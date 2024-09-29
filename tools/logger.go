package tools

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func PrintLogo() {
	levelColor := color.New(color.FgCyan)
	levelColor.Println(`
 /$$$$$$$  /$$                         /$$$$$$$              /$$    
| $$__  $$| $$                        | $$__  $$            | $$    
| $$  \ $$| $$ /$$   /$$ /$$$$$$/$$$$ | $$  \ $$  /$$$$$$  /$$$$$$  
| $$$$$$$ | $$| $$  | $$| $$_  $$_  $$| $$$$$$$  /$$__  $$|_  $$_/  
| $$__  $$| $$| $$  | $$| $$ \ $$ \ $$| $$__  $$| $$  \ $$  | $$    
| $$  \ $$| $$| $$  | $$| $$ | $$ | $$| $$  \ $$| $$  | $$  | $$ /$$
| $$$$$$$/| $$|  $$$$$$/| $$ | $$ | $$| $$$$$$$/|  $$$$$$/  |  $$$$/
|_______/ |__/ \______/ |__/ |__/ |__/|_______/  \______/    \___/  
`)

	levelColor.Println("ρσωєяє∂ ву: ѕкιвι∂ι ѕιgмα ¢σ∂є")

	levelColor = color.New(color.FgRed)
	levelColor.Println("[!] All risks are your responsibility. This tool is intended for educational purposes and to make your life easier.....")
}

func Logger(level, message string) {
	message = strings.TrimSpace(strings.ReplaceAll(message, "\n", ""))

	level = strings.ToLower(level)
	var levelColor *color.Color

	switch level {
	case "info":
		levelColor = color.New(color.FgWhite)
	case "error":
		levelColor = color.New(color.FgRed)
	case "success":
		levelColor = color.New(color.FgGreen)
	case "warning":
		levelColor = color.New(color.FgYellow)
	default:
		levelColor = color.New(color.FgWhite)
	}

	if level == "input" {
		levelColor.Printf("[+] %s", message)
	} else if level == "error" || level == "warning" {
		levelColor.Println(fmt.Sprintf("[!] %s", message))
	} else if _, err := strconv.Atoi(level); err == nil {
		levelColor.Println(fmt.Sprintf("[%s] %s", level, message))
	} else {
		levelColor.Println(fmt.Sprintf("[*] %s", message))
	}
}
