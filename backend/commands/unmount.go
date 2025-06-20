package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/global"
	Structs "MIA_PI_202001151_1VAC1S2025/structs"
	"MIA_PI_202001151_1VAC1S2025/utils"
	"fmt"
	"strings"
)

func Fn_Unmount(params string) {
	paramDefs := map[string]Structs.ParamDef{
		"-id": {Required: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		println(err.Error())
		return
	}
	id := strings.TrimSpace(parsed["-id"])
	Unmount(id)
}

// Mount -driveletter=<driveLetter> [-name=<name>]
func Unmount(id string) {
	driveLetter := strings.ToUpper(id[:1])
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
	var index int = -1
	for i, partition := range tempMBR.Partitions {
		if strings.Contains(string(partition.Id[:]), id) {
			if partition.Id == [4]byte{} {
				utils.ShowMessage("La partición con ID: "+id+" ya está desmontada.", true)
				return
			}
			index = i
			break
		}
	}

	if index == -1 {
		utils.ShowMessage("No se encontró la partición con ID: "+id, true)
		return
	}

	if index < 0 || index >= len(tempMBR.Partitions) {
		utils.ShowMessage("Índice de partición inválido: "+fmt.Sprint(index), true)
		return
	}
	tempMBR.Partitions[index].Id = [4]byte{}
	copy(tempMBR.Partitions[index].Status[:], "0")
	tempMBR.Partitions[index].Correlative = 0

	if err := utils.WriteObject(file, tempMBR, 0); err != nil {
		return
	}
	Structs.PrintPartitions(tempMBR, driveLetter)
	utils.ShowMessage("La partición ["+id+"] del disco ["+driveLetter+"] ha sido desmontada.", false)
}
