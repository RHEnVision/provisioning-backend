package main

type Parameter struct {
	Name        string      `json:"name" yaml:"name"`
	Description string      `json:"description" yaml:"description"`
	Required    bool        `json:"required" yaml:"required"`
	Default     interface{} `json:"default" yaml:"default"`
	Type        string      `json:"type" yaml:"type"`
	In          string      `json:"in" yaml:"in"`
}

var LimitQueryParam = Parameter{
	Name:        "limit",
	Description: "The number of items to return",
	Default:     100,
	Type:        "integer",
	Required:    false,
	In:          "query",
}

var OffsetQueryParam = Parameter{
	Name:        "offset",
	Description: "The number of items to skip before starting to collect the result set",
	Default:     0,
	Type:        "integer",
	Required:    false,
	In:          "query",
}

var TokenQueryParam = Parameter{
	Name:        "token",
	Description: "The token used for requesting the next page of results; empty token for the first page",
	Default:     "",
	Type:        "string",
	Required:    false,
	In:          "query",
}
