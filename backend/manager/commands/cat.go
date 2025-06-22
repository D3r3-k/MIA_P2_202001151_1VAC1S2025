package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	utilsApi "MIA_PI_202001151_1VAC1S2025/utils"
	"encoding/binary"
	"fmt"
	"path/filepath"
	"strings"
)

type CatResponse struct {
	Name        string `json:"Name"`
	Size        string `json:"Size"`
	CreatedAt   string `json:"CreatedAt"`
	Owner       string `json:"Owner"`
	Content     string `json:"Content"`
	Extension   string `json:"Extension"`
	Permissions string `json:"Permissions"`
}

func Fn_Cat(params string) (*CatResponse, error) {
	if globals.LoginSession.User == "" {
		return nil, fmt.Errorf("debe iniciar sesión primero")
	}

	files, err := utils.ParseCatParameters(params)
	if err != nil {
		return nil, fmt.Errorf("error al analizar parámetros: %w", err)
	}

	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		return nil, fmt.Errorf("la partición de la sesión no está montada")
	}

	res, err := cat(files, *part)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func cat(files []string, part Structs.Partition) (*CatResponse, error) {
	drive := strings.ToUpper(string(part.Id[0]))
	diskPath := globals.PathDisks + drive + ".dsk"
	file, err := utils.OpenFile(diskPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el disco: %w", err)
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return nil, fmt.Errorf("no se pudo leer el superbloque: %w", err)
	}

	for _, path := range files {
		if path == "" || path[0] != '/' {
			continue
		}

		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			continue
		}

		fileName := parts[len(parts)-1]
		folders := parts[:len(parts)-1]
		currentInode := int32(0)

		found := true
		for _, dir := range folders {
			if dir == "" {
				continue
			}
			ok, childInode := utils.SearchDirectoryEntry(file, &sb, currentInode, dir)
			if !ok {
				found = false
				break
			}
			currentInode = childInode
		}
		if !found {
			continue
		}

		ok, fileInodeNum := utils.SearchDirectoryEntry(file, &sb, currentInode, fileName)
		if !ok {
			continue
		}

		var inode Structs.Inode
		inodoSize := int32(binary.Size(Structs.Inode{}))
		if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start+fileInodeNum*inodoSize)); err != nil {
			continue
		}

		if inode.I_type[0] != '1' {
			continue
		}

		var totalContent []byte
		blockSize := int32(binary.Size(Structs.Fileblock{}))

		for i := 0; i < 14; i++ {
			blockNum := inode.I_block[i]
			if blockNum == -1 {
				continue
			}
			var blk Structs.Fileblock
			offset := int64(sb.S_block_start + blockSize*blockNum)
			if err := utils.ReadObject(file, &blk, offset); err != nil {
				continue
			}
			totalContent = append(totalContent, blk.B_content[:]...)
		}

		if inode.I_block[14] != -1 {
			var pointerBlock Structs.Pointerblock
			offsetPointer := int64(sb.S_block_start + blockSize*inode.I_block[14])
			if err := utils.ReadObject(file, &pointerBlock, offsetPointer); err == nil {
				for i := 0; i < 16; i++ {
					ptr := pointerBlock.B_pointers[i]
					if ptr == -1 {
						continue
					}
					var blk Structs.Fileblock
					offset := int64(sb.S_block_start + blockSize*ptr)
					if err := utils.ReadObject(file, &blk, offset); err != nil {
						continue
					}
					totalContent = append(totalContent, blk.B_content[:]...)
				}
			}
		}

		limit := int(inode.I_size)
		if limit > len(totalContent) {
			limit = len(totalContent)
		}
		utils.ShowMessageCustom(path, string(totalContent))
		user, _ := utils.GetUserAndGroupNames(drive, inode.I_uid, inode.I_gid)
		_size := utilsApi.ConvertSizeToString(int32(inode.I_size))
		_perms := utils.InodePermString(inode.I_perm[:])
		_date := strings.Trim(string(inode.I_ctime[:]), "\x00")
		_content := strings.Trim(string(totalContent[:limit]), "\x00")
		response := &CatResponse{
			Name:        fileName,
			Size:        _size,
			CreatedAt:   _date,
			Owner:       user,
			Content:     _content,
			Extension:   filepath.Ext(fileName),
			Permissions: _perms,
		}
		return response, nil
	}
	return nil, fmt.Errorf("no se encontraron archivos para mostrar")
}
