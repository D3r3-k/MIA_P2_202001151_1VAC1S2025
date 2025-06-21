package main

import (
	router "MIA_PI_202001151_1VAC1S2025/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	// Configuración del puerto y el logger
	PORT := "3001"
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// Inicialización del router
	r := router.SetupRouter()
	// Configuración del logger para mostrar el encabezado
	fmt.Print(encabezado(PORT))
	// Configuración de CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
	// Iniciar el servidor HTTP con CORS
	log.Fatal(http.ListenAndServe(":"+PORT, c.Handler(r)))
}

func encabezado(PORT string) string {
	header := ""
	header += fmt.Sprintf("%-21s╔═══════════════════════════════╗\n", "")
	header += "╔════════════════════╣ MIA P1F2 202001151 1VAC1S2025 ╠══════════════════╗\n"
	header += "║                    ╚═══════════════════════════════╝                  ║\n"
	header += "║ Escuchando en el puerto " + PORT + "...                                       ║\n"
	header += "╚═══════════════════════════════════════════════════════════════════════╝\n"
	return header

}
