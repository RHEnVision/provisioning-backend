# API Documentation

This service uses OpenAPI for API documentation.

We use [redoc](https://github.com/Redocly/redoc) for auto generating swagger UI based on the spec file
The docs locate under `<root>/docs` and the json spec under `<root>/openapi.json`
In addition, `openapi.json` serves under `/api/provisioning/v1/openapi.json`

## Adding new endpoint

We use a hybrid approach when schemas are generated from Go code and paths are manually maintained.

### Schemas

1. Create the endpoint's payload under `internal/payload`
2. Create an example(s) in `cmd/spec` package
3. Register the type and the example(s) in `cmd/spec/main.go` application

### Paths

Edit `/cmd/spec/path.yaml` utilizing the generated schemas and examples.

```yml
paths:
  /<NEW_ENDPOINT>:
    get:
      operationId: getResourceList
      description: ''
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/v1.MySchema'
              examples:
                  example1:
                      $ref: '#/components/examples/v1.MySchemaExample1'
                  example2:
                      $ref: '#/components/examples/v1.MySchemaExample2'
```

Operation naming convention:

* GET /resource - getResourceList
* POST /resource - createResource
* GET /resource/ID - getResourceById
* DELETE /resource/ID - removeResourceById

Make sure to assign OpenAPI "tags" to each endpoint.

### Errors

You can reuse and reference predefined error responses:

```yaml
 responses:
   "404":
      $ref: "#/components/responses/NotFound"
   "500":
      $ref: '#/components/responses/InternalError'
```

For creating a new response's ref, use `addResponse` function in `cmd/spec/main.go`.

## Generating specification

1. For updating the spec files run
   ```sh
   make generate-spec
   ```
   This updates the `openapi.gen.json` and `openapi.gen.yml` files under `/api` folder
2. Run the server and navigate to `/docs` to see the changes visually
3. Run `make validate-spec` to validate the spec
