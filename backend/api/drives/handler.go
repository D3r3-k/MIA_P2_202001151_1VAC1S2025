package drivesInfo

import (
	utilsApi "MIA_PI_202001151_1VAC1S2025/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	if len(data) == 0 {
		http.Error(w, "No disks found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type DriveInfoResponse struct {
	Name       string `json:"Name"`
	Path       string `json:"Path"`
	Size       string `json:"Size"`
	Fit        string `json:"Fit"`
	Partitions int    `json:"Partitions"`
}

func DriveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	driveletter := mux.Vars(r)["driveletter"]
	if driveletter == "" {
		http.Error(w, "Missing drive letter", http.StatusBadRequest)
		return
	}
	data, err := utilsApi.GetDriveInfo(driveletter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if data == (utilsApi.DriveInfo{}) {
		http.Error(w, "Drive not found", http.StatusNotFound)
		return
	}
	_size := utilsApi.ConvertSizeToString(data.Size)
	response := DriveInfoResponse{
		Name:       data.Name,
		Path:       data.Path,
		Size:       _size,
		Fit:        data.Fit,
		Partitions: data.Partitions,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

type DrivePartitionsResponse struct {
	Partitions []struct {
		Status     string `json:"status"`
		Type       string `json:"type"`
		Fit        string `json:"fit"`
		Size       int32  `json:"size"`
		Start      int32  `json:"start"`
		Name       string `json:"name"`
		ID         string `json:"id"`
		Path       string `json:"path"`
		Date       string `json:"date"`
		Filesystem string `json:"filesystem"`
	} `json:"partitions"`
}

func DrivePartitionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	driveletter := mux.Vars(r)["driveletter"]
	if driveletter == "" {
		http.Error(w, "Missing drive letter", http.StatusBadRequest)
		return
	}
	data, err := utilsApi.GetDrivePartitions(driveletter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if data == nil {
		http.Error(w, "Drive not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
