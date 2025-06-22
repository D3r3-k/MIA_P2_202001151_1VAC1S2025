package find

import (
	"MIA_PI_202001151_1VAC1S2025/manager/commands"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type FindRequest struct {
	Path string `json:"path"`
}

type FindResponse struct {
	Root []FileSystemItem `json:"Root"`
}

type FileSystemItem struct {
	ID       string           `json:"ID"`
	Name     string           `json:"Name"`
	Path     string           `json:"Path"`
	Type     string           `json:"Type"`
	Children []FileSystemItem `json:"Children,omitempty"`
}

func FindHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req FindRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ejecutar el comando find
	finds, err := commands.Fn_Find("-path=" + req.Path + " -name=*")
	if err != nil {
		http.Error(w, "Error processing find command: "+err.Error(), http.StatusInternalServerError)
		return
	}

	root := req.Path
	if root == "" {
		root = "/"
	}
	var convertFindResponse func([]commands.FindResponse) []FileSystemItem
	convertFindResponse = func(items []commands.FindResponse) []FileSystemItem {
		var result []FileSystemItem
		for i, item := range items {
			result = append(result, FileSystemItem{
				ID:       fmt.Sprintf("%d", i+1),
				Path:     item.Path,
				Name:     item.Name,
				Type:     strings.ToLower(item.Type),
				Children: convertFindResponse(item.Children),
			})
		}
		return result
	}

	formatted := FindResponse{
		Root: convertFindResponse(finds[0].Children),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(formatted)
}
