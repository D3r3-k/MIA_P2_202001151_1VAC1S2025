package find

import (
	"MIA_PI_202001151_1VAC1S2025/manager/commands"
	utilsApi "MIA_PI_202001151_1VAC1S2025/utils"
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    "Invalid request payload",
			Response: nil,
			Status:   "error",
		})
		return
	}

	output := commands.Fn_Find("-path=" + req.Path + " -name=*")
	if output.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    "Error processing find command: " + output.Error.Error(),
			Response: nil,
			Status:   "error",
		})
		return
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

	response := FindResponse{
		Root: convertFindResponse(output.Object.Children),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utilsApi.StandardResponse{
		Error:    "",
		Response: response,
		Status:   "success",
	})
}
