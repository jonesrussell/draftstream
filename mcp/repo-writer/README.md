# Repo Writer Service

A JSON-RPC 2.0 service for writing Jekyll blog drafts with proper front matter.

## Features

- Accepts JSON-RPC 2.0 requests with method `writeJekyllDraft`
- Automatically creates `_drafts` directory if it doesn't exist
- Generates proper Jekyll front matter with all metadata
- Sanitizes filenames for safe file system usage
- Returns success/error responses with proper JSON-RPC format

## Usage

### Start the server

```bash
cd mcp/repo-writer
go run main.go
```

The server will listen on port 8081.

### Request Format

Send a POST request to `http://localhost:8081/mcp` with the following JSON structure:

```json
{
  "jsonrpc": "2.0",
  "method": "writeJekyllDraft",
  "params": {
    "title": "Your Post Title",
    "tags": ["tag1", "tag2"],
    "categories": ["category1", "category2"],
    "series": "series-name",
    "summary": "One-line summary",
    "body": "Full markdown content",
    "path": "/path/to/your/blog"
  },
  "id": 1
}
```

### Response Format

**Success:**
```json
{
  "result": "written",
  "id": 1
}
```

**Error:**
```json
{
  "error": {
    "code": -32602,
    "message": "Invalid params"
  },
  "id": 1
}
```

### Generated File

The service will create a file at `${path}/_drafts/${sanitized-title}.md` with the following structure:

```yaml
---
layout: post
title: "Your Post Title"
date: 2024-01-15
categories: [category1, category2]
tags: [tag1, tag2]
series: series-name
summary: "One-line summary"
---

# Your markdown content here
```

## Testing

You can test the service using the provided `test_request.json` file:

```bash
curl -X POST http://localhost:8081/mcp \
  -H "Content-Type: application/json" \
  -d @test_request.json
```

## Error Codes

- `-32700`: Parse error
- `-32601`: Method not found
- `-32602`: Invalid params
- `-32603`: Internal error (file system operations) 