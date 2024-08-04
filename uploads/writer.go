package uploads

import "os"

type Writer interface {
	WriteChunk([]byte) (int, error)
	// Consume() error
	// ConsumeEntry() error
	Close() error
}

func NewWriter() Writer {
	return &TmpFileWriter{}
}

type TmpFileWriter struct {
	Files        []os.File
	BytesWritten int
}

func (t *TmpFileWriter) WriteChunk(b []byte) (int, error) {
	return 0, nil
}

func (t *TmpFileWriter) Close() error {
	return nil
}
