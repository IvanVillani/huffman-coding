package manager

import "fmt"

type ConsoleWriter struct{}

func (f *ConsoleWriter) Write(text string) error {
	fmt.Print(text)

	return nil
}
