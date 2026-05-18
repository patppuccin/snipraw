package console

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/patppuccin/snipraw/src/consts"
)

func Banner(msg string) string {
	return consts.AppBanner + "\n" + color.New(color.FgGreen).Sprint(msg) + "\n"
}

func logPrint(icon *color.Color, sym, msg string) {
	fmt.Println(icon.Sprint(sym), msg)
}

func Debug(msg string)   { logPrint(color.New(color.FgHiBlack), "(~)", msg) }
func Info(msg string)    { logPrint(color.New(color.FgBlue), "(i)", msg) }
func Warn(msg string)    { logPrint(color.New(color.FgYellow), "(!)", msg) }
func Error(msg string)   { logPrint(color.New(color.FgRed), "(x)", msg) }
func Success(msg string) { logPrint(color.New(color.FgGreen), "(✓)", msg) }
