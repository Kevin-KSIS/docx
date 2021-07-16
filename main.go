package main

import "docx/docx"

func main() {
	r, err := docx.ReadDocxFile("TestDocument.docx")
	if err != nil {
		panic(err)
	}

	docx1 := r.NewReplaceDocx()
	docx1.AddProperty("xxxx", "vvvvv")
	docx1.WriteToFile("./new_result_1.docx")
	r.Close()
}
