package stdio

import (
	"bufio"
	"fmt"
	"os"
)

type IStdio interface {
	Print(...interface{}) (n int, err error)
	Println(a ...interface{}) (n int, err error)
	Printf(format string, a ...interface{}) (n int, err error)
	Input() string
	Error(s string) (n int, err error)
}

var Io IStdio = ConsoleIO{}

func SetIO(io IStdio) {
	Io = io
}

type ConsoleIO struct {
}

func (ConsoleIO) Print(a ...interface{}) (n int, err error) {
	return fmt.Print(a...)
}

func (ConsoleIO) Println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

func (ConsoleIO) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}

func (ConsoleIO) Error(s string) (n int, err error) {
	return fmt.Fprintln(os.Stderr, s)
}

func (c ConsoleIO) Input() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
