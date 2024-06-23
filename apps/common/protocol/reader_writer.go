package protocol

import "io"

type ReaderWriter interface {
	Reader
	Writer
}

type defaultReaderWriter struct {
	reader *Reader
	writer *Writer
}

func (d *defaultReaderWriter) ReadHeader() (Header, error) {
	return (*d.reader).ReadHeader()
}

func (d *defaultReaderWriter) Read() (payload []byte, err error) {
	return (*d.reader).Read()
}

func (d *defaultReaderWriter) Write(payload []byte) error {
	return (*d.writer).Write(payload)
}

func NewReaderWriter(reader io.Reader, writer io.Writer) *ReaderWriter {

	var rw ReaderWriter = &defaultReaderWriter{
		reader: NewReader(reader),
		writer: NewWriter(writer),
	}

	return &rw
}
