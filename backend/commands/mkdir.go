package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/global"
	Structs "MIA_PI_202001151_1VAC1S2025/structs"
	"MIA_PI_202001151_1VAC1S2025/utils"
	"regexp"
	"strings"
)

func Fn_Mkdir(params string) {
	if globals.LoginSession.User == "" {
		utils.ShowMessage("Debe iniciar sesión primero.", true)
		return
	}
	paramDefs := map[string]Structs.ParamDef{
		"-path": {Required: true},
		"-r":    {Required: false, NotValue: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return
	}
	path := parsed["-path"]
	r := bool(parsed["-r"] == "true")
	if path == "" || path[0] != '/' {
		utils.ShowMessage("La ruta debe comenzar desde la raíz [/].", true)
		return
	}
	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		utils.ShowMessage("La partición de la sesión no está montada.", true)
		return
	}
	mkdir(path, r, *part)
}

// mkdir -path=<ruta> [-r]
func mkdir(path string, r bool, part Structs.Partition) {
	folders := strings.Split(path, "/")
	if len(folders) > 0 && folders[len(folders)-1] == "" {
		folders = folders[:len(folders)-1]
	}
	if len(folders) == 0 {
		utils.ShowMessage("Ruta inválida.", true)
		return
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_ ]{1,12}$`)
	for _, dir := range folders {
		if dir == "" {
			continue
		}
		if !re.MatchString(dir) {
			utils.ShowMessage("Los nombres de las carpetas deben tener entre 1 y 12 caracteres alfanuméricos, guion bajo o espacio: "+dir, true)
			return
		}
	}

	drive := strings.ToUpper(string(part.Id[0]))
	diskPath := globals.PathDisks + drive + ".dsk"
	file, err := utils.OpenFile(diskPath)
	if err != nil {
		utils.ShowMessage("No se pudo abrir el disco: "+diskPath, true)
		return
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		utils.ShowMessage("No se pudo leer el superbloque.", true)
		return
	}

	currentInode := int32(0)
	for idx, dir := range folders {
		if dir == "" {
			continue
		}
		exists, childInode := utils.SearchDirectoryEntry(file, &sb, currentInode, dir)
		if exists {
			currentInode = childInode
			continue
		}
		if !r && idx != len(folders)-1 {
			utils.ShowMessage("La carpeta ["+dir+"] no existe.\nUse -r para crearla junto a sus padres.", true)
			return
		}

		utils.WriteJournaling(sb, part, file, []byte("mkdir"), []byte(dir), []byte("-"))

		newInode, err := utils.CreateDirectory(file, &sb, currentInode, dir, part)
		if err != nil {
			utils.ShowMessage("Error al crear carpeta ["+dir+"]: "+err.Error(), true)
			return
		}
		currentInode = newInode
	}
	utils.ShowMessage("Carpeta creada/existente.", false)
}
