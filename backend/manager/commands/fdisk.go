package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"errors"
	"fmt"
	"strings"
)

func Fn_Fdisk(params string) (string, error) {
	paramDefs := map[string]Structs.ParamDef{
		"-size":        {Required: false},
		"-driveletter": {Required: true},
		"-name":        {Required: true},
		"-unit":        {Required: false, Default: "k"},
		"-type":        {Required: false, Default: "p"},
		"-fit":         {Required: false, Default: "wf"},
		"-delete":      {Required: false, Default: ""},
		"-add":         {Required: false, Default: ""},
	}

	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		return "", err
	}

	driveLetter := strings.ToUpper(parsed["-driveletter"])
	name := strings.TrimSpace(parsed["-name"])
	unit := strings.ToLower(parsed["-unit"])
	type_ := strings.ToLower(parsed["-type"])
	fit := strings.ToLower(parsed["-fit"])

	if parsed["-delete"] != "" {
		return FdiskDelete(parsed["-delete"], driveLetter, name)
	}

	if parsed["-add"] != "" {
		var add int
		if _, err := fmt.Sscanf(parsed["-add"], "%d", &add); err != nil {
			return "", errors.New("el valor de -add debe ser numérico")
		}
		return FdiskAdd(add, unit, driveLetter, name)
	}

	if parsed["-size"] == "" {
		return "", errors.New("el parámetro -size es obligatorio")
	}

	var size int
	if _, err := fmt.Sscanf(parsed["-size"], "%d", &size); err != nil {
		return "", errors.New("el valor de -size debe ser numérico")
	}

	return Fdisk(size, driveLetter, name, type_, fit, unit)
}

// Fdisk -size=<size> -driveletter=<driveLetter> -name=<name> -type=<p|e> -fit=<b|f|w> -unit=<k|m>
func Fdisk(size int, driveLetter, name, type_, fit, unit string) (string, error) {
	validaciones := map[string]struct {
		ok  bool
		msg string
	}{
		"size":        {size > 0, "El tamaño debe ser mayor a 0."},
		"driveLetter": {len(driveLetter) == 1 && driveLetter[0] >= 'A' && driveLetter[0] <= 'Z', "La letra del disco debe ser una sola letra entre A y Z."},
		"name":        {len(name) > 0 && utils.ValidateRegex(name, "^[A-Za-z0-9]+$"), "El nombre debe contener solo letras y números y no puede estar vacío."},
		"type":        {type_ == "p" || type_ == "e", "El tipo debe ser 'p' (primaria) o 'e' (extendida)."},
		"fit":         {fit == "bf" || fit == "ff" || fit == "wf", "El ajuste debe ser 'bf', 'ff' o 'wf'."},
		"unit":        {unit == "b" || unit == "k" || unit == "m", "La unidad debe ser 'b', 'k' o 'm'."},
	}

	for _, valid := range validaciones {
		if !valid.ok {
			return "", errors.New(valid.msg)
		}
	}

	size = utils.GetRealSize(size, unit)
	file, err := utils.OpenFile(globals.PathDisks + driveLetter + ".dsk")
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el disco [%s]: %v", driveLetter, err)
	}
	defer file.Close()

	var mbr Structs.MBR
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		return "", err
	}

	for i := 0; i < 4; i++ {
		existing := strings.TrimRight(string(mbr.Partitions[i].Name[:]), "\x00")
		if mbr.Partitions[i].Size != 0 && existing == name {
			return "", fmt.Errorf("ya existe una partición con el nombre [%s] en el disco [%s]", name, driveLetter)
		}
		if mbr.Partitions[i].Type[0] == 'e' && type_ == "e" {
			return "", fmt.Errorf("ya existe una partición extendida en el disco [%s]", driveLetter)
		}
	}

	freeSpaces := utils.GetFreeSpaces(&mbr, mbr.MbrSize)
	type spaceIdx struct {
		idx   int
		start int32
		size  int32
	}
	var candidates []spaceIdx

	for i := 0; i < 4; i++ {
		if mbr.Partitions[i].Size == 0 {
			for _, sp := range freeSpaces {
				if int32(size) <= sp[1] {
					candidates = append(candidates, spaceIdx{i, sp[0], sp[1]})
				}
			}
		}
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("no hay espacio para la partición [%s] en el disco [%s]", name, driveLetter)
	}

	_fit := strings.ToLower(string(mbr.Fit[:]))
	var selected spaceIdx
	switch _fit {
	case "f":
		selected = candidates[0]
	case "b":
		selected = candidates[0]
		for _, c := range candidates {
			if c.size < selected.size {
				selected = c
			}
		}
	case "w":
		selected = candidates[0]
		for _, c := range candidates {
			if c.size > selected.size {
				selected = c
			}
		}
	default:
		return "", fmt.Errorf("ajuste [%s] no reconocido. Debe ser <ff|bf|wf>", _fit)
	}

	p := &mbr.Partitions[selected.idx]
	p.Size = int32(size)
	copy(p.Name[:], name)
	copy(p.Fit[:], fit)
	copy(p.Status[:], "0")
	copy(p.Type[:], type_)
	p.Start = selected.start

	if err := utils.WriteObject(file, &mbr, 0); err != nil {
		return "", err
	}

	unitLabel := map[string]string{"b": "bytes", "k": "KB", "m": "MB"}[unit]
	Structs.PrintMBR(mbr, driveLetter)
	return fmt.Sprintf("Partición [%s] de tipo [%s] creada en disco [%s] con tamaño %d%s", name, type_, driveLetter, size, unitLabel), nil
}

// Fdisk -add=<add> -unit=<unit> -driveletter=<driveLetter> -name=<name>
func FdiskAdd(add int, unit, driveLetter, name string) (string, error) {
	if unit != "b" && unit != "k" && unit != "m" {
		return "", errors.New("la unidad debe ser <b|k|m>")
	}

	size := utils.GetRealSize(add, unit)
	path := globals.PathDisks + driveLetter + ".dsk"

	file, err := utils.OpenFile(path)
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el disco [%s]: %v", driveLetter, err)
	}
	defer file.Close()

	var mbr Structs.MBR
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		return "", err
	}

	found := false
	var sizeBefore, sizeAfter int32

	for i := 0; i < 4; i++ {
		part := &mbr.Partitions[i]
		existingName := strings.TrimRight(string(part.Name[:]), "\x00")

		if part.Size > 0 && existingName == name {
			found = true
			sizeBefore = part.Size

			if add > 0 {
				end := part.Start + part.Size
				newEnd := end + int32(size)

				for j := 0; j < 4; j++ {
					if i == j || mbr.Partitions[j].Size == 0 {
						continue
					}
					startJ := mbr.Partitions[j].Start
					if startJ > end && startJ < newEnd {
						conflict := strings.TrimRight(string(mbr.Partitions[j].Name[:]), "\x00")
						return "", fmt.Errorf("no se puede expandir: colisión con partición [%s]", conflict)
					}
				}

				if newEnd > mbr.MbrSize {
					return "", fmt.Errorf("no se puede expandir: se excede el tamaño del disco [%s]", driveLetter)
				}

				part.Size += int32(size)

			} else if add < 0 {
				if int32(-size) > part.Size {
					return "", fmt.Errorf("no se puede reducir la partición [%s] más allá de su tamaño actual", existingName)
				}

				oldEnd := part.Start + part.Size
				newEnd := oldEnd + int32(size) // size negativo
				zeroBuf := make([]byte, -size)
				if _, err := file.WriteAt(zeroBuf, int64(newEnd)); err != nil {
					return "", fmt.Errorf("error al limpiar espacio al reducir partición: %v", err)
				}
				part.Size += int32(size)
			}

			sizeAfter = part.Size
			break
		}
	}

	if !found {
		return "", fmt.Errorf("no se encontró una partición con el nombre [%s] en el disco [%s]", name, driveLetter)
	}

	if err := utils.WriteObject(file, &mbr, 0); err != nil {
		return "", fmt.Errorf("no se pudo escribir en el disco: %v", err)
	}

	Structs.PrintMBR(mbr, driveLetter)

	msg := fmt.Sprintf(
		"Partición [%s] del disco [%s] modificada: tamaño anterior = %d bytes, nuevo tamaño = %d bytes",
		name, driveLetter, sizeBefore, sizeAfter,
	)
	return msg, nil
}

// Fdisk -delete=full -<driveLetter>=<driveLetter> -name=<name>
func FdiskDelete(delete string, driveLetter string, name string) (string, error) {
	if delete != "full" {
		return "", errors.New("el parámetro -delete debe ser 'full'")
	}

	path := globals.PathDisks + driveLetter + ".dsk"
	file, err := utils.OpenFile(path)
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el disco [%s]: %v", driveLetter, err)
	}
	defer file.Close()

	var mbr Structs.MBR
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		return "", err
	}

	found := false
	for i := 0; i < 4; i++ {
		existingName := strings.TrimRight(string(mbr.Partitions[i].Name[:]), "\x00")
		if mbr.Partitions[i].Size != 0 && existingName == name {
			mbr.Partitions[i] = Structs.Partition{}
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("no se encontró una partición con el nombre [%s] en el disco [%s]", name, driveLetter)
	}

	if err := utils.WriteObject(file, &mbr, 0); err != nil {
		return "", fmt.Errorf("no se pudo guardar el MBR actualizado en el disco [%s]: %v", driveLetter, err)
	}

	Structs.PrintMBR(mbr, driveLetter)

	return fmt.Sprintf("Partición [%s] eliminada exitosamente del disco [%s]", name, driveLetter), nil
}
