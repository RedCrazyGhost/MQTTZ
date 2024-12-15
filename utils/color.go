package utils

import "fmt"

func fontColor(r, g, b uint8) string {
	return fmt.Sprintf("\033[1;38;2;%d;%d;%dm", r, g, b)
}

func UseColor(text string, r, g, b uint8) string {
	return fontColor(r, g, b) + text + "\033[0m"
}
