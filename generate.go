package main

// This file is only used to hold go generate commands

// APIs are generated for public and private endpoints.
// Endpoints tagged with doc are ignored since they are only in the swagger spec for documentation.
// They are actually handled outside the generated APIs.

//go:generate oapi-codegen -generate types -package api -o api/types.go swagger.yaml
//go:generate oapi-codegen -generate chi-server -include-tags public -exclude-tags doc -package public -o api/public/api.go swagger.yaml
//go:generate oapi-codegen -generate chi-server -include-tags public,private -exclude-tags doc -package private -o api/private/api.go swagger.yaml
