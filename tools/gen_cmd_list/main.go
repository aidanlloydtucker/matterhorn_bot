package main

import (
	"go/parser"
	"go/token"
	"go/ast"
	"os"
	"strings"
	"text/template"
	"github.com/urfave/cli"
	"errors"
	"bytes"
	"go/format"
)

type Handler struct {
	Name string
	Package string
}

func main() {
	app := cli.NewApp()

	app.Name = "Generate Command List"
	app.Usage = "Generates a file that initializes commands for the bot"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "out, o",
			Usage: "Output filepath for the command list",
			Value: "command_list.go",
		},
		cli.StringSliceFlag{
			Name: "dir, d",
			Usage: "Directories that are parsed",
			Value: &cli.StringSlice{"commands"},
		},
		cli.StringSliceFlag{
			Name: "package, p",
			Usage: "Packages to import in the file",
		},
	}

	app.Action = runApp
	app.Run(os.Args)
}

func runApp(c *cli.Context) error {
	if !c.IsSet("out") {
		return errors.New("Missing output file")
	}

	pkgDirList := c.StringSlice("dir")

	fset := token.NewFileSet()

	// List of handlers
	handlers := []Handler{}

	for _, pkgDir := range pkgDirList {
		// Parse the directory for packages
		pkgs, err := parser.ParseDir(fset, pkgDir, filterGoFiles, 0)
		if err != nil {
			return err
		}

		for _, pkg := range pkgs {
			// Goes through all files
			for _, f := range pkg.Files {
				// Gets all declarations in the file
				for _, d := range f.Decls {
					// Makes sure they are all type declarations
					gd, ok := d.(*ast.GenDecl)
					if !ok {
						continue
					}
					if gd.Tok != token.TYPE {
						continue
					}
					for _, spec := range gd.Specs {
						sp, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}
						if strings.HasSuffix(strings.ToLower(sp.Name.String()), "handler") {
							handlers = append(handlers, Handler{
								Name: sp.Name.String(),
								Package: pkg.Name,
							})
						}
					}
				}
			}
		}
	}

	tmpl, err := template.ParseFiles("tools/gen_cmd_list/cmd_list.tmpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, struct {
		Packages     []string
		Handlers []Handler
		Args string
	}{
		Packages:     c.StringSlice("packages"),
		Handlers: handlers,
		Args: strings.Join(os.Args[1:], " "),
	})
	if err != nil {
		return err
	}

	fmtTmpl, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	outFile, err := os.Create(c.String("out"))
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = outFile.Write(fmtTmpl)
	return err
}

func filterGoFiles(fi os.FileInfo) bool{
	return !strings.HasSuffix(fi.Name(), "_test.go")
}