package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/global"
	Structs "MIA_PI_202001151_1VAC1S2025/structs"
	"MIA_PI_202001151_1VAC1S2025/utils"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Fn_Mkgrp(params string) {
	if globals.LoginSession.User == "" {
		utils.ShowMessage("Debe iniciar sesión primero.", true)
		return
	}
	paramDefs := map[string]Structs.ParamDef{
		"-name": {Required: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return
	}
	name := parsed["-name"]
	if len(name) > 10 {
		utils.ShowMessage("El nombre del grupo no puede exceder los 10 caracteres.", true)
		return
	}
	if globals.LoginSession.User != "root" {
		utils.ShowMessage("Solo el usuario root puede crear grupos.", true)
		return
	}
	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		utils.ShowMessage("La partición de la sesión no está montada.", true)
		return
	}
	Mkgrp(name, part)
}

func Mkgrp(name string, part *Structs.Partition) {
	part = utils.GetPartitionById(string(part.Id[:]))
	if part == nil {
		utils.ShowMessage("La partición no está montada.", true)
		return
	}
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

	// Leer contenido actual
	var fullContent string
	for i := 0; i < 15; i++ {
		blockNum := inodeUser.I_block[i]
		if blockNum == -1 {
			continue
		}
		var blk Structs.Fileblock
		blockSize := int32(binary.Size(Structs.Fileblock{}))
		offset := int64(sb.S_block_start + blockSize*blockNum)
		if err := utils.ReadObject(file, &blk, offset); err != nil {
			utils.ShowMessage("No se pudo leer el bloque "+fmt.Sprint(blockNum)+" de users.txt.", true)
			return
		}
		fullContent += string(blk.B_content[:])
	}

	// Verificar existencia
	lines := strings.Split(fullContent, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) >= 3 && parts[1] == "G" && parts[2] == name && parts[0] != "0" {
			utils.ShowMessage("El grupo '"+name+"' ya existe.", true)
			return
		}
	}

	// Calcular siguiente ID
	nextGroupID := 1
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) < 3 || parts[1] != "G" || parts[0] == "0" {
			continue
		}
		id, err := strconv.Atoi(strings.Trim(parts[0], "\x00"))
		if err == nil && id >= nextGroupID {
			nextGroupID = id + 1
		}
	}

	// Agregar nuevo grupo usando AppendToFileBlocks
	newLine := fmt.Sprintf("%d,G,%s\n", nextGroupID, name)
	err = utils.AppendToFileBlocks(file, &sb, &inodeUser, []byte(newLine))
	if err != nil {
		utils.ShowMessage("Error al agregar grupo: "+err.Error(), true)
		return
	}

	// Actualizar inodo users.txt
	utils.WriteObject(file, inodeUser, int64(sb.S_inode_start+inodoSize*1))
	utils.ShowMessage("Grupo '"+name+"' creado exitosamente.", false)
}
