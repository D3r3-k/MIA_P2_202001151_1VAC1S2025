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

func Fn_Fdisk(params string) {
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
		fmt.Println(err)
		return
	}
	var size int
	driveLetter := strings.ToUpper(parsed["-driveletter"])
	name := strings.TrimSpace(parsed["-name"])
	unit := strings.ToLower(parsed["-unit"])
	type_ := strings.ToLower(parsed["-type"])
	fit := strings.ToLower(parsed["-fit"])
	if parsed["-delete"] != "" {
		FdiskDelete(parsed["-delete"], driveLetter, name)
		return
	}
	if parsed["-add"] != "" {
		var add int
		if _, err := fmt.Sscanf(parsed["-add"], "%d", &add); err != nil {
			utils.ShowMessage("El valor de -add debe ser numérico.", true)
			return
		}
		FdiskAdd(add, unit, driveLetter, name)
		return
	}
	if parsed["-size"] == "" {
		utils.ShowMessage("El parámetro -size es obligatorio.", true)
		return
	}
	if _, err := fmt.Sscanf(parsed["-size"], "%d", &size); err != nil {
		utils.ShowMessage("El valor de -size debe ser numérico.", true)
		return
	}
	Fdisk(size, driveLetter, name, type_, fit, unit)

}

// Fdisk -size=<size> -driveletter=<driveLetter> -name=<name> -type=<p|e> -fit=<b|f|w> -unit=<k|m>
func Fdisk(size int, driveLetter string, name string, type_ string, fit string, unit string) {
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
			utils.ShowMessage(valid.msg, true)
			return
		}
	}
	size = utils.GetRealSize(size, unit)

	file, err := utils.OpenFile(globals.PathDisks + driveLetter + ".dsk")
	if err != nil {
		return
	}
	defer file.Close()

	var tempMBR Structs.MBR
	if err := utils.ReadObject(file, &tempMBR, 0); err != nil {
		return
	}
	for i := 0; i < 4; i++ {
		existingName := strings.TrimRight(string(tempMBR.Partitions[i].Name[:]), "\x00")
		if tempMBR.Partitions[i].Size != 0 && existingName == name {
			utils.ShowMessage(fmt.Sprintf("Ya existe una partición con el nombre [%s] en el disco [%s].", name, driveLetter), true)
			return
		}
		if tempMBR.Partitions[i].Type[0] == 'e' {
			utils.ShowMessage(fmt.Sprintf("No se puede crear una partición extendida [%s] en el disco [%s].\nYa existe una partición extendida.", name, driveLetter), true)
			return
		}
	}

	freeSpaces := utils.GetFreeSpaces(&tempMBR, tempMBR.MbrSize)

	type spaceIdx struct {
		idx   int
		start int32
		size  int32
	}
	var candidates []spaceIdx

	for idx := 0; idx < 4; idx++ {
		if tempMBR.Partitions[idx].Size == 0 {
			for _, sp := range freeSpaces {
				if int32(size) <= sp[1] {
					candidates = append(candidates, spaceIdx{idx, sp[0], sp[1]})
				}
			}
		}
	}

	if len(candidates) == 0 {
		utils.ShowMessage(fmt.Sprintf("No hay espacio para la partición [%s].\nEspacio del disco [%s]: %d%s.", name, driveLetter, size, unit), true)
		return
	}

	_fitMbr := strings.ToLower(string(tempMBR.Fit[:]))
	var selected spaceIdx
	switch _fitMbr {
	case "f":
		// First fit: primer espacio donde quepa
		selected = candidates[0]
	case "b":
		// Best fit: espacio más pequeño posible
		best := candidates[0]
		for _, c := range candidates {
			if c.size < best.size {
				best = c
			}
		}
		selected = best
	case "w":
		// Worst fit: espacio más grande posible
		worst := candidates[0]
		for _, c := range candidates {
			if c.size > worst.size {
				worst = c
			}
		}
		selected = worst
	default:
		utils.ShowMessage(fmt.Sprintf("Ajuste [%s] no reconocido. Debe ser <ff|bf|wf>.\n", _fitMbr), true)
		return
	}

	// Asignar partición
	p := &tempMBR.Partitions[selected.idx]
	p.Size = int32(size)
	copy(p.Name[:], name)
	copy(p.Fit[:], fit)
	copy(p.Status[:], "0")
	copy(p.Type[:], type_)
	p.Start = selected.start

	// Guardar cambios
	if err := utils.WriteObject(file, &tempMBR, 0); err != nil {
		return
	}

	if type_ == "e" {
		type_ = "Extendida"
	} else if type_ == "p" {
		type_ = "Primaria"
	}
	if unit == "b" {
		unit = " bytes"
	} else if unit == "k" {
		unit = " Kb"
	} else if unit == "m" {
		unit = " Mb"
	}

	Structs.PrintMBR(tempMBR, driveLetter)
	utils.ShowMessage(fmt.Sprintf("Partición [%s] - tipo [%s].\nCreada exitosamente en el disco [%s] %d%s.", name, type_, driveLetter, size, unit), false)
}

// Fdisk -add=<add> -unit=<unit> -driveletter=<driveLetter> -name=<name>
func FdiskAdd(add int, unit string, driveLetter string, name string) {
	if unit != "b" && unit != "k" && unit != "m" {
		utils.ShowMessage("La unidad debe ser <b|k|m>.", true)
		return
	}
	// Convertir tamaño a bytes
	size := utils.GetRealSize(add, unit)

	path := globals.PathDisks + driveLetter + ".dsk"
	file, err := utils.OpenFile(path)
	if err != nil {
		return
	}
	defer file.Close()

	var tempMBR Structs.MBR
	if err := utils.ReadObject(file, &tempMBR, 0); err != nil {
		return
	}

	found := false
	for i := 0; i < 4; i++ {
		part := &tempMBR.Partitions[i]
		existingName := strings.TrimRight(string(part.Name[:]), "\x00")
		if part.Size > 0 && existingName == name {
			found = true
			if add > 0 {
				end := part.Start + part.Size
				newEnd := end + int32(size)
				for j := 0; j < 4; j++ {
					if i == j || tempMBR.Partitions[j].Size == 0 {
						continue
					}
					startJ := tempMBR.Partitions[j].Start
					if startJ > end && startJ < newEnd {
						utils.ShowMessage(fmt.Sprintf("No se puede expandir: colisión con partición [%s].", strings.TrimRight(string(tempMBR.Partitions[j].Name[:]), "\x00")), true)
						return
					}
				}
				if newEnd > tempMBR.MbrSize {
					utils.ShowMessage(fmt.Sprintf("No se puede expandir: se excede el tamaño del disco [%s].", driveLetter), true)
					return
				}
				part.Size += int32(size)

			} else if add < 0 {
				if int32(-size) > part.Size {
					utils.ShowMessage(fmt.Sprintf("No se puede reducir la partición [%s] más allá de su tamaño actual.", existingName), true)
					return
				}
				oldEnd := part.Start + part.Size
				newEnd := oldEnd + int32(size)
				zeroBuf := make([]byte, -size)
				file.WriteAt(zeroBuf, int64(newEnd))
				part.Size += int32(size)
			}
			break
		}
	}
	if !found {
		utils.ShowMessage(fmt.Sprintf("No se encontró una partición con el nombre [%s] en el disco [%s].", name, driveLetter), true)
		return
	}
	if err := utils.WriteObject(file, &tempMBR, 0); err != nil {
		return
	}

	Structs.PrintMBR(tempMBR, driveLetter)
	utils.ShowMessage(fmt.Sprintf("Partición [%s] actualizada exitosamente en el disco [%s].", name, driveLetter), false)
}

// Fdisk -delete=full -<driveLetter>=<driveLetter> -name=<name>
func FdiskDelete(delete string, driveLetter string, name string) {
	if delete != "full" {
		utils.ShowMessage("Error en el parametro: fdisk -<delete=full>.", true)
		return
	}
	file, err := utils.OpenFile(globals.PathDisks + driveLetter + ".dsk")
	if err != nil {
		utils.ShowMessage(fmt.Sprintf("No se pudo abrir el disco [%s]: %v", driveLetter, err), true)
		return
	}
	defer file.Close()

	utils.ShowMessage("Desea eliminar el disco "+driveLetter+"?\nS: para reemplazar\nN: para cancelar", true)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">>> ")
	scanner.Scan()
	response := strings.TrimSpace(scanner.Text())
	if strings.ToUpper(response) != "S" {
		utils.ShowMessage("Operación cancelada por el usuario.", false)
		return
	}

	var tempMBR Structs.MBR
	if err := utils.ReadObject(file, &tempMBR, 0); err != nil {
		return
	}
	found := false
	for i := 0; i < 4; i++ {
		existingName := strings.TrimRight(string(tempMBR.Partitions[i].Name[:]), "\x00")
		if tempMBR.Partitions[i].Size != 0 && existingName == name {
			tempMBR.Partitions[i] = Structs.Partition{}
			found = true
			break
		}
	}
	if !found {
		utils.ShowMessage(fmt.Sprintf("No se encontró una partición con el nombre [%s] en el disco [%s].", name, driveLetter), true)
		return
	}
	if err := utils.WriteObject(file, &tempMBR, 0); err != nil {
		return
	}
	Structs.PrintMBR(tempMBR, driveLetter)
	utils.ShowMessage(fmt.Sprintf("Partición [%s] eliminada exitosamente del disco [%s].", name, driveLetter), false)
}
