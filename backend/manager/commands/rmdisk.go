package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"errors"
	"fmt"
	"os"
	"strings"
)

// rmdisk -driveletter=<letter>
func Fn_rmdisk(params string) (string, error) {
	paramDefs := map[string]Structs.ParamDef{
		"-driveletter": {Required: true},
	}

	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		return "", err
	}

	driveLetter := strings.ToUpper(parsed["-driveletter"])

	// Validación de letra
	if len(driveLetter) != 1 || driveLetter[0] < 'A' || driveLetter[0] > 'Z' {
		return "", errors.New("la letra de unidad debe ser una sola letra entre A y Z")
	}

	return rmdisk(driveLetter)
}

func rmdisk(driveLetter string) (string, error) {
	fileName := fmt.Sprintf("%s%s.dsk", globals.PathDisks, driveLetter)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return "", fmt.Errorf("el disco [%s] no existe", driveLetter)
	}

	if err := os.Remove(fileName); err != nil {
		return "", fmt.Errorf("no se pudo eliminar el disco [%s]: %v", driveLetter, err)
	}

	return fmt.Sprintf("Disco [%s] eliminado exitosamente", driveLetter), nil
}
