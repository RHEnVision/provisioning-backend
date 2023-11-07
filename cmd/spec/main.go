package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"gopkg.in/yaml.v3"
)

// addPayloads - MAKE SURE THE TYPE HAS JSON/YAML Go STRUCT TAGS (or "map key XXX not found" error occurs)
func addPayloads(gen *APISchemaGen) {
	gen.addSchema("v1.PubkeyRequest", &payloads.PubkeyRequest{})
	gen.addSchema("v1.PubkeyResponse", &payloads.PubkeyResponse{})
	gen.addSchema("v1.SourceResponse", &payloads.SourceResponse{})
	gen.addSchema("v1.InstanceTypeResponse", &payloads.InstanceTypeResponse{})
	gen.addSchema("v1.GenericReservationResponse", &payloads.GenericReservationResponse{})
	gen.addSchema("v1.NoopReservationResponse", &payloads.NoopReservationResponse{})
	gen.addSchema("v1.AWSReservationRequest", &payloads.AWSReservationRequest{})
	gen.addSchema("v1.AWSReservationResponse", &payloads.AWSReservationResponse{})
	gen.addSchema("v1.AzureReservationRequest", &payloads.AzureReservationRequest{})
	gen.addSchema("v1.AzureReservationResponse", &payloads.AzureReservationResponse{})
	gen.addSchema("v1.GCPReservationRequest", &payloads.GCPReservationRequest{})
	gen.addSchema("v1.GCPReservationResponse", &payloads.GCPReservationResponse{})
	gen.addSchema("v1.AvailabilityStatusRequest", &payloads.AvailabilityStatusRequest{})
	gen.addSchema("v1.AccountIDTypeResponse", &payloads.AccountIdentityResponse{})
	gen.addSchema("v1.SourceUploadInfoResponse", &payloads.SourceUploadInfoResponse{})
	gen.addSchema("v1.LaunchTemplatesResponse", &payloads.LaunchTemplateResponse{})

	gen.addSchema("v1.ListSourceResponse", &payloads.SourceListResponse{})
	gen.addSchema("v1.ListPubkeyResponse", &payloads.PubkeyListResponse{})
	gen.addSchema("v1.ListInstaceTypeResponse", &payloads.InstanceTypeListResponse{})
	gen.addSchema("v1.ListGenericReservationResponse", &payloads.GenericReservationListResponse{})
	gen.addSchema("v1.ListLaunchTemplateResponse", &payloads.LaunchTemplateListResponse{})
}

func addExamples(gen *APISchemaGen) {
	gen.addExample("v1.PubkeyRequestExample", PubkeyRequest)
	gen.addExample("v1.PubkeyResponseExample", PubkeyResponse)
	gen.addExample("v1.PubkeyListResponseExample", PubkeyListResponse)
	gen.addExample("v1.SourceListResponseExample", SourceListResponse)
	gen.addExample("v1.SourceUploadInfoAWSResponse", SourceUploadInfoAWSResponse)
	gen.addExample("v1.SourceUploadInfoAzureResponse", SourceUploadInfoAzureResponse)
	gen.addExample("v1.LaunchTemplateListResponse", LaunchTemplateListResponse)
	gen.addExample("v1.AvailabilityStatusRequest", AvailabilityStatusRequest)
	gen.addExample("v1.GenericReservationResponsePayloadSuccessExample", GenericReservationResponsePayloadSuccessExample)
	gen.addExample("v1.GenericReservationResponsePayloadPendingExample", GenericReservationResponsePayloadPendingExample)
	gen.addExample("v1.GenericReservationResponsePayloadFailureExample", GenericReservationResponsePayloadFailureExample)
	gen.addExample("v1.GenericReservationResponsePayloadListExample", GenericReservationResponsePayloadListExample)
	gen.addExample("v1.AwsReservationRequestPayloadExample", AwsReservationRequestPayloadExample)
	gen.addExample("v1.AwsReservationResponsePayloadPendingExample", AwsReservationResponsePayloadPendingExample)
	gen.addExample("v1.AwsReservationResponsePayloadDoneExample", AwsReservationResponsePayloadDoneExample)
	gen.addExample("v1.AzureReservationRequestPayloadExample", AzureReservationRequestPayloadExample)
	gen.addExample("v1.AzureReservationResponsePayloadPendingExample", AzureReservationResponsePayloadPendingExample)
	gen.addExample("v1.AzureReservationResponsePayloadDoneExample", AzureReservationResponsePayloadDoneExample)
	gen.addExample("v1.GCPReservationRequestPayloadExample", GCPReservationRequestPayloadExample)
	gen.addExample("v1.GCPReservationResponsePayloadPendingExample", GCPReservationResponsePayloadPendingExample)
	gen.addExample("v1.GCPReservationResponsePayloadDoneExample", GCPReservationResponsePayloadDoneExample)
	gen.addExample("v1.NoopReservationResponsePayloadExample", NoopReservationResponsePayloadExample)
	gen.addExample("v1.InstanceTypesAWSResponse", InstanceTypesAWSResponse)
	gen.addExample("v1.InstanceTypesAzureResponse", InstanceTypesAzureResponse)
	gen.addExample("v1.InstanceTypesGCPResponse", InstanceTypesGCPResponse)
}

func addParameters(gen *APISchemaGen) {
	gen.addQueryParameter("Limit", LimitQueryParam)
	gen.addQueryParameter("Offset", OffsetQueryParam)
	gen.addQueryParameter("Token", TokenQueryParam)
}

// addErrorSchemas all generic errors, that can be returned.
func addErrorSchemas(gen *APISchemaGen) {
	// error payloads
	gen.addSchema("v1.ResponseError", &payloads.ResponseError{})

	// errors
	gen.addResponse("NotFound", "The requested resource was not found", "#/components/schemas/v1.ResponseError", ResponseNotFoundErrorExample)
	gen.addResponse("InternalError", "The server encountered an internal error", "#/components/schemas/v1.ResponseError", ResponseErrorGenericExample)
	gen.addResponse("BadRequest", "The request's parameters are not valid", "#/components/schemas/v1.ResponseError", ResponseBadRequestErrorExample)
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
	s.Components.Examples = make(map[string]*openapi3.ExampleRef)
	s.Components.Parameters = make(map[string]*openapi3.ParameterRef)
	return s
}

// Schema customizer allowing tagging with description and nullable to work
var enableNullableAndDescriptionOpts = openapi3gen.SchemaCustomizer(
	func(_name string, _t reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {
		if tag.Get("nullable") == "true" {
			schema.Nullable = true
		}
		if desc, ok := tag.Lookup("description"); ok && desc != "-" {
			schema.Description = desc
		}
		return nil
	},
)

func (s *APISchemaGen) addSchema(name string, model interface{}) {
	schema, err := openapi3gen.NewSchemaRefForValue(model, s.Components.Schemas, enableNullableAndDescriptionOpts)
	if err != nil {
		panic(err)
	}
	s.Components.Schemas[name] = schema
}

func (s *APISchemaGen) addResponse(name string, description string, ref string, example interface{}) {
	response := openapi3.NewResponse().WithDescription(description).WithJSONSchemaRef(&openapi3.SchemaRef{Ref: ref})
	response.Content.Get("application/json").Examples = make(map[string]*openapi3.ExampleRef)
	response.Content.Get("application/json").Examples["error"] = &openapi3.ExampleRef{Value: openapi3.NewExample(example)}
	s.Components.Responses[name] = &openapi3.ResponseRef{Value: response}
}

func (s *APISchemaGen) addExample(name string, value interface{}) {
	// verify all fields has both json and yaml struct flags
	rval := reflect.TypeOf(value)
	checkTags(rval)

	example := openapi3.NewExample(value)
	s.Components.Examples[name] = &openapi3.ExampleRef{Value: example}
}

func (s *APISchemaGen) addQueryParameter(name string, value Parameter) {
	checkTags(reflect.TypeOf(value))

	param := &openapi3.Parameter{
		Name:        value.Name,
		In:          value.In,
		Description: value.Description,
		Required:    value.Required,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Default: value.Default,
				Type:    value.Type,
			},
		},
	}
	s.Components.Parameters[name] = &openapi3.ParameterRef{Value: param}
}

//nolint:goerr113
func checkTags(rval reflect.Type) {
	if rval.Kind() == reflect.Array || rval.Kind() == reflect.Slice {
		checkTags(rval.Elem())
		return
	}

	if rval.Kind() != reflect.Struct {
		fmt.Printf("unable to check type %s of kind %s for struct tags, skipped\n", rval.Name(), rval.Kind().String())
		return
	}

	for i := 0; i < rval.NumField(); i++ {
		for _, tagName := range []string{"json", "yaml"} {
			if _, ok := rval.Field(i).Tag.Lookup(tagName); !ok {
				panic(fmt.Errorf("type %s does not have struct flag '%s'", rval.Name(), tagName))
			}
		}
	}
}

func main() {
	// build schema part
	gen := NewSchemaGenerator()
	addErrorSchemas(gen)
	addPayloads(gen)
	addExamples(gen)
	addParameters(gen)

	// store schema part as buffer
	schemasYaml, err := yaml.Marshal(&gen)
	if err != nil {
		panic(err)
	}

	// load endpoints and info part from file
	bufferYAML, err := os.ReadFile("./cmd/spec/path.yaml")
	if err != nil {
		panic(err)
	}

	// append both into single schema
	bufferYAML = append(bufferYAML, schemasYaml...)

	// load full schema
	loadedSchema, err := openapi3.NewLoader().LoadFromData(bufferYAML)
	if err != nil {
		panic(err)
	}

	// update version in the full schema and store it again
	if len(os.Args) >= 2 {
		loadedSchema.Info.Version = os.Args[1]
		bufferYAML, err = yaml.Marshal(&loadedSchema)
		if err != nil {
			panic(err)
		}
	}

	// validate it
	err = loadedSchema.Validate(context.Background())
	if err != nil {
		panic(err)
	}

	// and store the full schema as JSON and YAML
	bufferJSON, err := json.MarshalIndent(loadedSchema, "", "  ")
	if err != nil {
		panic(err)
	}
	tmp := make([]byte, len(bufferJSON), len(bufferJSON)+1)
	copy(tmp, bufferJSON)
	tmp = append(tmp, '\n')
	bufferJSON = tmp

	err = os.WriteFile("./api/openapi.gen.json", bufferJSON, 0o644) // #nosec G306
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./api/openapi.gen.yaml", bufferYAML, 0o644) // #nosec G306
	if err != nil {
		panic(err)
	}
}
