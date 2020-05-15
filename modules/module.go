package modules

import (
	"fmt"
)

type Module interface {
	Output() string
}

const (
	goodColor = "#69f0ae"
	badColor  = "#ff5555"
)

func SimpleOutput(fullText string) string {
	return fmt.Sprintf("{\"full_text\":\" %s \"}", fullText)
}

func ColoredOutput(fullText string, color string) string {
	return fmt.Sprintf("{\"full_text\":\" %s \",\"color\":\"%s\"}", fullText, color)
}

func GoodOutput(fullText string) string {
	return ColoredOutput(fullText, goodColor)
}

func BadOutput(fullText string) string {
	return ColoredOutput(fullText, badColor)
}
