package helpers

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

func WriteFile(filename string, data []byte) error {
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}
	return nil
}

func CopyPaste(from []byte, to string) error {
	return WriteFile(to, from)
}

func CopyPasteTemplate(from []byte, to string, args interface{}) error {
	tpl := template.New(to)
	tpl, err := tpl.Parse(string(from))
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if err = tpl.Execute(buf, args); err != nil {
		return err
	}

	return WriteFile(to, buf.Bytes())
}
