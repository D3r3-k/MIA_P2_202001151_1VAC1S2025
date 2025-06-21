package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"fmt"
	"sort"
	"strings"
)

func Fn_Mount(params string) {
	paramDefs := map[string]Structs.ParamDef{
		"-driveletter": {Required: true},
		"-name":        {Required: false},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		println(err.Error())
		return
	}
	driveLetter := strings.ToUpper(parsed["-driveletter"])
	name := strings.TrimSpace(parsed["-name"])
	Mount(driveLetter, name)
}

// Mount -driveletter=<driveLetter> [-name=<name>]
func Mount(driveLetter string, name string) {
	path := globals.PathDisks + driveLetter + ".dsk"
	file, err := utils.OpenFile(path)
	if err != nil {
		utils.ShowMessage(fmt.Sprintf("No se pudo abrir el disco [%s]: %v", driveLetter, err), true)
		return
	}
	defer file.Close()

	var tempMBR Structs.MBR
	if err := utils.ReadObject(file, &tempMBR, 0); err != nil {
		return
	}
	var index int = -1
	var emptyId [4]byte
	correlativos := []int{}
	for i := 0; i < 4; i++ {
		part := tempMBR.Partitions[i]
		if strings.Contains(string(part.Name[:]), name) {
			if part.Id != emptyId {
				utils.ShowMessage(fmt.Sprintf("La partición [%s] ya está montada en la unidad [%s].", name, driveLetter), true)
				return
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
		utils.ShowMessage(fmt.Sprintf("No se encontró una partición con el nombre [%s] en el disco [%s].", name, driveLetter), true)
		return
	}

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

	copy(tempMBR.Partitions[index].Id[:], id)
	copy(tempMBR.Partitions[index].Status[:], "1")
	tempMBR.Partitions[index].Correlative = int32(correlativo)

	if err := utils.WriteObject(file, tempMBR, 0); err != nil {
		return
	}
	Structs.PrintPartitions(tempMBR, driveLetter)
	utils.ShowMessage(fmt.Sprintf("Partición [%s] montada exitosamente en la unidad [%s].", name, driveLetter), false)
}
