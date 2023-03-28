# 14. Union types in OpenSPEC

Authors: Lukáš Zapletal


## Status

_Proposed_


## Problem Statement

OpenAPI specification contains `oneOf`/`anyOf` keyword that can be used used to implement union type. The Go language, however, does not support union type at all. It is difficult to work with union types through the OpenAPI implementation (kin-openapi) we currently use because there is no native support for this for generating components.

* https://en.wikipedia.org/wiki/Union_type
* https://swagger.io/docs/specification/data-models/oneof-anyof-allof-not/
* https://github.com/getkin/kin-openapi

## Goals

In our project, we need to be able to return REST payloads of different types (e.g. AWS, Azure, GCP) in a reliable and convinient way, ideally keeping the generating OpenAPI components from our Go code.

## Non-goals

The intention is only to set an expectation and standard for union types (oneOf/anyOf), not all OpenAPI features.

## Current Architecture

Image builder currently have at least one payload that makes use of anyOf. The problem is that image builder does not use generating componts from the code, so creating the scheme was possible. Writing and maintaining handler in Go is complicated, because this OpenAPI feature is implemented via raw JSON

* https://github.com/osbuild/image-builder/blob/00a28a72ffb7936e302974e75831a4373b6dd1e7/internal/v1/api.yaml#L591-L604
* https://github.com/osbuild/image-builder/blob/a1677b19f43e5265767a08a74d830abcdcddc7af/internal/v1/handler.go#L581-L712

## Proposed Architecture

Instead of giving up on our hybrid OpenAPI generation approach and starting with manual edits using anyOf:

```yaml
UnionResponse:
  type: object
  required:
    - type
    - data
  properties:
    type:
      $ref: '#/components/schemas/UnionTypes'
    data:
      anyOf:
        - $ref: '#/components/schemas/AWSData'
        - $ref: '#/components/schemas/GCPData'
        - $ref: '#/components/schemas/AzureData'
UnionTypes:
  type: string
  enum: ['aws', 'gcp', 'azure']
```

which would lead to handling code with raw JSON unmarshaling, the proposed solution is to use the following approach without anyOf feature:

```yaml
UnionResponse:
  type: object
  required:
    - type
  properties:
    type:
      $ref: '#/components/schemas/UnionTypes'
    aws:
      $ref: '#/components/schemas/AWSData'
    gcp:
      $ref: '#/components/schemas/GCPData'
    azure:
      $ref: '#/components/schemas/AzureData'
UnionTypes:
  type: string
  enum: ['aws', 'gcp', 'azure']
AWSData:
  type: object
  # ...
GCPData:
  type: object
  # ...
AzureData:
  type: object
  # ...
```

Examples how payloads would look like. For AWS:

```json
{
    "type": "aws",
    "aws": {...}
}
```

For GCP:

```json
{
    "type": "gcp",
    "gcp": {...}
}
```

The handling code can unmarshal the whole struct once and processing will be straightforward. This will work with any language or stack since it is not using any specific features.

## Challenges

N/A

## Alternatives Considered

N/A

## Dependencies

N/A
