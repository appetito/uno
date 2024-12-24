package gen

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"
	"golang.org/x/text/cases"
	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	Name        string
	Description string
	Request     string
	Response    string
	Example     string
}

type Type struct {
	Name        string
	Description string
	Fields      []Field
	Example     string
}

type Field struct {
	Name        string
	Type        string
	Description string
	Example     string
}

type Service struct {
	Namespace   string     `yaml:"namespace"`
	Name        string     `yaml:"name"`
	Module      string     `yaml:"module"`
	Version     string     `yaml:"version"`
	Description string     `yaml:"description"`
	Repository  string     `yaml:"repository"`
	Keywords    []string   `yaml:"keywords"`
	Endpoints   []Endpoint
	Types       []Type
}

const ApiTpl = `
package api

import (
    "encoding/json"
    "context"

    "github.com/appetito/uno"
    "github.com/nats-io/nats.go"
)

const (
    NS = "{{ .Namespace }}"
    NAME = "{{ .Name }}"
    SERVICE_NAME = NS + "." + NAME

{{ range .Endpoints }}
    //{{ .Description }}
    {{upper .Name }} = "{{ .Name }}"
{{ end }}
{{ range .Endpoints }}
    //{{ .Description }}
    {{upper .Name }}_ENDPOINT = SERVICE_NAME + "." + {{upper .Name }}
{{ end }}

)
type (
{{ range .Types }}
 //{{ .Description }}
 {{.Name }} struct {
{{- range .Fields }}
    {{.Name}} {{.Type}} ` + "`json:\"{{snake .Name}}\"`" + `
{{- end}}
 }
{{ end }}
)


type {{ camel .Name }}Client struct {
    uc *uno.UnoClient
}


func New{{ camel .Name }}Client(nc *nats.Conn, cfg *uno.UnoClientConfig) *{{ camel .Name }}Client {
    return &{{ camel .Name }}Client{
        uc: uno.NewUnoClient(nc, cfg),
    }
}


{{ range .Endpoints }}
//{{ .Description }}
func (c *{{ camel $.Name }}Client) {{ .Name }}(ctx context.Context, request {{.Request}}) (response {{.Response}}, err error) {
    reply, err := c.uc.RequestJSON(ctx, {{upper .Name }}_ENDPOINT, request)
    if err != nil {
        return response, err
    }
    err = json.Unmarshal(reply.Data, &response)
    return response, err
}
{{ end }}
`


const ServiceTpl = `
package service

import (

    "github.com/appetito/uno"
    "github.com/nats-io/nats.go"

	"github.com/rs/zerolog/log"

	"{{.Module }}/api"
	"{{.Module }}/internal/config"
	"{{.Module }}/internal/handlers"

)

func New(cfg *config.Config) uno.Service {

	log.Info().Str("URL", cfg.NatsServers).Msg("Connecting to NATS")
	nc, err := nats.Connect(cfg.NatsServers)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}

	svc, err := uno.AddService(nc, uno.Config{
		Name:       "{{.Namespace}}" + "_" +  "{{.Name}}",
		Version:     "{{.Version}}",
		Description: "{{.Name}}",
		Interceptors: []uno.InterceptorFunc{
			uno.NewPanicInterceptor,
			uno.NewMetricsInterceptor,
			uno.NewTracingInterceptor, 
			uno.NewLoggingInterceptor,   
		},
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to add service")
	}

	root := svc.AddGroup(api.SERVICE_NAME)

{{range .Endpoints }}
	root.AddEndpoint(api.{{upper .Name}}, uno.AsStructHandler[{{apiref .Request}}](handlers.{{.Name}}Handler))
{{ end }}
	
	return svc
}
`

const HandlersTpl = `
package handlers

import (

    "github.com/appetito/uno"

	"{{.Module }}/api"

)

{{ range .Endpoints }}
//{{ .Description }}
func {{ .Name }}Handler(r uno.Request, request {{apiref .Request}}){
	var response {{apiref .Response}}
	r.RespondJSON(response)
}
{{ end }}
`

const ConfigTpl = `
package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	NatsServers string ` + "`" + `env:"NATS_SERVERS" env-default:"localhost:4222"` + "`" + `
}



func GetConfig() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	return &cfg
}
`

const CmdMainTpl = `
package main

import (

	"{{.Module }}/internal/config"
	"{{.Module }}/internal/service"
)

func main(){
	cfg := config.GetConfig()
	svc := service.New(cfg)
	svc.ServeForever()
}
`

func GenProject(filepath string) {

	yalmData, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("failed to read YAML file: %v", err)
	}

	defn := Service{}

	err = yaml.Unmarshal(yalmData, &defn)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML file: %v", err)
	}

	CreateGoModFile(defn)
	CreatePackage("api", "api", ApiTpl, defn, false)
	CreatePackage("internal/handlers", "handlers", HandlersTpl, defn, false)
	CreatePackage("internal/service", "service", ServiceTpl, defn, false)
	CreatePackage("internal/config", "config", ConfigTpl, defn, false)
	CreatePackage("cmd", "main", CmdMainTpl, defn, false)
}


func GenAPI(filepath string) {

	yalmData, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("failed to read YAML file: %v", err)
	}

	defn := Service{}

	err = yaml.Unmarshal(yalmData, &defn)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML file: %v", err)
	}

	CreatePackage("api", "api", ApiTpl, defn, true)
}


func CreatePackage(dir, name string, tpl string, defn Service, overwrite bool) {
	funcMap := template.FuncMap{
		"title": cases.Title,
		"snake": strcase.SnakeCase,
		"upper": strcase.UpperSnakeCase,
		"camel": strcase.UpperCamelCase,
		"apiref": AsApiRef,
	}

	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("api").Funcs(funcMap).Parse(tpl)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	pkgDir := dir
	pkgFilePath := dir + "/" + name + ".go"

	if isDirExist(pkgDir) && !overwrite {
		log.Fatalf("package %s already exists", name)
	}

	err = os.MkdirAll(pkgDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new file
	f, err := os.Create(pkgFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Run the template to verify the output.
	err = tmpl.Execute(f, defn)
	if err != nil {
		log.Fatalf("%s package gen: execution error: %s", name, err)
	}
	log.Println(name + " package created")
}


func CreateGoModFile(defn Service) {
	if isDirExist("go.mod") {
		log.Fatalf("go.mod already exists")
	}
	// Create a new file
	f, err := os.Create("go.mod")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Run the template to verify the output.
	_, err = f.WriteString("module "+ defn.Module + "\n")
	if err != nil {
		log.Fatalf("go.mod gen: execution error: %s", err)
	}
	log.Println("go.mod created")
}


func AsApiRef(s string) string {
	if strings.HasPrefix(s, "[]") {
		return fmt.Sprintf("[]api.%s", s[2:])
	}
	return fmt.Sprintf("api.%s", s)
}


func isDirExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}