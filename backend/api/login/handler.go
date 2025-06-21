package login

import (
	"MIA_PI_202001151_1VAC1S2025/manager/commands"
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	PartitionID string `json:"partition_id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserData struct {
		Username     string `json:"username"`
		Group        string `json:"group"`
		Partition_id string `json:"partition_id"`
	} `json:"user_data"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if loginReq.PartitionID == "" || loginReq.Username == "" || loginReq.Password == "" {
		http.Error(w, "Missing partition_id, username, or password", http.StatusBadRequest)
		return
	}

	res := commands.Fn_Login("-user=" + loginReq.Username + " -pass=" + loginReq.Password + " -id=" + loginReq.PartitionID)

	if !res {
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}
	nombre, grupo := utils.GetUserAndGroupNames(loginReq.PartitionID, globals.LoginSession.UID, globals.LoginSession.GID)
	var loginRes LoginResponse
	loginRes.Token = "true"
	loginRes.UserData.Username = nombre
	loginRes.UserData.Group = grupo
	loginRes.UserData.Partition_id = loginReq.PartitionID
	json.NewEncoder(w).Encode(loginRes)
}

// [Logout]
type LogoutResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commands.Fn_Logout("")

	logoutRes := LogoutResponse{Status: true, Message: "Sesión cerrada correctamente."}
	json.NewEncoder(w).Encode(logoutRes)
}
