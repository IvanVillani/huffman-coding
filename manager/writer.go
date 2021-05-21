package manager

type Writer interface {
	Write(text string) error
}
