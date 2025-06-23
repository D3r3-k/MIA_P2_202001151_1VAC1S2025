package drivesInfo

import (
	utilsApi "MIA_PI_202001151_1VAC1S2025/utils"
	"encoding/json"
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

	totalDisks := utilsApi.CountDisks()
	totalPartitions, totalSize := utilsApi.CalculateTotalPartitions()

	response := DrivesInfoResponse{
		TotalDisks:      totalDisks,
		TotalPartitions: totalPartitions,
		TotalSize:       totalSize,
	}

	json.NewEncoder(w).Encode(utilsApi.StandardResponse{
		Response: response,
		Status:   "success",
	})
}

func DrivesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := utilsApi.GetDiskInfo()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    "Ha ocurrido un error al obtener la informacion de los discos",
			Status:   "error",
			Response: []interface{}{},
		})
		return
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    "No hay discos disponibles",
			Status:   "error",
			Response: []interface{}{},
		})
		return
	}

	json.NewEncoder(w).Encode(utilsApi.StandardResponse{
		Response: data,
		Status:   "success",
	})
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
	response := DriveInfoResponse{
		Name:       "",
		Path:       "",
		Size:       "",
		Fit:        "",
		Partitions: 0,
	}

	driveletter := mux.Vars(r)["driveletter"]
	if driveletter == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    "Missing drive letter",
			Status:   "error",
			Response: response,
		})
		return
	}

	data, err := utilsApi.GetDriveInfo(driveletter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    err.Error(),
			Status:   "error",
			Response: response,
		})
		return
	}

	if data == (utilsApi.DriveInfo{}) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    "Drive not found",
			Status:   "error",
			Response: response,
		})
		return
	}

	_size := utilsApi.ConvertSizeToString(data.Size)
	response = DriveInfoResponse{
		Name:       data.Name,
		Path:       data.Path,
		Size:       _size,
		Fit:        data.Fit,
		Partitions: data.Partitions,
	}

	json.NewEncoder(w).Encode(utilsApi.StandardResponse{
		Response: response,
		Status:   "success",
	})
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    "Missing drive letter",
			Status:   "error",
			Response: []interface{}{},
		})
		return
	}

	data, err := utilsApi.GetDrivePartitions(driveletter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    err.Error(),
			Status:   "error",
			Response: []interface{}{},
		})
		return
	}

	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(utilsApi.StandardResponse{
			Error:    "Drive not found",
			Status:   "error",
			Response: []interface{}{},
		})
		return
	}

	json.NewEncoder(w).Encode(utilsApi.StandardResponse{
		Response: data,
		Status:   "success",
	})
}
