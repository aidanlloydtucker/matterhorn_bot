package main

import (
	"html/template"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		panic("error! must have 2 arguments: [command name] [output]")
	}

	args := os.Args[1:]

	tmpl, err := template.ParseFiles("tools/generate_cmd/cmd.tmpl")
	if err != nil {
		panic(err)
	}

	f, err := os.Create(args[1])
	if err != nil {
		panic(err)
	}

	tmpl.Execute(f, struct {
		CommandName     string
		CommandNameCaps string
	}{
		CommandName:     args[0],
		CommandNameCaps: strings.Title(args[0]),
	})
}
