// +build ignore

//generate command: go run gen/generate.go --html views/*.html --tmpl=views/setting.tmpl --output=setting --package=setting
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

//var -
var (
	KEYWORD = "$$HTML$$"
	NAME    = "$$NAME$$"
	ID      = "$$ID$$"
	PACKAGE = "$$PACKAGE$$"
	TMPLS   = "$$TMPLS$$"
)

//options - parameter
type options struct {
	Input  string `long:"html" description:"html file"`
	Output string `long:"output" description:"output go file"`
	Tmpl   string `long:"tmpl" description:"template file"`
	Pkg    string `long:"package" description:"go package name"`
}

func main() {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(-1)
	}
	if opts.Pkg == "" {
		opts.Pkg = "main"
	}

	//read template
	var tmplsrc string = generateTemplate
	if opts.Tmpl != "" {
		t, err := ioutil.ReadFile(opts.Tmpl)
		if err != nil {
			log.Printf("Cannot open %q: %v\n", opts.Tmpl, err)
			os.Exit(-1)
		}
		tmplsrc = string(t)
	}
	tmplsrc = strings.ReplaceAll(tmplsrc, `\x60`, "`")

	var summary string
	var dir string
	if !strings.HasSuffix(opts.Output, ".go") {
		dir = opts.Output
		summary = filepath.Join(dir, "template.go")
	} else {
		//check output directory
		dir = filepath.Dir(opts.Output)
		if !Exists(dir) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				log.Println("Failed to create output directory", err)
				os.Exit(-1)
			}
		}
		summary = filepath.Join(dir, "template.go")
	}

	filenames := make([]string, 0)
	//glob
	if strings.Index(opts.Input, "*") > -1 || strings.Index(opts.Input, "?") > -1 {
		names, err := filepath.Glob(opts.Input)
		if err != nil {
			log.Printf("Cannot find input file. %v", err)
			os.Exit(-1)
		}
		if len(names) == 0 {
			log.Printf("template: pattern matches no files: %#q", opts.Input)
			os.Exit(-1)
		}
		filenames = append(filenames, names...)
	} else {
		filenames = append(filenames, opts.Input)
	}

	tmpls := ""
	for _, filename := range filenames {
		//read html
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Cannot open %q: %v\n", filename, err)
			os.Exit(-1)
		}
		tmpl := tmplsrc

		//Package
		tmpl = strings.ReplaceAll(tmpl, PACKAGE, opts.Pkg)

		base := Name(filename)
		title := Title(base)
		tmpl = strings.ReplaceAll(tmpl, ID, title)

		//define templates
		dtName := title + "Name"
		dtTmpl := title + "Template"
		tmpls = tmpls + defineTmpl(dtName, dtTmpl)

		//Name
		name := filepath.Base(filename)
		tmpl = strings.ReplaceAll(tmpl, NAME, name)

		//html body
		tmpl = strings.ReplaceAll(tmpl, KEYWORD, string(b))

		gofile := base + ".go"
		output := filepath.Join(dir, gofile)
		//write to go file
		err = ioutil.WriteFile(output, []byte(tmpl), 0644)
		if err != nil {
			log.Printf("Cannot write %v: %v\n", output, err)
			os.Exit(-1)
		}
	}

	if len(tmpls) > 0 {
		tmpl := generateSummary
		//Package
		tmpl = strings.ReplaceAll(tmpl, PACKAGE, opts.Pkg)
		//define templates
		tmpl = strings.ReplaceAll(tmpl, TMPLS, tmpls)

		err := ioutil.WriteFile(summary, []byte(tmpl), 0644)
		if err != nil {
			log.Printf("Cannot write %v: %v\n", summary, err)
			os.Exit(-1)
		}
	}
}

//Exists -
func Exists(name string) bool {
	_, err := os.Stat(name)

	return !os.IsNotExist(err)
}

//Name -
func Name(name string) string {
	fileName := filepath.Base(name)
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

//Title -
func Title(name string) string {
	return strings.Title(strings.ToLower(name))
}

var generateTemplate = `
// Code generated .* DO NOT EDIT\.$
package $$PACKAGE$$

//Name - page Name
var $$ID$$Name string = \x60$$NAME$$\x60

//Template - page template
var $$ID$$Template string = \x60$$HTML$$\x60
`

func defineTmpl(name, data string) string {
	return fmt.Sprintf("{Name: %s, Data: %s},", name, data)
}

var generateSummary = `// Code generated .* DO NOT EDIT\.$
package $$PACKAGE$$

import (
	"text/template"
)

//DefineTemplate -
type DefineTemplate struct {
	Name string
	Data string
}

//NewTemplate - page template
func NewTemplate() (*template.Template, error) {
	tmpls := []DefineTemplate{
		$$TMPLS$$
	}

	var t *template.Template
	var w *template.Template
	for _, d := range tmpls {
		if t == nil {
			t = template.New(d.Name)
		}
		if d.Name == t.Name() {
			w = t
		} else {
			w = t.New(d.Name)
		}
		_, err := w.Parse(d.Data)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
`
