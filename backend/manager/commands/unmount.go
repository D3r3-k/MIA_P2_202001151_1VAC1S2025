package commands

import (
	"errors"
	"fmt"
	"strings"

	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
)

func Fn_Unmount(params string) (string, error) {
	paramDefs := map[string]Structs.ParamDef{
		"-id": {Required: true},
	}

	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		return "", err
	}

	id := strings.TrimSpace(parsed["-id"])
	return Unmount(id)
}

func Unmount(id string) (string, error) {
	if len(id) < 2 {
		return "", errors.New("el ID proporcionado es inválido")
	}

	driveLetter := strings.ToUpper(id[:1])
	path := globals.PathDisks + driveLetter + ".dsk"

	file, err := utils.OpenFile(path)
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el disco [%s]: %v", driveLetter, err)
	}
	defer file.Close()

	var mbr Structs.MBR
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		return "", fmt.Errorf("no se pudo leer el MBR del disco [%s]: %v", driveLetter, err)
	}

	index := -1
	for i, partition := range mbr.Partitions {
		partID := strings.TrimRight(string(partition.Id[:]), "\x00")
		if partID == id {
			if partition.Id == [4]byte{} {
				return "", fmt.Errorf("la partición con ID [%s] ya está desmontada", id)
			}
			index = i
			break
		}
	}

	if index == -1 {
		return "", fmt.Errorf("no se encontró una partición con ID [%s] en el disco [%s]", id, driveLetter)
	}

	// Limpiar datos de montaje
	mbr.Partitions[index].Id = [4]byte{}
	copy(mbr.Partitions[index].Status[:], "0")
	mbr.Partitions[index].Correlative = 0

	if err := utils.WriteObject(file, mbr, 0); err != nil {
		return "", fmt.Errorf("no se pudo escribir los cambios al disco [%s]: %v", driveLetter, err)
	}

	Structs.PrintPartitions(mbr, driveLetter)
	return fmt.Sprintf("Partición [%s] desmontada exitosamente del disco [%s]", id, driveLetter), nil
}
