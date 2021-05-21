package manager

import "os"

type FileWriter struct {
	File os.File
}

func (f *FileWriter) Write(text string) error {
	_, err := f.File.WriteString(text)

	if err != nil {
		return err
	}

	return nil
}
