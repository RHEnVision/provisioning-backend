package main

import "github.com/RHEnVision/provisioning-backend/internal/payloads"

var ResponseErrorGenericExample = payloads.ResponseError{
	TraceId:   "b57f7b78c",
	Error:     "error: this can be pretty long string",
	Version:   "df8a489",
	BuildTime: "2023-04-14_17:15:02",
}

var ResponseNotFoundErrorExample = payloads.ResponseError{
	TraceId:   "b57f7b78c",
	Error:     "error: resource not found: details can be long",
	Version:   "df8a489",
	BuildTime: "2023-04-14_17:15:02",
}

var ResponseBadRequestErrorExample = payloads.ResponseError{
	TraceId:   "b57f7b78c",
	Error:     "error: bad request: details can be long",
	Version:   "df8a489",
	BuildTime: "2023-04-14_17:15:02",
}

var ResponseErrorUserFriendlyExample = payloads.ResponseError{
	Message:   "vCPU limit reached, contact AWS support",
	TraceId:   "b57f7b78c",
	Error:     "cannot run instances: cannot run instances: operation error EC2: RunInstances, https response error StatusCode: 400, RequestID: af6e10a0-75c9-47e1-a83c-872494250322, VcpuLimitExceeded: You have requested more vCPU capacity than your current vCPU limit of 128 allows for the instance bucket that the specified instance type belongs to. Please visit http://aws.amazon.com/contact-us/ec2-request to request an adjustment to this limit.",
	Version:   "df8a489",
	BuildTime: "2023-04-14_17:15:02",
}
