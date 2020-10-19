package main

import (
	"bufio"
	"flag"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kjasir/spec"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	sourcePtr := flag.String("file", "", "filename")
	tmplPtr := flag.String("template", "", "template")
	flag.Parse()

	file, sourcePtrErr := ioutil.ReadFile(*sourcePtr)
	check(sourcePtrErr)

	fileByte := []byte(string(file))
	swagger, swaggerErr := openapi3.NewSwaggerLoader().LoadSwaggerFromData(fileByte)
	check(swaggerErr)
	design := spec.Transform(swagger)

	tmpl := template.Must(template.ParseFiles(*tmplPtr))
	for _, resource := range design.Resources {
		filename := strings.ReplaceAll(resource.ResourceDefinition, " ", "_")
		fileTarget, fileTargetErr := os.Create(filename + ".md")
		check(fileTargetErr)
		defer fileTarget.Close()
		writer := bufio.NewWriter(fileTarget)
		execErr := tmpl.Execute(writer, resource)
		check(execErr)
		writer.Flush()
	}
}
