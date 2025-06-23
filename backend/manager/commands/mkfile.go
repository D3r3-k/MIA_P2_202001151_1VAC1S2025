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

func Fn_Mkfile(params string) (string, error) {
	if globals.LoginSession.User == "" {
		utils.ShowMessage("Debe iniciar sesión primero.", true)
		return "", fmt.Errorf("debe iniciar sesión primero")
	}
	paramDefs := map[string]Structs.ParamDef{
		"-path": {Required: true},
		"-r":    {Required: false, NotValue: true},
		"-size": {Required: false},
		"-cont": {Required: false},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return "", err
	}
	path := parsed["-path"]
	r := bool(parsed["-r"] == "true")
	sizeStr := parsed["-size"]
	cont := parsed["-cont"]
	if path[0] != '/' {
		utils.ShowMessage("La ruta debe comenzar desde la raiz [/].", true)
		return "", fmt.Errorf("la ruta debe comenzar desde la raiz [/]")
	}
	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		utils.ShowMessage("La partición de la sesión no está montada.", true)
		return "", fmt.Errorf("la partición de la sesión no está montada")
	}

	return mkfile(path, sizeStr, cont, r, *part)
}

// mkfile -path=<ruta> [-r] -size=<tamaño> -cont=<contenido>
func mkfile(path string, sizeStr string, cont string, r bool, part Structs.Partition) (string, error) {
	folders := strings.Split(path, "/")
	lastPart := folders[len(folders)-1]
	if lastPart == "" {
		utils.ShowMessage("La ruta termina en '/', se esperaba un archivo, no una carpeta.", true)
		return "", fmt.Errorf("la ruta termina en '/', se esperaba un archivo, no una carpeta")
	}
	if !strings.Contains(lastPart, ".") {
		utils.ShowMessage("Debe especificar un archivo con formato (ejemplo: nombre.txt).", true)
		return "", fmt.Errorf("debe especificar un archivo con formato (ejemplo: nombre.txt)")
	}
	fileName := lastPart
	if len(fileName) > 12 {
		utils.ShowMessage("El nombre del archivo (incluyendo su extensión) debe tener como máximo 12 caracteres: "+fileName, true)
		return "", fmt.Errorf("el nombre del archivo (incluyendo su extensión) debe tener como máximo 12 caracteres: %s", fileName)
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_ .]+$`)
	if !re.MatchString(fileName) {
		utils.ShowMessage("Los nombres de los archivos deben tener solo caracteres alfanuméricos, guion bajo, espacio o punto: "+fileName, true)
		return "", fmt.Errorf("los nombres de los archivos deben tener solo caracteres alfanuméricos, guion bajo, espacio o punto: %s", fileName)
	}

	folders = folders[:len(folders)-1]

	drive := strings.ToUpper(string(part.Id[0]))
	diskPath := globals.PathDisks + drive + ".dsk"
	file, err := utils.OpenFile(diskPath)
	if err != nil {
		utils.ShowMessage("No se pudo abrir el disco: "+diskPath, true)
		return "", fmt.Errorf("no se pudo abrir el disco: %s", diskPath)
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		utils.ShowMessage("No se pudo leer el superbloque.", true)
		return "", fmt.Errorf("no se pudo leer el superbloque: %v", err)
	}

	currentInode := int32(0)
	for _, dir := range folders {
		if dir == "" {
			continue
		}
		exists, childInode := utils.SearchDirectoryEntry(file, &sb, currentInode, dir)
		if exists {
			currentInode = childInode
		} else {
			if r {
				newInode, err := utils.CreateDirectory(file, &sb, currentInode, dir, part)
				if err != nil {
					utils.ShowMessage("Error al crear carpeta ["+dir+"]: "+err.Error(), true)
					return "", err
				}
				currentInode = newInode
			} else {
				utils.ShowMessage("La carpeta ["+dir+"] no existe. Use -r para crearla.", true)
				return "", fmt.Errorf("la carpeta [%s] no existe, use -r para crearla", dir)
			}
		}
	}

	// Buscar si existe el archivo y si lo vas a reemplazar
	exists, oldInodeNum := utils.SearchDirectoryEntry(file, &sb, currentInode, fileName)
	var inodeNum int32
	var inode Structs.Inode
	inodoSize := int32(binary.Size(Structs.Inode{}))

	if exists {
		inodeNum = oldInodeNum
		utils.FreeFileBlocks(file, &sb, oldInodeNum)
		if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start+inodeNum*inodoSize)); err != nil {
			utils.ShowMessage("No se pudo leer el inodo viejo: "+err.Error(), true)
			return "", fmt.Errorf("no se pudo leer el inodo viejo: %v", err)
		}
		for i := 0; i < 15; i++ {
			inode.I_block[i] = -1
		}
		inode.I_size = 0
	} else {
		tmp := utils.GetFreeInode(file, sb)
		if tmp == -1 {
			utils.ShowMessage("No hay inodos libres.", true)
			return "", fmt.Errorf("no hay inodos libres")
		}
		inodeNum = tmp
		inode.I_uid = globals.LoginSession.UID
		inode.I_gid = globals.LoginSession.GID
		inode.I_size = 0
		copy(inode.I_atime[:], utils.GetCurrentTimeString(16))
		copy(inode.I_ctime[:], utils.GetCurrentTimeString(16))
		copy(inode.I_mtime[:], utils.GetCurrentTimeString(16))
		inode.I_type[0] = '1'
		copy(inode.I_perm[:], "664")
		for i := 0; i < 15; i++ {
			inode.I_block[i] = -1
		}
	}

	var content []byte
	if cont != "" {
		data, err := os.ReadFile(cont)
		if err != nil {
			utils.ShowMessage("No se pudo leer el archivo de contenido: "+err.Error(), true)
			return "", err
		}
		content = data
	}
	if len(content) == 0 && sizeStr != "" {
		sz, err := strconv.Atoi(sizeStr)
		if err != nil || sz < 0 {
			utils.ShowMessage("El tamaño debe ser un número entero no negativo.", true)
			return "", fmt.Errorf("el tamaño debe ser un número entero no negativo")
		}
		patron := []byte("0123456789")
		for len(content) < sz {
			resto := sz - len(content)
			if resto >= 10 {
				content = append(content, patron...)
			} else {
				content = append(content, patron[:resto]...)
			}
		}
	}

	utils.WriteJournaling(sb, part, file, []byte("mkfile"), []byte(path), content)

	if err := utils.AppendToFileBlocks(file, &sb, &inode, content); err != nil {
		utils.ShowMessage("Error al escribir bloques del archivo: "+err.Error(), true)
		return "", err
	}

	offsetInode := int64(sb.S_inode_start + inodeNum*inodoSize)
	utils.WriteObject(file, inode, offsetInode)
	utils.WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_inode_start+inodeNum))

	if !exists {
		err = utils.AddDirectoryEntry(file, &sb, currentInode, fileName, inodeNum)
		if err != nil {
			utils.ShowMessage("Error al agregar entrada: "+err.Error(), true)
			return "", err
		}
		sb.S_free_inodes_count--
	}
	utils.WriteObject(file, sb, int64(part.Start))
	utils.ShowMessage("Archivo ["+fileName+"] creado exitosamente.", false)
	return "Archivo [" + fileName + "] creado exitosamente.", nil
}
