package status

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"encoding/json"
	"net/http"
)

type StatusResponse struct {
	Status    string `json:"status"`
	AuthToken struct {
		Username     string `json:"username"`
		Group        string `json:"group"`
		Partition_id string `json:"partition_id"`
	} `json:"authToken"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data StatusResponse
	if globals.LoginSession.User == "" {
		data.Status = "false"
		data.AuthToken = struct {
			Username     string `json:"username"`
			Group        string `json:"group"`
			Partition_id string `json:"partition_id"`
		}{}
	} else {
		nombre, grupo := utils.GetUserAndGroupNames(string(globals.LoginSession.PartitionID[:]), globals.LoginSession.UID, globals.LoginSession.GID)
		data.Status = "true"
		data.AuthToken = struct {
			Username     string `json:"username"`
			Group        string `json:"group"`
			Partition_id string `json:"partition_id"`
		}{
			Username:     nombre,
			Group:        grupo,
			Partition_id: string(globals.LoginSession.PartitionID[:]),
		}
	}
	json.NewEncoder(w).Encode(StatusResponse{Status: data.Status, AuthToken: data.AuthToken})
}
