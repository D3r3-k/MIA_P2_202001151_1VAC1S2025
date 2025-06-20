package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/global"
	Structs "MIA_PI_202001151_1VAC1S2025/structs"
	"MIA_PI_202001151_1VAC1S2025/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// rmdisk -driveletter=<letter>
func Fn_rmdisk(params string) {
	paramDefs := map[string]Structs.ParamDef{
		"-driveletter": {Required: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		fmt.Println(err)
		return
	}
	driveLetter := strings.ToUpper(parsed["-driveletter"])
	utils.ShowMessage("Desea eliminar el disco "+driveLetter+"?\nS: para reemplazar\nN: para cancelar", true)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">>> ")
	scanner.Scan()
	response := strings.TrimSpace(scanner.Text())
	if strings.ToUpper(response) != "S" {
		utils.ShowMessage("Operación cancelada por el usuario.", false)
		return
	}
	Rmdisk(driveLetter)
}

func Rmdisk(driveLetter string) {
	driveLetter = strings.ToUpper(driveLetter)
	if len(driveLetter) != 1 || driveLetter[0] < 'A' || driveLetter[0] > 'Z' {
		utils.ShowMessage("La letra de unidad debe ser una sola letra entre A y Z.", true)
		return
	}
	fileName := fmt.Sprintf(globals.PathDisks+"%s.dsk", driveLetter)
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		utils.ShowMessage(fmt.Sprintf("El disco [%s] no existe.", driveLetter), true)
		return
	}
	err = os.Remove(fileName)
	if err != nil {
		utils.ShowMessage(fmt.Sprintf("No se pudo eliminar el disco [%s]: %v", driveLetter, err), true)
		return
	}
	utils.ShowMessage(fmt.Sprintf("Disco [%s] eliminado exitosamente.", driveLetter), false)
}
