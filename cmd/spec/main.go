package main

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"gopkg.in/yaml.v2"
)

// addPayloads - MAKE SURE THE TYPE HAS JSON/YAML Go STRUCT TAGS (or "map key XXX not found" error occurs)
func addPayloads(gen *APISchemaGen) {
	gen.addSchema("v1.PubkeyRequest", &payloads.PubkeyRequest{})
	gen.addSchema("v1.PubkeyResponse", &payloads.PubkeyResponse{})
	gen.addSchema("v1.SourceResponse", &payloads.SourceResponse{})
	gen.addSchema("v1.InstanceTypeResponse", &payloads.InstanceTypeResponse{})
	gen.addSchema("v1.ReservationResponse", &payloads.GenericReservationResponsePayload{})
	gen.addSchema("v1.NoopReservationResponse", &payloads.NoopReservationResponsePayload{})
	gen.addSchema("v1.AWSReservationRequest", &payloads.AWSReservationRequestPayload{})
	gen.addSchema("v1.AWSReservationResponse", &payloads.AWSReservationResponsePayload{})
	gen.addSchema("v1.AzureReservationRequest", &payloads.AzureReservationRequestPayload{})
	gen.addSchema("v1.AzureReservationResponse", &payloads.AzureReservationResponsePayload{})
	gen.addSchema("v1.AvailabilityStatusRequest", &payloads.AvailabilityStatusRequest{})
	gen.addSchema("v1.AccountIDTypeResponse", &payloads.AccountIdentityResponse{})
	gen.addSchema("v1.SourceUploadInfoResponse", &payloads.SourceUploadInfoResponse{})
	gen.addSchema("v1.LaunchTemplatesResponse", &payloads.LaunchTemplateResponse{})
}

// addErrorSchemas all generic errors, that can be returned.
func addErrorSchemas(gen *APISchemaGen) {
	// error payloads
	gen.addSchema("v1.ResponseError", &payloads.ResponseError{})

	// errors
	gen.addResponse("NotFound", "The requested resource was not found", "#/components/schemas/v1.ResponseError")
	gen.addResponse("InternalError", "The server encountered an internal error", "#/components/schemas/v1.ResponseError")
	gen.addResponse("BadRequest", "The request's parameters are not valid", "#/components/schemas/v1.ResponseError")
}

type APISchemaGen struct {
	Components openapi3.Components `json:"components,omitempty" yaml:"components,omitempty"`
	Servers    openapi3.Servers    `json:"servers,omitempty" yaml:"servers,omitempty"`
}

func NewSchemaGenerator() *APISchemaGen {
	s := &APISchemaGen{}
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

	return s
}

// Schema customizer allowing tagging with nullable to work
var enableNullableOpt = openapi3gen.SchemaCustomizer(
	func(_name string, _t reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {
		if tag.Get("nullable") == "true" {
			schema.Nullable = true
		}
		return nil
	},
)

func (s *APISchemaGen) addSchema(name string, model interface{}) {
	schema, err := openapi3gen.NewSchemaRefForValue(model, s.Components.Schemas, enableNullableOpt)
	if err != nil {
		panic(err)
	}
	s.Components.Schemas[name] = schema
}

func (s *APISchemaGen) addResponse(name string, description string, ref string) {
	response := openapi3.NewResponse().WithDescription(description).WithJSONSchemaRef(&openapi3.SchemaRef{Ref: ref})
	s.Components.Responses[name] = &openapi3.ResponseRef{Value: response}
}

func main() {
	gen := NewSchemaGenerator()
	addErrorSchemas(gen)
	addPayloads(gen)

	bufferYAML, err := os.ReadFile("./cmd/spec/path.yaml")
	if err != nil {
		panic(err)
	}
	schemasYaml, err := yaml.Marshal(&gen)
	if err != nil {
		panic(err)
	}
	bufferYAML = append(bufferYAML, schemasYaml...)

	// Load the full yaml schema to dump it to JSON with indent as a whole.
	// This also helps validate the spec.
	loadedSchema, err := openapi3.NewLoader().LoadFromData(bufferYAML)
	if err != nil {
		panic(err)
	}

	bufferJSON, err := json.MarshalIndent(loadedSchema, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./api/openapi.gen.json", bufferJSON, 0o644) // #nosec G306
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./api/openapi.gen.yaml", bufferYAML, 0o644) // #nosec G306
	if err != nil {
		panic(err)
	}
}
