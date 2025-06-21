package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func Fn_Find(params string) {
	if globals.LoginSession.User == "" {
		utils.ShowMessage("Debe iniciar sesión primero.", true)
		return
	}
	paramDefs := map[string]Structs.ParamDef{
		"-path": {Required: true},
		"-name": {Required: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return
	}
	startPath := parsed["-path"]
	namePattern := parsed["-name"]

	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		utils.ShowMessage("La partición de la sesión no está montada.", true)
		return
	}

	// Convertir patrón de nombre a regex
	regexPattern := "^" + strings.ReplaceAll(strings.ReplaceAll(namePattern, ".", `\.`), "*", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "?", ".")
	regexPattern += "$"
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		utils.ShowMessage("Patrón inválido: "+namePattern, true)
		return
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
	inodoInicio := int32(0)
	if startPath != "/" {
		parts := strings.Split(startPath, "/")
		for _, dir := range parts {
			if dir == "" {
				continue
			}
			ok, childInode := utils.SearchDirectoryEntry(file, &sb, inodoInicio, dir)
			if !ok {
				utils.ShowMessage("No existe el directorio: "+dir, true)
				return
			}
			inodoInicio = childInode
		}
	}
	// ENCABEZADO: FIND
	fmt.Println("╔═══════════════════════[Busqueda]══════════════════════════╗")
	find(file, &sb, inodoInicio, "/", regex, 0, "", true)
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
}

// Recursivo: muestra solo los archivos/carpetas que cumplan el patrón
func find(file *os.File, sb *Structs.Superblock, inodeNum int32, path string, regex *regexp.Regexp, depth int, prefix string, isLast bool) {
	inodoSize := int32(binary.Size(Structs.Inode{}))
	blockSize := int32(binary.Size(Structs.Folderblock{}))
	var inode Structs.Inode
	if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start+inodeNum*inodoSize)); err != nil {
		return
	}

	base := path
	if strings.HasSuffix(base, "/") && len(base) > 1 {
		base = base[:len(base)-1]
	}
	nodeName := base[strings.LastIndex(base, "/")+1:]

	show := inode.I_type[0] == '0' || regex.MatchString(nodeName)
	branch := ""
	if depth == 0 {
		fmt.Printf("║ %-40s#%s\n", "/", string(inode.I_perm[:3]))
	} else if show {
		if isLast {
			branch = "└── "
		} else {
			branch = "├── "
		}
		name := strings.Trim(nodeName, "\x00")
		line := prefix + branch + name
		output := fmt.Sprintf("%-40s#%s", line, string(inode.I_perm[:3]))
		fmt.Printf("║ %-57s \n", output)
	}
	if inode.I_type[0] != '0' {
		return
	}
	type hijo struct {
		nombre string
		inodo  int32
	}
	var hijos []hijo
	for i := 0; i < 14; i++ {
		blockNum := inode.I_block[i]
		if blockNum == -1 {
			continue
		}
		var fb Structs.Folderblock
		offset := int64(sb.S_block_start + blockNum*blockSize)
		if err := utils.ReadObject(file, &fb, offset); err != nil {
			continue
		}
		for _, entry := range fb.B_content {
			nombre := strings.Trim(string(entry.B_name[:]), "\x00")
			if nombre == "" || nombre == "." || nombre == ".." || entry.B_inodo == -1 {
				continue
			}
			hijos = append(hijos, hijo{nombre, entry.B_inodo})
		}
	}
	for i, h := range hijos {
		newPrefix := prefix
		if depth > 0 {
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
		}
		last := (i == len(hijos)-1)
		find(file, sb, h.inodo, path+h.nombre+"/", regex, depth+1, newPrefix, last)
	}
}
