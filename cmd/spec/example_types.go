package main

import (
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
)

var InstanceTypesAWSResponse = []payloads.InstanceTypeResponse{{
	Name:               "c5a.8xlarge",
	VCPUs:              32,
	Cores:              16,
	MemoryMiB:          65536,
	EphemeralStorageGB: 0,
	Supported:          true,
	Architecture:       "x86_64",
	AzureDetail:        nil,
}}

var InstanceTypesAzureResponse = []payloads.InstanceTypeResponse{{
	Name:               "Standard_M128s",
	VCPUs:              128,
	Cores:              64,
	MemoryMiB:          2000000,
	EphemeralStorageGB: 4096,
	Supported:          true,
	Architecture:       "x86_64",
	AzureDetail: &clients.InstanceTypeDetailAzure{
		GenV1: true,
		GenV2: true,
	},
}}

var InstanceTypesGCPResponse = []payloads.InstanceTypeResponse{{
	Name:               "e2-highcpu-16",
	VCPUs:              16,
	MemoryMiB:          15623,
	EphemeralStorageGB: 0,
	Supported:          true,
	Architecture:       "x86_64",
	AzureDetail:        nil,
}}
