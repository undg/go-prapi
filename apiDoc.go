package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIDocumentation represents the structure of our API documentation
type APIDocumentation struct {
	Version     string               `json:"version"`
	Description string               `json:"description"`
	Endpoints   []EndpointDoc        `json:"endpoints"`
	Actions     []ActionDoc          `json:"actions"`
	Schemas     map[string]SchemaDoc `json:"schemas"`
}

// EndpointDoc represents documentation for an API endpoint
type EndpointDoc struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

// ActionDoc represents documentation for an API action
type ActionDoc struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SchemaDoc represents documentation for a data schema
type SchemaDoc struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
}

// JSONSchema represents the root of the JSON Schema
type JSONSchema struct {
	Schema      string                 `json:"$schema"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Definitions map[string]interface{} `json:"definitions"`
}

func RenderSchemaJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	schema := GetSchemaJSON()
	jsonData, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>API Documentation</title>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; padding: 20px; }
			pre { background-color: #f4f4f4; padding: 15px; border-radius: 5px; overflow-x: auto; }
		</style>
	</head>
	<body>
		<h1>API Documentation</h1>
		<pre><code>%s</code></pre>
	</body>
	</html>
	`

	fmt.Fprintf(w, html, string(jsonData))
}

// GetSchema returns the JSON Schema for the API
func GetSchemaJSON() interface{} {
	schema := JSONSchema{
		Schema:      "http://json-schema.org/draft-07/schema#",
		Type:        "object",
		Title:       "PulseAudio WebSocket API",
		Description: "API for controlling PulseAudio via WebSocket",
		Definitions: map[string]interface{}{
			"Action": map[string]interface{}{
				"type": "string",
				"enum": []string{
					"GetCards",
					"GetOutputs",
					"GetVolume",
					"GetSchema",
					"GetMute",
					"SetVolume",
					"SetMute",
					"ToggleMute",
				},
			},
			"Message": map[string]interface{}{
				"type": "object",
				"properties": map[string]map[string]interface{}{
					"action": {
						"$ref": "#/definitions/Action",
					},
					"value": {
						"type": "any",
					},
				},
				"required": []string{"action"},
			},
			"Response": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"actionIn": map[string]interface{}{
						"$ref": "#/definitions/Action",
					},
					"status": map[string]interface{}{
						"type": "integer",
					},
					"value": map[string]interface{}{
						"type": "any",
					},
					"error": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"actionIn", "status"},
			},
		},
	}

	return schema
}
