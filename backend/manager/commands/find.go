package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type FindResponse struct {
	Name     string         `json:"Name"`
	Type     string         `json:"Type"` // Folder o File
	Path     string         `json:"Path"`
	Size     string         `json:"Size,omitempty"`
	Children []FindResponse `json:"Children,omitempty"`
}

func Fn_Find(params string) ([]FindResponse, error) {
	if globals.LoginSession.User == "" {
		return nil, fmt.Errorf("debe iniciar sesión primero")
	}

	paramDefs := map[string]Structs.ParamDef{
		"-path": {Required: true},
		"-name": {Required: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		return nil, err
	}

	startPath := parsed["-path"]
	namePattern := parsed["-name"]

	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		return nil, fmt.Errorf("la partición de la sesión no está montada")
	}

	regexPattern := "^" + strings.ReplaceAll(strings.ReplaceAll(namePattern, ".", `\.`), "*", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "?", ".") + "$"
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, fmt.Errorf("patrón inválido: %s", namePattern)
	}

	drive := strings.ToUpper(string(part.Id[0]))
	diskPath := globals.PathDisks + drive + ".dsk"
	file, err := utils.OpenFile(diskPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el disco: %s", diskPath)
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return nil, fmt.Errorf("no se pudo leer el superbloque")
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
				return nil, fmt.Errorf("no existe el directorio: %s", dir)
			}
			inodoInicio = childInode
		}
	}

	result := []FindResponse{}
	rootResp := buildFindTree(file, &sb, inodoInicio, startPath, regex)
	if rootResp != nil {
		result = append(result, *rootResp)
	}
	return result, nil
}

func buildFindTree(file *os.File, sb *Structs.Superblock, inodeNum int32, path string, regex *regexp.Regexp) *FindResponse {
	inodoSize := int32(binary.Size(Structs.Inode{}))
	blockSize := int32(binary.Size(Structs.Folderblock{}))

	var inode Structs.Inode
	if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start+inodeNum*inodoSize)); err != nil {
		return nil
	}

	base := path
	if strings.HasSuffix(base, "/") && len(base) > 1 {
		base = base[:len(base)-1]
	}
	nodeName := base[strings.LastIndex(base, "/")+1:]
	if nodeName == "" {
		nodeName = "/"
	}

	// Evaluar si debe incluirse por nombre
	if nodeName != "/" && !regex.MatchString(nodeName) {
		return nil
	}

	resp := FindResponse{
		Name: nodeName,
		Path: path,
		Type: "File",
	}

	if inode.I_type[0] == '0' {
		resp.Type = "Folder"
	} else {
		resp.Size = strconv.Itoa(int(inode.I_size)) + " B"
	}

	// Si es carpeta, buscar hijos
	if inode.I_type[0] == '0' {
		var hijos []FindResponse
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
				childPath := path
				if !strings.HasSuffix(childPath, "/") {
					childPath += "/"
				}
				childPath += nombre
				child := buildFindTree(file, sb, entry.B_inodo, childPath, regex)
				if child != nil {
					hijos = append(hijos, *child)
				}
			}
		}
		if len(hijos) > 0 {
			resp.Children = hijos
		}
	}
	return &resp
}
