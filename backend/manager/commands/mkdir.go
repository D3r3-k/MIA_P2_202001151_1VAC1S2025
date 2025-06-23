package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"fmt"
	"regexp"
	"strings"
)

func Fn_Mkdir(params string) (string, error) {
	if globals.LoginSession.User == "" {
		utils.ShowMessage("Debe iniciar sesión primero.", true)
		return "", fmt.Errorf("debe iniciar sesión primero")
	}
	paramDefs := map[string]Structs.ParamDef{
		"-path": {Required: true},
		"-r":    {Required: false, NotValue: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return "", err
	}
	path := parsed["-path"]
	r := bool(parsed["-r"] == "true")
	if path == "" || path[0] != '/' {
		utils.ShowMessage("La ruta debe comenzar desde la raíz [/].", true)
		return "", fmt.Errorf("la ruta debe comenzar desde la raíz [/]")
	}
	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		utils.ShowMessage("La partición de la sesión no está montada.", true)
		return "", fmt.Errorf("la partición de la sesión no está montada")
	}
	return mkdir(path, r, *part)
}

// mkdir -path=<ruta> [-r]
func mkdir(path string, r bool, part Structs.Partition) (string, error) {
	folders := strings.Split(path, "/")
	if len(folders) > 0 && folders[len(folders)-1] == "" {
		folders = folders[:len(folders)-1]
	}
	if len(folders) == 0 {
		return "", fmt.Errorf("ruta inválida")
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9_ ]{1,12}$`)
	for _, dir := range folders {
		if dir == "" {
			continue
		}
		if !re.MatchString(dir) {
			return "", fmt.Errorf("los nombres de las carpetas deben tener entre 1 y 12 caracteres alfanuméricos, guion bajo o espacio: %s", dir)
		}
	}

	drive := strings.ToUpper(string(part.Id[0]))
	diskPath := globals.PathDisks + drive + ".dsk"
	file, err := utils.OpenFile(diskPath)
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el disco: %s", diskPath)
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return "", fmt.Errorf("no se pudo leer el superbloque: %v", err)
	}

	currentInode := int32(0)
	var lastCreated string
	for idx, dir := range folders {
		if dir == "" {
			continue
		}
		exists, childInode := utils.SearchDirectoryEntry(file, &sb, currentInode, dir)
		if exists {
			currentInode = childInode
			lastCreated = "" // no fue creada
			continue
		}
		if !r && idx != len(folders)-1 {
			return "", fmt.Errorf("la carpeta [%s] no existe, use -r para crearla junto a sus padres", dir)
		}

		utils.WriteJournaling(sb, part, file, []byte("mkdir"), []byte(dir), []byte("-"))

		newInode, err := utils.CreateDirectory(file, &sb, currentInode, dir, part)
		if err != nil {
			return "", fmt.Errorf("error al crear carpeta [%s]: %v", dir, err)
		}
		currentInode = newInode
		lastCreated = dir
	}
	if lastCreated == "" {
		return fmt.Sprintf("La carpeta [%s] ya existe.", folders[len(folders)-1]), nil
	}
	return fmt.Sprintf("Carpeta [%s] creada exitosamente.", lastCreated), nil
}
