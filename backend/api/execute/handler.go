package execute

import (
	"MIA_PI_202001151_1VAC1S2025/manager/lib"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type ExecuteRequest struct {
	Command string `json:"command"`
}

func ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ExecuteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "JSON inválido"})
		return
	}
	if req.Command == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Comando no proporcionado"})
		return
	}
	lines := strings.Split(req.Command, "\n")
	for _, line := range lines {
		log.Printf("Ejecutando línea de comando: %s", line)
		if line == "" || line[0] == '#' || len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(line)
		comando, parametros := lib.GetCommands(line)
		lib.AnalyzeCommand(comando, parametros)
	}

	log.Printf("Ejecutando comando: %s", req.Command)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": "Comando recibido: " + req.Command})
}
