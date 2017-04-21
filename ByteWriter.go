package go_networking

import (
	"io"
	"mime/multipart"
	"os"
	"net/textproto"
	"fmt"
	"strings"
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
func (bw *byteWriter) writeByteField(fieldname string, value []byte, filename string, contentType string) error {
	p, err := bw.CreateFormFileWithContentType(fieldname, filename, contentType)
	if err != nil {
		return err
	}
	_, err = p.Write(value)
	return err
}

// Create and write a form field
func (bw *byteWriter) writeFileField(fieldname string, path string, filename string, contentType string) error {
	file, err := os.Open(path)
	p, err := bw.CreateFormFileWithContentType(fieldname, filename, contentType)
	if err != nil {
		return err
	}

	_, err = io.Copy(p, file)
	return err
}

// CreateFormFile is a convenience wrapper around CreatePart. It creates
// a new form-data header with the provided field name, file name, and content type.
func (bw *byteWriter) CreateFormFileWithContentType(fieldname, filename, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
	fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
	escapeQuotes(fieldname), escapeQuotes(filename)))
	if len(contentType) > 0 {
		h.Set("Content-Type", contentType)
	} else {
		h.Set("Content-Type", "application/octet-stream")
	}

	return bw.w.CreatePart(h)
}

// Close finishes the multipart message and writes the trailing
// boundary end line to the output.
func (bw *byteWriter) close() {
	bw.w.Close()
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
