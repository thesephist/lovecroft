package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

func useTemplate(name string) (*template.Template, error) {
	file, err := os.Open(fmt.Sprintf("./templates/%s.html", name))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return template.New(name).Parse(string(data))
}
