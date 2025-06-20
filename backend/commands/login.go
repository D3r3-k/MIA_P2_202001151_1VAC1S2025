package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/global"
	Structs "MIA_PI_202001151_1VAC1S2025/structs"
	"MIA_PI_202001151_1VAC1S2025/utils"
	"encoding/binary"
	"fmt"
	"strings"
)

func Fn_Login(params string) {
	paramDefs := map[string]Structs.ParamDef{
		"-user": {Required: true},
		"-pass": {Required: true},
		"-id":   {Required: true},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return
	}
	user := parsed["-user"]
	pass := parsed["-pass"]
	id := parsed["-id"]
	part := utils.GetPartitionById(id)
	if part == nil {
		return
	}
	Login(user, pass, part)

}

// login -user=<username> -pass=<password> -id=<partition_id>
func Login(user string, pass string, part *Structs.Partition) {
	if globals.LoginSession.User != "" {
		utils.ShowMessage("Ya hay una sesión iniciada con el usuario: "+globals.LoginSession.User, true)
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

	// Leer superbloque
	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		utils.ShowMessage("No se pudo leer el superbloque.", true)
		return
	}
	// Leer inodo 1 (users.txt)
	inodoSize := int32(binary.Size(Structs.Inode{}))
	var inodeUser Structs.Inode
	if err := utils.ReadObject(file, &inodeUser, int64(sb.S_inode_start+inodoSize*1)); err != nil {
		utils.ShowMessage("No se pudo leer el inodo de users.txt.", true)
		return
	}
	// Leer TODOS los bloques directos usados
	blockSize := int32(binary.Size(Structs.Fileblock{}))
	var fullContent string
	for i := 0; i < 15; i++ {
		blockNum := inodeUser.I_block[i]
		if blockNum == -1 {
			continue
		}
		var blk Structs.Fileblock
		offset := int64(sb.S_block_start + blockSize*int32(blockNum))
		if err := utils.ReadObject(file, &blk, offset); err != nil {
			utils.ShowMessage("No se pudo leer el bloque "+fmt.Sprint(blockNum)+" de users.txt.", true)
			return
		}
		fullContent += string(blk.B_content[:])
	}
	// Buscar usuario y grupo
	var uid, gid int32
	var groupUser string
	found := false
	lines := strings.Split(fullContent, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) >= 5 && parts[0] != "0" && parts[1] == "U" && parts[3] == user && parts[4] == pass {
			found = true
			uid, err = utils.StringToInt32(strings.Trim(parts[0], "\x00"))
			if err != nil {
				utils.ShowMessage("Error al convertir UID: "+err.Error(), true)
				return
			}
			groupUser = strings.Trim(parts[2], "\x00")
			break
		}
	}
	if !found {
		utils.ShowMessage("Usuario o contraseña incorrecta.", true)
		return
	}
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) >= 3 && parts[0] != "0" && parts[1] == "G" && strings.Trim(parts[2], "\x00") == groupUser {
			gid, err = utils.StringToInt32(strings.Trim(parts[0], "\x00"))
			if err != nil {
				utils.ShowMessage("Error al convertir GID: "+err.Error(), true)
				return
			}
			break
		}
	}
	if gid == 0 {
		utils.ShowMessage("No se encontró el grupo del usuario: "+groupUser, true)
		return
	}

	// Guardar sesión
	globals.LoginSession = Structs.LoginSession{
		UID:         uid,
		GID:         gid,
		User:        user,
		Password:    pass,
		PartitionID: part.Id,
	}
	utils.ShowMessage("Sesión iniciada correctamente.", false)
}
