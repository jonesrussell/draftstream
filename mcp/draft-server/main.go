// draft-server/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Request struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	ID     int             `json:"id"`
}

type Response struct {
	Result interface{} `json:"result"`
	ID     int         `json:"id"`
}

func main() {
	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method == "generateMarkdown" {
			var input struct {
				Title string `json:"title"`
				Notes string `json:"notes"`
			}
			json.Unmarshal(req.Params, &input)

			// Mock markdown generation (replace with AI call)
			md := "# " + input.Title + "\n\n" + input.Notes

			resp := Response{Result: md, ID: req.ID}
			json.NewEncoder(w).Encode(resp)
		}
	})

	log.Println("Draft server listening on :8080...")
	http.ListenAndServe(":8080", nil)
}
