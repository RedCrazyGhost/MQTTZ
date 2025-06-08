package color

import "fmt"

type Color struct {
	R, G, B uint8
}

func NewColor(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b}
}

func (c Color) Prefix() string {
	return fmt.Sprintf("\033[1;38;2;%d;%d;%dm", c.R, c.G, c.B)
}

func (c Color) Suffix() string {
	return "\033[0m"
}

func (c Color) Text(text string) string {
	return c.Prefix() + text + c.Suffix()
}
