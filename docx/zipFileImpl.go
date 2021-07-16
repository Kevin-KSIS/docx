package docx

import "archive/zip"

//Type for zip files read from disk
type ZipFile struct {
	data *zip.ReadCloser
}

func (d ZipFile) files() []*zip.File {
	return d.data.File
}

func (d ZipFile) close() error {
	return d.data.Close()
}
