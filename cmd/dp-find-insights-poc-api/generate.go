package main

// This file is only used to hold go generate commands

// APIs are generated for public and private endpoints.
// Endpoints tagged with doc are ignored since they are only in the swagger spec for documentation.
// They are actually handled outside the generated APIs.

//go:generate oapi-codegen -generate types,chi-server -include-tags public,private,spec -exclude-tags doc -package api -o ../../api/api.go ../../swagger.yaml
