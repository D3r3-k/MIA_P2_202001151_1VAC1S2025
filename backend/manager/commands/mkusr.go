package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Fn_Mkusr(params string) (string, error) {
	if globals.LoginSession.User == "" {
		utils.ShowMessage("Debe iniciar sesión primero.", true)
		return "", fmt.Errorf("debe iniciar sesión primero")
	}

	paramDefs := map[string]Structs.ParamDef{
		"-user": {Required: true},
		"-pass": {Required: true},
		"-grp":  {Required: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return "", err
	}

	user := parsed["-user"]
	pass := parsed["-pass"]
	grp := parsed["-grp"]

	if globals.LoginSession.User != "root" {
		utils.ShowMessage("Solo el usuario root puede crear usuarios.", true)
		return "", fmt.Errorf("solo el usuario root puede crear usuarios")
	}

	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		utils.ShowMessage("La partición de la sesión no está montada.", true)
		return "", fmt.Errorf("la partición de la sesión no está montada")
	}

	return Mkusr(user, pass, grp, part)
}

// mkusr -user=<usuario> -pass=<contraseña> -grp=<grupo>
func Mkusr(user string, pass string, grp string, part *Structs.Partition) (string, error) {
	drive := strings.ToUpper(string(part.Id[0]))
	path := globals.PathDisks + drive + ".dsk"

	file, err := utils.OpenFile(path)
	if err != nil {
		utils.ShowMessage("No se pudo abrir el disco: "+path, true)
		return "", err
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		utils.ShowMessage("No se pudo leer el superbloque.", true)
		return "", err
	}

	inodoSize := int32(binary.Size(Structs.Inode{}))
	var inodeUser Structs.Inode
	if err := utils.ReadObject(file, &inodeUser, int64(sb.S_inode_start+inodoSize*1)); err != nil {
		utils.ShowMessage("No se pudo leer el inodo de users.txt.", true)
		return "", err
	}

	// Leer contenido actual de users.txt
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
			return "", err
		}
		fullContent += string(blk.B_content[:])
	}

	// Verificar existencia del usuario y grupo
	lines := strings.Split(fullContent, "\n")
	groupExists := false
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) >= 3 && parts[1] == "U" && parts[3] == user && parts[0] != "0" {
			utils.ShowMessage("El usuario ["+user+"] ya existe.", true)
			return "", fmt.Errorf("el usuario [%s] ya existe", user)
		} else if len(parts) >= 3 && parts[1] == "G" && parts[2] == grp && parts[0] != "0" {
			groupExists = true
		}
	}
	if !groupExists {
		utils.ShowMessage("El grupo ["+grp+"] no existe. Debe crearse antes.", true)
		return "", fmt.Errorf("el grupo [%s] no existe", grp)
	}

	// Determinar siguiente ID de usuario
	nextUserID := 1
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) < 3 || parts[1] != "U" || parts[0] == "0" {
			continue
		}
		id, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err == nil && id >= nextUserID {
			nextUserID = id + 1
		}
	}

	newLine := fmt.Sprintf("%d,U,%s,%s,%s\n", nextUserID, grp, user, pass)

	fullContent = strings.TrimRight(fullContent, "\x00\n\r ") + "\n" + newLine
	fullContent = strings.TrimRight(fullContent, "\x00\n\r ") + "\n" + newLine

	utils.FreeFileBlocks(file, &sb, 1)
	inodeUser.I_size = 0
	for i := 0; i < 15; i++ {
		inodeUser.I_block[i] = -1
	}

	err = utils.AppendToFileBlocks(file, &sb, &inodeUser, []byte(fullContent))
	if err != nil {
		utils.ShowMessage("Error al reescribir users.txt: "+err.Error(), true)
		return "", err
	}

	utils.WriteJournaling(sb, *part, file, []byte("mkusr"), part.Name[:], []byte(newLine))

	utils.WriteObject(file, inodeUser, int64(sb.S_inode_start+inodoSize*1))
	utils.WriteObject(file, sb, int64(part.Start))

	utils.ShowMessage("Usuario ["+user+"] creado exitosamente", false)
	return "Usuario [" + user + "] creado exitosamente", nil
}
