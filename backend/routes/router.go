package router

import (
	"MIA_PI_202001151_1VAC1S2025/api/cat"
	drivesInfo "MIA_PI_202001151_1VAC1S2025/api/drives"
	"MIA_PI_202001151_1VAC1S2025/api/execute"
	"MIA_PI_202001151_1VAC1S2025/api/find"
	"MIA_PI_202001151_1VAC1S2025/api/login"
	"MIA_PI_202001151_1VAC1S2025/api/status"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	// [Status] verificar el estado del servidor y la sesion del usuario
	api.HandleFunc("/status", status.StatusHandler).Methods("GET")
	// [Execute] ejecutar comandos en el servidor
	api.HandleFunc("/execute", execute.ExecuteHandler).Methods("POST")
	// [Drives] obtener informacion de los discos y particiones
	api.HandleFunc("/drives", drivesInfo.DrivesHandler).Methods("GET")
	api.HandleFunc("/drives/info", drivesInfo.DrivesInfoHandler).Methods("GET")
	api.HandleFunc("/drives/{driveletter}", drivesInfo.DriveHandler).Methods("GET")
	api.HandleFunc("/drives/{driveletter}/partitions", drivesInfo.DrivePartitionsHandler).Methods("GET")
	// [Login/Logout] manejar el inicio y cierre de sesion
	api.HandleFunc("/login", login.LoginHandler).Methods("POST")
	api.HandleFunc("/logout", login.LogoutHandler).Methods("POST")
	// [Find] buscar archivos en el servidor
	api.HandleFunc("/find", find.FindHandler).Methods("POST")
	// [Cat] mostrar el contenido de un archivo
	api.HandleFunc("/cat", cat.CatHandler).Methods("POST")
	return r
}
