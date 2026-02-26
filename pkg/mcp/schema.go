package mcp

import (
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
)

var (
	getFeedSchema = &jsonschema.Schema{
		Type: "object",
	}

	getEntrySchema = &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"entryId": {
				Type:        "integer",
				Description: "id of the entry",
			},
		},
		Required: []string{"entryId"},
	}

	getEntriesSchema = &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"status": {
				Type:        "string",
				Enum:        []any{"unread", "read", "removed"},
				Description: "status of the entry ",
			},
			"search": {
				Type:        "string",
				Description: "search term query",
			},
			"limit": {
				Type:        "integer",
				Description: "number of entries",
				Default:     json.RawMessage("100"),
			},
			"category": {
				Type:        "string",
				Description: "category of the entries",
			},
			"feed": {
				Type:        "string",
				Description: "feed of the entries",
			},
			"publishedAfter": {
				Type:        "string",
				Format:      "date-time",
				Description: "filter entries published after this date in ISO 8601 format",
			},
			"publishedBefore": {
				Type:        "string",
				Format:      "date-time",
				Description: "filter entries published before this date in ISO 8601 format",
			},
		},
	}

	entrySchema = &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"id": {
				Type:        "integer",
				Description: "id of the entry",
			},
			"title": {
				Type:        "string",
				Description: "title of the entry",
			},
			"url": {
				Type:        "string",
				Description: "url of the entry",
			},
			"content": {
				Type:        "string",
				Description: "content of the entry",
			},
			"status": {
				Type:        "string",
				Description: "statues of the entry",
				Enum:        []any{"unread", "read", "removed"},
			},
			"createdAt": {
				Type:        "string",
				Format:      "date-time",
				Description: "time the entry created at",
			},
		},
	}

	entriesSchema = &jsonschema.Schema{
		Type:  "array",
		Items: entrySchema,
	}

	feedSchema = &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"id": {
				Type:        "integer",
				Description: "id of the entry",
			},
			"title": {
				Type:        "string",
				Description: "title of the entry",
			},
			"url": {
				Type:        "string",
				Description: "url of the entry",
			},
		},
	}

	feedsSchema = &jsonschema.Schema{
		Type:  "array",
		Items: feedSchema,
	}
)
