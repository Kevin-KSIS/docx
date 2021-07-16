package docx

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"
)


type Docx struct {
	files   []*zip.File
	content string
	links   string
	headers map[string]string
	footers map[string]string
	newProperties string
}

func (d *Docx) GetContent() string {
	return d.content
}

func (d *Docx) SetContent(content string) {
	d.content = content
}

func (d *Docx) ReplaceRaw(oldString string, newString string, num int) {
	d.content = strings.Replace(d.content, oldString, newString, num)
}

func (d *Docx) Replace(oldString string, newString string, num int) (err error) {
	oldString, err = encode(oldString)
	if err != nil {
		return err
	}
	newString, err = encode(newString)
	if err != nil {
		return err
	}
	d.content = strings.Replace(d.content, oldString, newString, num)

	return nil
}

func (d *Docx) ReplaceLink(oldString string, newString string, num int) (err error) {
	oldString, err = encode(oldString)
	if err != nil {
		return err
	}
	newString, err = encode(newString)
	if err != nil {
		return err
	}
	d.links = strings.Replace(d.links, oldString, newString, num)

	return nil
}

func (d *Docx) ReplaceHeader(oldString string, newString string) (err error) {
	return replaceHeaderFooter(d.headers, oldString, newString)
}

func (d *Docx) ReplaceFooter(oldString string, newString string) (err error) {
	return replaceHeaderFooter(d.footers, oldString, newString)
}

func (d *Docx) WriteToFile(path string) (err error) {
	var target *os.File
	target, err = os.Create(path)
	if err != nil {
		return
	}
	defer target.Close()
	err = d.Write(target)
	return
}

func (d *Docx) Write(ioWriter io.Writer) (err error) {
	w := zip.NewWriter(ioWriter)
	var writer io.Writer
	for _, file := range d.files {
		var readCloser io.ReadCloser

		writer, err = w.Create(file.Name)
		if err != nil {
			return err
		}
		readCloser, err = file.Open()
		if err != nil {
			return err
		}
		if file.Name == "word/document.xml" {
			writer.Write([]byte(d.content))
		} else if file.Name == "word/_rels/document.xml.rels" {
			writer.Write([]byte(d.links))
		} else if strings.Contains(file.Name, "header") && d.headers[file.Name] != "" {
			writer.Write([]byte(d.headers[file.Name]))
		} else if strings.Contains(file.Name, "footer") && d.footers[file.Name] != "" {
			writer.Write([]byte(d.footers[file.Name]))
		} else {
			writer.Write(streamToByte(readCloser))
		}
	}

	//  modified
	content := generateCustomXML(d.newProperties)
	fp := strings.NewReader(content)

	writer, err = w.Create("docProps/custom.xml")
	_, err = writer.Write(streamToByte(fp))
	if err != nil {
		return err
	}
	//

	w.Close()
	return
}

/* MODIFIED */
func (d *Docx) AddProperty(key string, value string) {
	xmlString := fmt.Sprintf("<%s>%s</%s>", key, value, key)
	d.newProperties = d.newProperties + xmlString
}

func generateCustomXML(newProp string) string {
	return "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n<Properties\n\t" +
		"xmlns=\"http://schemas.openxmlformats.org/officeDocument/2006/extended-properties\"\n\t" +
		"xmlns:vt=\"http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes\">" +
		newProp +
		"</Properties>"
}