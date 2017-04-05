package go_networking

import (
	"io"
	"mime/multipart"
	"os"
)

// `byteWriter` a wrapper struct to read the progress of bytes written.
// Implementation for `Writer`
type byteWriter struct {
	w *multipart.Writer
}

// Constructor for `byteWriter`
func newByteWriter(w io.Writer) *byteWriter {
	return &byteWriter{
		w: multipart.NewWriter(w),
	}
}

// Create and write a byte field
func (bw *byteWriter) writeByteField(fieldname string, value []byte, filename string) error {
	p, err := bw.w.CreateFormFile(fieldname, filename)
	if err != nil {
		return err
	}

	_, err = p.Write(value)
	return err
}

// Create and write a form field
func (bw *byteWriter) writeFileField(fieldname string, path string, filename string) error {
	file, err := os.Open(path)
	p, err := bw.w.CreateFormFile(fieldname, filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(p, file)
	return err
}

// Close finishes the multipart message and writes the trailing
// boundary end line to the output.
func (bw *byteWriter) close() {
	bw.w.Close()
}
