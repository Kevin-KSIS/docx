package docx

import "archive/zip"

//Contains functions to work with data from a zip file
type ZipData interface {
	files() []*zip.File
	close() error
}
