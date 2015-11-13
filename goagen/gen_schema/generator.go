package genschema

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/raphael/goa/design"
	"github.com/raphael/goa/goagen/codegen"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Generator is the application code generator.
type Generator struct {
	genfiles []string
}

// Generate is the generator entry point called by the meta generator.
func Generate(api *design.APIDefinition) ([]string, error) {
	g, err := NewGenerator()
	if err != nil {
		return nil, err
	}
	return g.Generate(api)
}

// NewGenerator returns the application code generator.
func NewGenerator() (*Generator, error) {
	app := kingpin.New("Main generator", "application JSON schema generator")
	codegen.RegisterFlags(app)
	NewCommand().RegisterFlags(app)
	_, err := app.Parse(os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf(`invalid command line: %s. Command line was "%s"`,
			err, strings.Join(os.Args, " "))
	}
	return new(Generator), nil
}

// JSONSchemaDir is the path to the directory where the schema controller is generated.
func JSONSchemaDir() string {
	return filepath.Join(codegen.OutputDir, "schema")
}

// Generate produces the skeleton main.
func (g *Generator) Generate(api *design.APIDefinition) ([]string, error) {
	os.RemoveAll(JSONSchemaDir())
	os.MkdirAll(JSONSchemaDir(), 0755)
	s := APISchema(api)
	js, err := s.JSON()
	if err != nil {
		return nil, err
	}
	schemaFile := filepath.Join(JSONSchemaDir(), "schema.json")
	if err := ioutil.WriteFile(schemaFile, js, 0755); err != nil {
		return nil, err
	}
	g.genfiles = append(g.genfiles, schemaFile)
	controllerFile := filepath.Join(JSONSchemaDir(), "schema.go")
	tmpl, err := template.New("schema").Parse(jsonSchemaTmpl)
	if err != nil {
		panic(err.Error()) // bug
	}
	gg := codegen.NewGoGenerator(controllerFile)
	imports := []*codegen.ImportSpec{
		codegen.SimpleImport("github.com/raphael/goa"),
	}
	gg.WriteHeader(fmt.Sprintf("%s JSON Hyper-schema", api.Name), "schema", imports)
	data := map[string]interface{}{
		"schema": string(js),
	}
	err = tmpl.Execute(gg, data)
	if err != nil {
		g.Cleanup()
		return nil, err
	}
	if err := gg.FormatCode(); err != nil {
		g.Cleanup()
		return nil, err
	}
	g.genfiles = []string{controllerFile, schemaFile}
	return g.genfiles, nil
}

// Cleanup removes all the files generated by this generator during the last invokation of Generate.
func (g *Generator) Cleanup() {
	for _, f := range g.genfiles {
		os.Remove(f)
	}
	g.genfiles = nil
}

const jsonSchemaTmpl = `
// MountController mounts the API JSON schema controller under "/schema.json".
func MountController(app *goa.Application) {
	app.Logger.Info("mount", "ctrl", "Schema", "action", "Show", "route", "GET /schema.json")
	app.Router.Handle("GET", "/schema.json", app.NewHTTPRouterHandle("Schema", "Show", getSchema))
}

// getSchema is the httprouter handle that returns the API JSON hyper schema.
// func getSchema(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
func getSchema(ctx *goa.Context) error {
	ctx.Header().Set("Content-Type", "application/schema+json")
	ctx.Header().Set("Cache-Control", "public, max-age=3600")
	return ctx.Respond(200, []byte(schema))
}

// Generated schema
const schema = ` + "`" + `{{.schema}}` + "`" + `
`