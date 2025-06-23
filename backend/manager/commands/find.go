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
	Name        string         `json:"Name"`
	Type        string         `json:"Type"` // Folder o File
	Path        string         `json:"Path"`
	Size        string         `json:"Size,omitempty"`
	Permissions string         `json:"Permissions"`
	Children    []FindResponse `json:"Children,omitempty"`
}

type FindResult struct {
	Tree   string        // Árbol representado como string
	Object *FindResponse // Objeto raíz estructurado
	Error  error         // Error si ocurre
}

func Fn_Find(params string) FindResult {
	if globals.LoginSession.User == "" {
		return FindResult{Error: fmt.Errorf("debe iniciar sesión primero")}
	}

	paramDefs := map[string]Structs.ParamDef{
		"-path": {Required: true},
		"-name": {Required: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		return FindResult{Error: err}
	}

	startPath := parsed["-path"]
	namePattern := parsed["-name"]

	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		return FindResult{Error: fmt.Errorf("la partición de la sesión no está montada")}
	}

	// Preparar regex
	regexPattern := "^" + strings.ReplaceAll(strings.ReplaceAll(namePattern, ".", `\.`), "*", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "?", ".") + "$"
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return FindResult{Error: fmt.Errorf("patrón inválido: %s", namePattern)}
	}

	// Abrir disco
	drive := strings.ToUpper(string(part.Id[0]))
	diskPath := globals.PathDisks + drive + ".dsk"
	file, err := utils.OpenFile(diskPath)
	if err != nil {
		return FindResult{Error: fmt.Errorf("no se pudo abrir el disco: %s", diskPath)}
	}
	defer file.Close()

	// Leer superbloque
	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return FindResult{Error: fmt.Errorf("no se pudo leer el superbloque")}
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
				return FindResult{Error: fmt.Errorf("no existe el directorio: %s", dir)}
			}
			inodoInicio = childInode
		}
	}

	root := buildFindTree(file, &sb, inodoInicio, startPath, regex)
	if root == nil {
		return FindResult{Tree: "", Object: nil, Error: nil}
	}

	treeStr := buildTreeString(*root, "", true)
	return FindResult{Tree: treeStr, Object: root, Error: nil}
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

	if nodeName != "/" && !regex.MatchString(nodeName) {
		return nil
	}

	resp := FindResponse{
		Name:        nodeName,
		Path:        path,
		Permissions: string(inode.I_perm[:3]),
		Type:        "File",
	}

	if inode.I_type[0] == '0' {
		resp.Type = "Folder"
	} else {
		resp.Size = strconv.Itoa(int(inode.I_size)) + " B"
	}

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
				childPath := strings.TrimSuffix(path, "/") + "/" + nombre
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

func buildTreeString(node FindResponse, prefix string, isLast bool) string {
	branch := "└── "
	if !isLast {
		branch = "├── "
	}
	line := fmt.Sprintf("%s%s%s\t\t[%s]", prefix, branch, node.Name, node.Permissions)
	lines := []string{line}

	for i, child := range node.Children {
		childPrefix := prefix
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
		lines = append(lines, buildTreeString(child, childPrefix, i == len(child.Children)-1))
	}
	return strings.Join(lines, "\n")
}
