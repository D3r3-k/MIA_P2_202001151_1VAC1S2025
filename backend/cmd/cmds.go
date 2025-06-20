package cmd

import (
	globals "MIA_PI_202001151_1VAC1S2025/global"
	"MIA_PI_202001151_1VAC1S2025/lib"
	Structs "MIA_PI_202001151_1VAC1S2025/structs"
	"MIA_PI_202001151_1VAC1S2025/utils"
	"bufio"
	"fmt"
	"strings"
)

func Cls() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("|====================| MIA PI 202001151 1VAC1S2025 |====================|")
	fmt.Println("| [exit] para salir                                                     |")
	fmt.Println("| [help] para ver los comandos disponibles                              |")
	fmt.Println("| [cls] para limpiar la pantalla                                        |")
	fmt.Println("|=======================================================================|")
}

func Help() {
	fmt.Println("Comandos disponibles:")
	for cmd, desc := range globals.Commands {
		fmt.Printf("[>] %s:\t %s\n", cmd, desc)
	}
}

func Execute(params string) {
	paramDefs := map[string]Structs.ParamDef{
		"-path": {Required: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return
	}
	path := parsed["-path"]

	if !strings.HasSuffix(strings.ToLower(path), ".sdaa") {
		utils.ShowMessage("Solo se permiten archivos con extensión .sdaa", true)
		return
	}

	file, err := utils.OpenFile(path)
	if err != nil {
		utils.ShowMessage("Error al abrir el archivo: "+err.Error(), true)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' || len(line) == 0 {
			continue
		}
		fmt.Println(">> ", line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(line)
		input := strings.TrimSpace(line)
		if strings.ToLower(input) == "exit" {
			break
		}
		comando, parametros := lib.GetCommands(input)
		lib.AnalyzeCommand(comando, parametros)
	}
}
