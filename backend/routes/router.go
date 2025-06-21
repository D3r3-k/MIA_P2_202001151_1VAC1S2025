package router

import (
	drivesInfo "MIA_PI_202001151_1VAC1S2025/api/drives"
	"MIA_PI_202001151_1VAC1S2025/api/execute"
	"MIA_PI_202001151_1VAC1S2025/api/login"
	"MIA_PI_202001151_1VAC1S2025/api/status"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/status", status.StatusHandler).Methods("GET")
	api.HandleFunc("/execute", execute.ExecuteHandler).Methods("POST")
	api.HandleFunc("/drives", drivesInfo.DrivesHandler).Methods("GET")
	api.HandleFunc("/drives/info", drivesInfo.DrivesInfoHandler).Methods("GET")
	api.HandleFunc("/login", login.LoginHandler).Methods("POST")

	return r
}
