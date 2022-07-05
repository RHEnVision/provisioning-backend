# API Documentation

This service uses openapi v.3 for documenting the service's API

## Docs

We use [redoc](https://github.com/Redocly/redoc) for auto generating swagger UI based on the spec file
The docs locate under `<root>/docs` and the json spec under `<root>/openapi.json`
In addition, `openapi.json` serves under `/api/provisioning/v1/openapi.json`

## Adding new endpoint

We chose an hybrid approach for adding a new endpoint

### How to add new type/schema
1. Create the endpoint's payload under `internal/payload`
2. Register it in `cmd/spec/main.go` script:
```go
gen.addSchema("v1.<YOUR_PAYLOAD>", &payloads.<YOUR_PAYLOAD>{})
```
3. use the `addSchema` for registering new errors payloads if needed

### How to add a new path
Edit `/cmd/spec/path.yaml` for adding the new route manually
It is recommended to use a dedicated openapi editor plugin for your IDE for fast editing

```yml
paths:
  /<NEW_ENDPOINT>:
    get:
      operationId: getListOfNewStuff
      description: ''
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/v1.<YOUR_PAYLOAD>' # a reference to the registered type 
```

### How to add error responses to endpoints
You can reuse and reference predefined error responses:

 ```yaml
 # path.yaml
 responses:
  #"200": "See above example"
   "404":
      $ref: "#/components/responses/NotFound"
   "500":
      $ref: '#/components/responses/InternalError'
 ```
For creating a new response's ref, use `addResponse` function in `cmd/spec/main.go`:
```go
// addResponse(name, description, schema ref)
gen.addResponse("NotAuthorized", "The user is not authorized", "#/components/schemas/v1.ResponseError")
```
Then, consume it in your paths:
```yaml
 # path.yaml
 responses:
  #"200": "See above example"
   "403":
      $ref: "#/components/responses/NotAuthorized"
```

### Running all together

1. For updating the spec files run
   ```sh
   make generate-spec
   ```
   This updates the `openapi.gen.json` and `openapi.gen.yml` files under `/api` folder
2. Run the server and navigate to `/docs` to see the changes visually
3. Run `make validate-spec` to validate the spec
