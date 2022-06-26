# API Documentation

This service uses openapi v.3 for documenting the service's API

## Docs

We use [redoc](https://github.com/Redocly/redoc) for auto generating swagger UI based on the spec file
The docs locates under `<root>/docs` and the json spec under `<root>/spec.json`
In addition, you can get the `openapi.json` file under `/api/provisioning/openapi.json`

## Adding new endpoint

We chose an hybrid approach for adding a new endpoint

### How to add new type/schema
1. Create the endpoint's payload under `internal/payload`
2. Register it in `cmd/spec/main.go` script:
```go
gen.addSchema("v1.<YOUR_PAYLOAD>", &payloads.<YOUR_PAYLOAD>{})
```
3. use the `addSchema` for registering new errors payloads if needed

### How to add new new path
Edit `/cmd/spec/paths.yaml` for adding the new route manually
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

### Running all together

1. For updating the spec files run
   ```sh
   make generate-spec
   ```
   This updates the `openapi.gen.json` and `openapi.gen.yml` files under `/api` folder
2. Run the server and navigate to `/docs` to see the changes visually