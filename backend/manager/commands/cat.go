package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"encoding/binary"
	"strings"
)

func Fn_Cat(params string) {
	if globals.LoginSession.User == "" {
		utils.ShowMessage("Debe iniciar sesión primero.", true)
		return
	}
	files, err := utils.ParseCatParameters(params)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return
	}
	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		utils.ShowMessage("La partición de la sesión no está montada.", true)
		return
	}
	cat(files, *part)
}

// cat <-file1,-file2,...>=<ruta1,ruta2,...>
func cat(files []string, part Structs.Partition) {
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

	for _, path := range files {
		if path == "" || path[0] != '/' {
			utils.ShowMessage("Ruta inválida: "+path, true)
			continue
		}
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			utils.ShowMessage("Ruta inválida.", true)
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
				utils.ShowMessage("Directorio no encontrado: "+dir, true)
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
			utils.ShowMessage("Archivo no encontrado: "+fileName, true)
			continue
		}
		inodoSize := int32(binary.Size(Structs.Inode{}))
		var inode Structs.Inode
		if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start+fileInodeNum*inodoSize)); err != nil {
			utils.ShowMessage("No se pudo leer el inodo del archivo.", true)
			continue
		}
		if inode.I_type[0] != '1' {
			utils.ShowMessage("El archivo no es un archivo regular.", true)
			continue
		}
		blockSize := int32(binary.Size(Structs.Fileblock{}))
		var totalContent []byte

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
		utils.ShowMessageCustom(path, string(totalContent[:limit]))
	}
}
