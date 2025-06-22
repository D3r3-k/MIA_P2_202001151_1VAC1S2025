package cat

import (
	"MIA_PI_202001151_1VAC1S2025/manager/commands"
	"encoding/json"
	"net/http"
)

type CatRequest struct {
	Path string `json:"path"`
}

type CatResponse struct {
	ID          string        `json:"ID"`
	Name        string        `json:"Name"`
	Type        string        `json:"Type"`
	Path        string        `json:"Path"`
	Children    []CatResponse `json:"Children,omitempty"`
	Size        string        `json:"Size"`
	CreatedAt   string        `json:"CreatedAt"`
	Owner       string        `json:"Owner"`
	Content     string        `json:"Content"`
	Extension   string        `json:"Extension"`
	Permissions string        `json:"Permissions"`
}

func CatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	cat, err := commands.Fn_Cat("-file1=" + req.Path)
	if err != nil {
		http.Error(w, "Error processing cat command: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if cat == nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cat)
}
