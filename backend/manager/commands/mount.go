package commands

import (
	"fmt"
	"sort"
	"strings"

	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
)

func Fn_Mount(params string) (string, error) {
	paramDefs := map[string]Structs.ParamDef{
		"-driveletter": {Required: true},
		"-name":        {Required: false},
	}

	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		return "", err
	}

	driveLetter := strings.ToUpper(parsed["-driveletter"])
	name := strings.TrimSpace(parsed["-name"])
	return Mount(driveLetter, name)
}

// Mount -driveletter=<driveLetter> [-name=<name>]
func Mount(driveLetter string, name string) (string, error) {
	path := globals.PathDisks + driveLetter + ".dsk"
	file, err := utils.OpenFile(path)
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el disco [%s]: %v", driveLetter, err)
	}
	defer file.Close()

	var mbr Structs.MBR
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		return "", fmt.Errorf("error leyendo el MBR del disco [%s]: %v", driveLetter, err)
	}

	var index int = -1
	var emptyId [4]byte
	var correlativos []int

	for i := 0; i < 4; i++ {
		part := mbr.Partitions[i]
		partName := strings.TrimRight(string(part.Name[:]), "\x00")

		if partName == name {
			if part.Id != emptyId {
				return "", fmt.Errorf("la partición [%s] ya está montada en la unidad [%s]", name, driveLetter)
			}
			index = i
		}

		idStr := string(part.Id[:])
		if len(idStr) >= 3 && part.Id != emptyId {
			correlativo := int(idStr[1] - '0')
			correlativos = append(correlativos, correlativo)
		}
	}

	if index == -1 {
		return "", fmt.Errorf("no se encontró una partición con el nombre [%s] en el disco [%s]", name, driveLetter)
	}

	// Calcular nuevo correlativo disponible
	correlativo := 1
	if len(correlativos) > 0 {
		sort.Ints(correlativos)
		for i, v := range correlativos {
			if v != i+1 {
				correlativo = i + 1
				break
			}
			if i == len(correlativos)-1 {
				correlativo = v + 1
			}
		}
	}

	id := strings.ToUpper(driveLetter) + fmt.Sprintf("%d", correlativo) + "51"

	copy(mbr.Partitions[index].Id[:], id)
	copy(mbr.Partitions[index].Status[:], "1")
	mbr.Partitions[index].Correlative = int32(correlativo)

	if err := utils.WriteObject(file, mbr, 0); err != nil {
		return "", fmt.Errorf("no se pudo guardar el estado de montaje: %v", err)
	}

	Structs.PrintPartitions(mbr, driveLetter)

	return fmt.Sprintf("Partición [%s] montada exitosamente en la unidad [%s] con ID [%s]", name, driveLetter, id), nil
}
