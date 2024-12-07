package gen

import (
	"log"
	"os"
	"text/template"

	"github.com/stoewer/go-strcase"
	"golang.org/x/text/cases"
	"gopkg.in/yaml.v3"
)


type Endpoint struct{
	Name string
	Description string
	Request string
	Response string
	Example string
}

type Type struct{
	Name string
	Description string
	Fields []Field
	Example string
}

type Field struct{
	Name string
	Type string
	Description string
	Example string
}

type Service struct {
	Namespace string `yaml: "namespace"`
	Name string `yaml: "name"`
	Version string `yaml: "version"`
	Description string `yaml: "description"`
	Repository string `yaml: "repository"`
	Keywords []string `yaml: "keywords"`
	Endpoints []Endpoint 
	Types []Type
 
}


const ApiTpl = `
package api

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

const (
    NS = "{{ .Namespace }}"
    NAME = "{{ .Name }}"
    SERVICE_NAME = NS + NAME

{{ range .Endpoints }}
    //{{ .Description }}
    {{upper .Name }} = "{{ .Name }}"
{{ end }}
{{ range .Endpoints }}
    //{{ .Description }}
    {{upper .Name }}_ENDPOINT = SERVICE_NAME + "." + "{{ .Name }}"
{{ end }}

)
type (
{{ range .Types }}
 //{{ .Description }}
 {{.Name }}struct {
 {{ range .Fields }}
	{{.Name}} {{.Type}} ` + "`json:\"{{snake .Name}}\"`" + `
 {{ end }}
 }
{{ end }}
)
`




func Gen(filepath string) {

	yalmData, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("failed to read YAML file: %v", err)  
	}

	defn := Service{}

	err = yaml.Unmarshal(yalmData, &defn)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML file: %v", err)
	}

	// empJSON, err := json.MarshalIndent(defn, "", "  ")
    // if err != nil {
    //     log.Fatalf(err.Error())
	// 	return
    // }
    // fmt.Println(string(empJSON))
	// strcase.SnakeCase("FooBar")

	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"title": cases.Title,
		"snake": strcase.SnakeCase,
		"upper": strcase.UpperSnakeCase,
	}

	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("titleTest").Funcs(funcMap).Parse(ApiTpl)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, defn)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
}