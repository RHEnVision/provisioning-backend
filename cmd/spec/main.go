package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/ghodss/yaml"
)

type APISchemaGen struct {
	Components openapi3.Components
	Servers    openapi3.Servers
}

func (s *APISchemaGen) init() {
	s.Servers = openapi3.Servers{
		&openapi3.Server{
			Description: "Local development",
			URL:         "http://0.0.0.0:{port}/api/{applicationName}",
			Variables: map[string]*openapi3.ServerVariable{
				"applicationName": {Default: "provisioning"},
				"port":            {Default: "8000"},
			},
		},
	}
	s.Components = openapi3.NewComponents()
	s.Components.Schemas = make(map[string]*openapi3.SchemaRef)
	s.Components.Responses = make(map[string]*openapi3.ResponseRef)
}

func (s *APISchemaGen) addSchema(name string, model interface{}) {
	schema, err := openapi3gen.NewSchemaRefForValue(model, s.Components.Schemas)
	checkErr(err)
	s.Components.Schemas[name] = schema
}

func (s *APISchemaGen) addResponse(name string, description string, ref string) {
	response := openapi3.NewResponse().WithDescription(description).WithJSONSchemaRef(&openapi3.SchemaRef{Ref: ref})
	s.Components.Responses[name] = &openapi3.ResponseRef{Value: response}
}

func main() {
	gen := APISchemaGen{}
	gen.init()
	// payloads
	gen.addSchema("v1.PubkeyRequest", &payloads.PubkeyRequest{})
	gen.addSchema("v1.PubkeyResponse", &payloads.PubkeyResponse{})
	gen.addSchema("v1.SourceResponse", &payloads.SourceResponse{})
	gen.addSchema("v1.InstanceTypeResponse", &payloads.InstanceTypeResponse{})
	gen.addSchema("v1.ReservationResponse", &payloads.GenericReservationResponsePayload{})
	gen.addSchema("v1.NoopReservationResponse", &payloads.NoopReservationResponsePayload{})
	gen.addSchema("v1.AWSReservationRequest", &payloads.AWSReservationRequestPayload{})
	gen.addSchema("v1.AWSReservationResponse", &payloads.AWSReservationResponsePayload{})

	// error payloads
	gen.addSchema("v1.ResponseError", &payloads.ResponseError{})

	// errors
	gen.addResponse("NotFound", "The specified resource was not found", "#/components/schemas/v1.ResponseError")
	gen.addResponse("InternalError", "The server encountered with an internal error", "#/components/schemas/v1.ResponseError")
	gen.addResponse("BadRequest", "The request's parameters are not sufficient", "#/components/schemas/v1.ResponseError")

	type Swagger struct {
		Components openapi3.Components `json:"components,omitempty" yaml:"components,omitempty"`
		Servers    openapi3.Servers    `json:"servers,omitempty" yaml:"servers,omitempty"`
	}

	swagger := Swagger{}
	swagger.Servers = gen.Servers
	swagger.Components = gen.Components

	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(swagger)
	checkErr(err)

	schema, err := yaml.JSONToYAML(b.Bytes())
	checkErr(err)

	paths, err := ioutil.ReadFile("./cmd/spec/path.yaml")
	checkErr(err)

	b = &bytes.Buffer{}
	b.Write(paths)
	b.Write(schema)

	doc, err := openapi3.NewLoader().LoadFromData(b.Bytes())
	checkErr(err)

	jsonB, err := json.MarshalIndent(doc, "", "  ")
	checkErr(err)
	err = ioutil.WriteFile("./api/openapi.gen.json", jsonB, 0o644) // #nosec G306
	checkErr(err)
	err = ioutil.WriteFile("./api/openapi.gen.yaml", b.Bytes(), 0o644) // #nosec G306
	checkErr(err)
	fmt.Println("Spec was generated successfully")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
