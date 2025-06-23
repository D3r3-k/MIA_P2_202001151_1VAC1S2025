package execute

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	"MIA_PI_202001151_1VAC1S2025/manager/lib"
)

type ExecuteRequest struct {
	Commands string `json:"commands"`
}

type ExecuteResponse struct {
	Error    string `json:"error,omitempty"`
	Response string `json:"response"`
	Status   string `json:"status"`
}

func ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	globals.Output = ""

	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ExecuteResponse{
			Error:    "JSON inválido",
			Response: "",
			Status:   "error",
		})
		return
	}

	if strings.TrimSpace(req.Commands) == "" {
		log.Println("Error: Comando no proporcionado")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ExecuteResponse{
			Error:    "Comando no proporcionado",
			Response: "",
			Status:   "error",
		})
		return
	}

	lines := strings.Split(req.Commands, "\n")
	status := "success"
	var firstError string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		comando, parametros := lib.GetCommands(line)
		if !lib.AnalyzeExistCommand(comando) {
			msg := fmt.Sprintf("Comando no reconocido: %s", comando)
			log.Printf("[Error] %s", msg)
			if firstError == "" {
				firstError = "Error en el comando: " + comando
			}
			status = "error"
			continue
		}

		log.Printf("Ejecutando línea de comando: %s", line)
		output, err := lib.AnalyzeCommand(comando, parametros)
		if err != nil {
			msg := err.Error()
			log.Printf("Error al ejecutar el comando '%s': %v", comando, err)
			globals.Output += "[Error] " + msg + "\n"
			if firstError == "" {
				firstError = "Error en el comando: " + comando
			}
			status = "error"
			continue
		}

		globals.Output += output + "\n"
		log.Printf("Resultado de la ejecución: %s", output)
	}

	if status == "success" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(ExecuteResponse{
		Error:    firstError,
		Response: globals.Output,
		Status:   status,
	})
}
