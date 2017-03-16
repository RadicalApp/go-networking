package go_networking

import (
	"io"
	"mime/multipart"
	"os"
)

type byteWriter struct {
	w *multipart.Writer
}

func newByteWriter(w io.Writer) *byteWriter {
	return &byteWriter{
		w: multipart.NewWriter(w),
	}
}

func (bw *byteWriter) writeByteField(fieldname string, value []byte, filename string) error {
	p, err := bw.w.CreateFormFile(fieldname, filename)
	if err != nil {
		return err
	}

	_, err = p.Write(value)
	return err
}

func (bw *byteWriter) writeFileField(fieldname string, path string, filename string) error {
	file, err := os.Open(path)
	p, err := bw.w.CreateFormFile(fieldname, filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(p, file)
	return err
}

func (bw *byteWriter) close() {
	bw.w.Close()
}
