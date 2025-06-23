package cat

import (
	"MIA_PI_202001151_1VAC1S2025/manager/commands"
	utilsApi "MIA_PI_202001151_1VAC1S2025/utils"
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    "Invalid request payload",
			Response: nil,
			Status:   "error",
		})
		return
	}

	cat, err := commands.Fn_Cat("-file1=" + req.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    "Error processing cat command: " + err.Error(),
			Response: nil,
			Status:   "error",
		})
		return
	}

	if cat == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    "File not found",
			Response: nil,
			Status:   "error",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utilsApi.StandardResponse{
		Error:    "",
		Response: cat,
		Status:   "success",
	})
}
