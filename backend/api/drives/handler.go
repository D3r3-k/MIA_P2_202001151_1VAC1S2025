package drivesInfo

import (
	utilsApi "MIA_PI_202001151_1VAC1S2025/utils"
	"encoding/json"
	"net/http"
)

type DrivesInfoResponse struct {
	TotalDisks      int    `json:"totalDisks"`
	TotalPartitions int    `json:"totalPartitions"`
	TotalSize       string `json:"totalSize"`
}

func DrivesInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	totalDisks := 0
	totalPartitions := 0
	totalSize := "N/A"
	totalDisks = utilsApi.CountDisks()
	totalPartitions, totalSize = utilsApi.CalculateTotalPartitions()

	response := DrivesInfoResponse{
		TotalDisks:      totalDisks,
		TotalPartitions: totalPartitions,
		TotalSize:       totalSize,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DrivesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := utilsApi.GetDiskInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// retornar la data como array de objetos JSON
	if len(data) == 0 {
		http.Error(w, "No disks found", http.StatusNotFound)
		return
	}
	// Enviar la respuesta como JSON
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
