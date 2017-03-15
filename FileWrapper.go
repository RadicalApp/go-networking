package go_networking

// FileWrapper contains all necessary file data for upload.
type fileWrapper struct {
	data        []byte // raw data
	name        string // file name
	contentType string // content type
	path        string
}

func newFileWrapper() fileWrapper {
	return fileWrapper{
		name:        "dusty.dust",
		contentType: "application/octet-stream",
	}
}

// FileWrapperBuilder implements builder pattern to construct a FileWrapper
type FileWrapperBuilder struct {
	fw fileWrapper
}

func NewFileWrapperBuilder() *FileWrapperBuilder {
	return &FileWrapperBuilder{
		fw: newFileWrapper(),
	}
}

func (fwb *FileWrapperBuilder) SetData(data []byte) *FileWrapperBuilder {
	fwb.fw.data = data
	return fwb
}

func (fwb *FileWrapperBuilder) SetName(name string) *FileWrapperBuilder {
	fwb.fw.name = name
	return fwb
}

func (fwb *FileWrapperBuilder) SetContentType(contentType string) *FileWrapperBuilder {
	fwb.fw.contentType = contentType
	return fwb
}

func (fwb *FileWrapperBuilder) SetPath(path string) *FileWrapperBuilder {
	fwb.fw.path = path
	return fwb
}

func (fwb *FileWrapperBuilder) Build() fileWrapper {
	return fwb.fw
}
