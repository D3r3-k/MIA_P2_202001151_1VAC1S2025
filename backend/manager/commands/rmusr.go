package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"encoding/binary"
	"fmt"
	"strings"
)

func Fn_Rmusr(params string) {
	if globals.LoginSession.User == "" {
		utils.ShowMessage("Debe iniciar sesión primero.", true)
		return
	}
	paramDefs := map[string]Structs.ParamDef{
		"-user": {Required: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return
	}
	user := strings.TrimSpace(parsed["-user"])

	if globals.LoginSession.User != "root" {
		utils.ShowMessage("Solo el usuario root puede eliminar usuarios.", true)
		return
	}
	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		utils.ShowMessage("La partición de la sesión no está montada.", true)
		return
	}
	Rmusr(user, part)
}

// rmusr -user=<usuario>
func Rmusr(user string, part *Structs.Partition) {
	drive := strings.ToUpper(string(part.Id[0]))
	path := globals.PathDisks + drive + ".dsk"
	file, err := utils.OpenFile(path)
	if err != nil {
		utils.ShowMessage("No se pudo abrir el disco: "+path, true)
		return
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		utils.ShowMessage("No se pudo leer el superbloque.", true)
		return
	}
	inodoSize := int32(binary.Size(Structs.Inode{}))
	var inodeUser Structs.Inode
	if err := utils.ReadObject(file, &inodeUser, int64(sb.S_inode_start+inodoSize*1)); err != nil {
		utils.ShowMessage("No se pudo leer el inodo de users.txt.", true)
		return
	}
	blockSize := int32(binary.Size(Structs.Fileblock{}))

	var fullContent string
	var bloquesUsados []int32
	for i := 0; i < 15; i++ {
		blockNum := inodeUser.I_block[i]
		if blockNum == -1 {
			continue
		}
		var blk Structs.Fileblock
		offset := int64(sb.S_block_start + blockSize*blockNum)
		if err := utils.ReadObject(file, &blk, offset); err != nil {
			utils.ShowMessage("No se pudo leer el bloque "+fmt.Sprint(blockNum)+" de users.txt.", true)
			return
		}
		fullContent += string(blk.B_content[:])
		bloquesUsados = append(bloquesUsados, blockNum)
	}

	lines := strings.Split(fullContent, "\n")
	found := false
	for i, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) >= 5 && parts[1] == "U" && parts[3] == user && parts[0] != "0" {
			parts[0] = "0"
			lines[i] = strings.Join(parts, ",")
			found = true
			break
		}
	}
	if !found {
		utils.ShowMessage("No existe un usuario activo con el nombre '"+user+"'.", true)
		return
	}

	finalContent := strings.Join(lines, "\n")
	contentBytes := []byte(finalContent)
	numBloques := (len(contentBytes) + 63) / 64
	bloquesNecesarios := numBloques

	for i := 0; i < bloquesNecesarios; i++ {
		start := i * 64
		end := start + 64
		if end > len(contentBytes) {
			end = len(contentBytes)
		}
		var blk Structs.Fileblock
		copy(blk.B_content[:], contentBytes[start:end])
		offset := int64(sb.S_block_start + blockSize*int32(bloquesUsados[i]))
		utils.WriteObject(file, blk, offset)
	}

	inodeUser.I_size = int32(len(contentBytes))
	utils.WriteObject(file, inodeUser, int64(sb.S_inode_start+inodoSize*1))
	utils.WriteJournaling(sb, *part, file, []byte("rmgrp"), part.Name[:], []byte("remove_"+user))

	utils.ShowMessage("Usuario ["+user+"] eliminado exitosamente.", false)
}
