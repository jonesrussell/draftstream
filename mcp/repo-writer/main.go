// repo-writer/main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Request struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	ID     int             `json:"id"`
}

type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  *Error      `json:"error,omitempty"`
	ID     int         `json:"id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type WriteJekyllDraftParams struct {
	Title      string   `json:"title"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
	Series     string   `json:"series"`
	Summary    string   `json:"summary"`
	Body       string   `json:"body"`
	Path       string   `json:"path"`
}

func main() {
	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, -32700, "Parse error", req.ID)
			return
		}

		if req.Method == "writeJekyllDraft" {
			handleWriteJekyllDraft(w, req)
		} else {
			respondWithError(w, -32601, "Method not found", req.ID)
		}
	})

	log.Println("Repo writer server listening on :8081...")
	http.ListenAndServe(":8081", nil)
}

func handleWriteJekyllDraft(w http.ResponseWriter, req Request) {
	var params WriteJekyllDraftParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		respondWithError(w, -32602, "Invalid params", req.ID)
		return
	}

	// Validate required fields
	if params.Title == "" || params.Path == "" {
		respondWithError(w, -32602, "Title and path are required", req.ID)
		return
	}

	// Create _drafts directory if it doesn't exist
	draftsDir := filepath.Join(params.Path, "_drafts")
	if err := os.MkdirAll(draftsDir, 0755); err != nil {
		respondWithError(w, -32603, fmt.Sprintf("Failed to create drafts directory: %v", err), req.ID)
		return
	}

	// Generate filename from title
	filename := sanitizeFilename(params.Title) + ".md"
	filepath := filepath.Join(draftsDir, filename)

	// Generate front matter
	frontMatter := generateFrontMatter(params)

	// Combine front matter and body
	content := frontMatter + "\n" + params.Body

	// Write file
	if err := ioutil.WriteFile(filepath, []byte(content), 0644); err != nil {
		respondWithError(w, -32603, fmt.Sprintf("Failed to write file: %v", err), req.ID)
		return
	}

	// Log confirmation
	log.Printf("Draft written to: %s", filepath)

	// Respond with success
	resp := Response{Result: "written", ID: req.ID}
	json.NewEncoder(w).Encode(resp)
}

func generateFrontMatter(params WriteJekyllDraftParams) string {
	now := time.Now()
	date := now.Format("2006-01-02")

	frontMatter := fmt.Sprintf(`---
layout: post
title: "%s"
date: %s`, params.Title, date)

	if len(params.Categories) > 0 {
		categories := fmt.Sprintf(`["%s"]`, params.Categories[0])
		for i := 1; i < len(params.Categories); i++ {
			categories = fmt.Sprintf(`%s, "%s"`, categories, params.Categories[i])
		}
		frontMatter += fmt.Sprintf("\ncategories: [%s]", categories)
	}

	if len(params.Tags) > 0 {
		tags := fmt.Sprintf(`["%s"]`, params.Tags[0])
		for i := 1; i < len(params.Tags); i++ {
			tags = fmt.Sprintf(`%s, "%s"`, tags, params.Tags[i])
		}
		frontMatter += fmt.Sprintf("\ntags: [%s]", tags)
	}

	if params.Series != "" {
		frontMatter += fmt.Sprintf("\nseries: %s", params.Series)
	}

	if params.Summary != "" {
		frontMatter += fmt.Sprintf("\nsummary: \"%s\"", params.Summary)
	}

	frontMatter += "\n---"
	return frontMatter
}

func sanitizeFilename(title string) string {
	// Simple sanitization - replace spaces with hyphens and remove special chars
	// This is a basic implementation; you might want to use a more robust library
	var result []rune
	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result = append(result, r)
		} else if r == ' ' {
			result = append(result, '-')
		}
	}
	return string(result)
}

func respondWithError(w http.ResponseWriter, code int, message string, id int) {
	resp := Response{
		Error: &Error{Code: code, Message: message},
		ID:    id,
	}
	json.NewEncoder(w).Encode(resp)
}
