package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/global"
	Structs "MIA_PI_202001151_1VAC1S2025/structs"
	"MIA_PI_202001151_1VAC1S2025/utils"
	"fmt"
	"os"
	"strings"
)

func Fn_Mkdisk(params string) {
	paramDefs := map[string]Structs.ParamDef{
		"-size": {Required: true},
		"-fit":  {Required: false, Default: "ff"},
		"-unit": {Required: false, Default: "m"},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return
	}
	var size int
	if _, err := fmt.Sscanf(parsed["-size"], "%d", &size); err != nil {
		return
	}
	fit := strings.ToLower(parsed["-fit"])
	unit := strings.ToLower(parsed["-unit"])
	Mkdisk(size, fit, unit)
}

// mkdisk -size=<size> -fit=<ff|bf|wf> -unit=<k|m>
func Mkdisk(size int, fit, unit string) {
	validaciones := map[string]struct {
		ok  bool
		msg string
	}{
		"size": {size > 0, "El tamaño debe ser mayor a 0."},
		"fit":  {fit == "ff" || fit == "bf" || fit == "wf", "El ajuste debe ser <'ff'|'bf'|'wf'>."},
		"unit": {unit == "k" || unit == "m", "La unidad debe ser <'k'|'m'>."},
	}
	for _, valid := range validaciones {
		if !valid.ok {
			utils.ShowMessage(valid.msg, true)
			return
		}
	}

	// Calcular tamaño real en bytes
	realSize := size
	if unit == "k" {
		realSize *= 1024
	} else {
		realSize *= 1024 * 1024
	}

	// Buscar un nombre de disco disponible
	var fileName string
	var driveletter rune
	for c := 'A'; c <= 'Z'; c++ {
		fileName = fmt.Sprintf(globals.PathDisks+"%c.dsk", c)
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			driveletter = c
			break
		}
	}

	// Crear archivo de disco con 0 binarios
	if err := utils.CreateFile(fileName); err != nil {
		return
	}
	file, err := utils.OpenFile(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	zeroBuffer := make([]byte, 1024)
	for i := 0; i < realSize/1024; i++ {
		if err := utils.WriteObject(file, zeroBuffer, int64(i*1024)); err != nil {
			return
		}
	}

	var mbr Structs.MBR
	mbr.MbrSize = int32(realSize)
	mbr.Signature = utils.GenerateRandomSignature()
	copy(mbr.Fit[:], fit)
	copy(mbr.CreationDate[:], utils.GetCurrentTimeString(10))
	if err := utils.WriteObject(file, mbr, 0); err != nil {
		return
	}

	Structs.PrintMBR(mbr, string(driveletter))
	unitStr := ""
	sizeStr := 0
	if unit == "m" {
		sizeStr = realSize / (1024 * 1024)
		unitStr = "MB"
	} else {
		sizeStr = realSize / 1024
		unitStr = "KB"
	}
	utils.ShowMessage(fmt.Sprintf("Disco creado: %s (%d %s)", fileName, sizeStr, unitStr), false)
}
