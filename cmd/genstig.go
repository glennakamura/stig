// Generates STIG code from templates.
package main

//go:generate go run gentmpl.go -src=templates -dst=tmpl
//go:generate go fmt tmpl/templates.go

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"stig/cmd/tmpl"
	"stig/xccdf"
	"strings"
	"text/template"
)

func wrapLine(line string, width int) string {
	var text string
	for len(line) > width {
		index := strings.LastIndex(line[:width+1], " ")
		if index < 0 {
			if index = strings.Index(line, " "); index < 0 {
				break
			}
		}
		text += strings.TrimRight(line[:index], " ") + "\n"
		line = line[index+1:]
	}
	return text + strings.TrimRight(line, " ")
}

func wrapText(text string, width int) string {
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		lines = append(lines, wrapLine(line, width))
	}
	return strings.Join(lines, "\n")
}

func main() {
	var ext, lang string
	var width int
	flag.StringVar(&ext, "ext", ".sh", "output file extension")
	flag.StringVar(&lang, "lang", "bash", "generated code language")
	flag.IntVar(&width, "width", 70, "line wrap width")
	flag.Parse()

	for _, stig := range flag.Args() {
		benchmarkXML, err := ioutil.ReadFile(stig)
		if err != nil {
			log.Fatal(err)
		}

		var benchmark xccdf.Benchmark
		if err = xml.Unmarshal(benchmarkXML, &benchmark); err != nil {
			log.Fatal(err)
		}

		funcMap := template.FuncMap{
			"benchmark": func() xccdf.Benchmark {
				return benchmark
			},
			"replace": strings.Replace,
			"wrap": func(text string) string {
				return wrapText(text, width)
			},
		}
		stigTemplate := template.Must(
			template.New(lang).Funcs(funcMap).Parse(
				tmpl.Templates[lang],
			),
		)

		for _, group := range benchmark.Groups {
			file, err := os.Create(group.Id + ext)
			if err != nil {
				log.Fatal(err)
			}
			stigTemplate.Execute(file, group)
			file.Close()
		}
	}
}
