package uploads

import (
	"fmt"
	"os"
)

type Writer interface {
	WriteChunk(string, []byte) (int, error)
	Consume(string, func(path string)) error
}

func NewTmpWriter() Writer {
	return &TmpFileWriter{
		Files: make(map[string]*os.File),
	}
}

type TmpFileWriter struct {
	Files map[string]*os.File
}

func (t *TmpFileWriter) WriteChunk(ref string, b []byte) (int, error) {
	tmpFile, exists := t.Files[ref]
	if !exists {
		var err error
		tmpFile, err = os.CreateTemp("", "upload_*.tmp")
		if err != nil {
			return 0, err
		}
		t.Files[ref] = tmpFile
	}

	n, err := tmpFile.Write(b)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (t *TmpFileWriter) Consume(ref string, f func(path string)) error {
	file, exists := t.Files[ref]
	if !exists {
		return fmt.Errorf("file not found")
	}

	f(file.Name())

	err := file.Close()
	if err != nil {
		return err
	}

	err = os.Remove(file.Name())
	if err != nil {
		return err
	}

	delete(t.Files, ref)

	return nil
}
