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

// Fn_Mkdisk procesa los parámetros y ejecuta la creación del disco.
func Fn_Mkdisk(params string) (string, error) {
	paramDefs := map[string]Structs.ParamDef{
		"-size": {Required: true},
		"-fit":  {Required: false, Default: "ff"},
		"-unit": {Required: false, Default: "m"},
	}

	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		return "", err
	}

	var size int
	if _, err := fmt.Sscanf(parsed["-size"], "%d", &size); err != nil {
		return "", errors.New("el parámetro -size debe ser un número entero")
	}

	fit := strings.ToLower(parsed["-fit"])
	unit := strings.ToLower(parsed["-unit"])

	return mkdisk(size, fit, unit)
}

// mkdisk valida y crea el disco virtual en base a los parámetros recibidos.
func mkdisk(size int, fit, unit string) (string, error) {
	if size <= 0 {
		return "", errors.New("el tamaño debe ser mayor a 0")
	}
	if fit != "ff" && fit != "bf" && fit != "wf" {
		return "", errors.New("el ajuste debe ser <'ff'|'bf'|'wf'>")
	}
	if unit != "k" && unit != "m" {
		return "", errors.New("la unidad debe ser <'k'|'m'>")
	}

	// Calcular tamaño real en bytes
	realSize := size
	if unit == "k" {
		realSize *= 1024
	} else {
		realSize *= 1024 * 1024
	}

	// Buscar una letra de unidad disponible (archivo no existente)
	var fileName string
	var driveLetter rune
	found := false
	for c := 'A'; c <= 'Z'; c++ {
		fileName = fmt.Sprintf("%s%c.dsk", globals.PathDisks, c)
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			driveLetter = c
			found = true
			break
		}
	}
	if !found {
		return "", errors.New("no hay letras disponibles para crear más discos")
	}

	// Crear archivo de disco
	if err := utils.CreateFile(fileName); err != nil {
		return "", fmt.Errorf("error creando archivo: %v", err)
	}

	file, err := utils.OpenFile(fileName)
	if err != nil {
		return "", fmt.Errorf("error abriendo archivo: %v", err)
	}
	defer file.Close()

	zeroBuffer := make([]byte, 1024)
	for i := 0; i < realSize/1024; i++ {
		if err := utils.WriteObject(file, zeroBuffer, int64(i*1024)); err != nil {
			return "", fmt.Errorf("error escribiendo bytes: %v", err)
		}
	}

	// Construir y escribir el MBR
	var mbr Structs.MBR
	mbr.MbrSize = int32(realSize)
	mbr.Signature = utils.GenerateRandomSignature()
	copy(mbr.Fit[:], fit)
	copy(mbr.CreationDate[:], utils.GetCurrentTimeString(10))

	if err := utils.WriteObject(file, mbr, 0); err != nil {
		return "", fmt.Errorf("error escribiendo MBR: %v", err)
	}

	Structs.PrintMBR(mbr, string(driveLetter))

	// Armar mensaje final
	unitStr := "KB"
	sizeStr := realSize / 1024
	if unit == "m" {
		unitStr = "MB"
		sizeStr = realSize / (1024 * 1024)
	}

	msg := fmt.Sprintf("Disco creado: %s (%d %s)", fileName, sizeStr, unitStr)
	return msg, nil
}
