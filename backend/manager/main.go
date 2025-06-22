package main

import (
	"MIA_PI_202001151_1VAC1S2025/manager/cmd"
	"MIA_PI_202001151_1VAC1S2025/manager/lib"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Printf("%-21s╔═════════════════════════════╗\n", "")
	fmt.Println("╔════════════════════╣ MIA PI 202001151 1VAC1S2025 ╠════════════════════╗")
	fmt.Println("║                    ╚═════════════════════════════╝                    ║")
	fmt.Println("║ [help] para ver los comandos disponibles                              ║")
	fmt.Println("║ [cls] para limpiar la pantalla                                        ║")
	fmt.Println("║ [exit] para salir                                                     ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════════════╝")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">>> ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		comando, parametros := lib.GetCommands(input)
		comando = strings.ToLower(comando)
		comando = strings.TrimSpace(comando)
		switch comando {
		case "exit":
			utils.ShowMessage("Saliendo del programa...", false)
			return
		case "execute":
			cmd.Execute(parametros)
			continue
		case "cls":
			cmd.Cls()
			continue
		case "help":
			cmd.Help()
			continue
		default:
			lib.AnalyzeCommand(comando, parametros)
		}
	}
}
