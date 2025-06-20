package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/global"
	Structs "MIA_PI_202001151_1VAC1S2025/structs"
	"MIA_PI_202001151_1VAC1S2025/utils"
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Fn_Mkfile(params string) {
	if globals.LoginSession.User == "" {
		utils.ShowMessage("Debe iniciar sesión primero.", true)
		return
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
		return
	}
	path := parsed["-path"]
	r := bool(parsed["-r"] == "true")
	sizeStr := parsed["-size"]
	cont := parsed["-cont"]
	if path[0] != '/' {
		utils.ShowMessage("La ruta debe comenzar desde la raiz [/].", true)
		return
	}
	part := utils.GetPartitionById(string(globals.LoginSession.PartitionID[:]))
	if part == nil {
		utils.ShowMessage("La partición de la sesión no está montada.", true)
		return
	}

	mkfile(path, sizeStr, cont, r, *part)
}

// mkfile -path=<ruta> [-r] -size=<tamaño> -cont=<contenido>
func mkfile(path string, sizeStr string, cont string, r bool, part Structs.Partition) {
	folders := strings.Split(path, "/")
	lastPart := folders[len(folders)-1]
	if lastPart == "" {
		utils.ShowMessage("La ruta termina en '/', se esperaba un archivo, no una carpeta.", true)
		return
	}
	if !strings.Contains(lastPart, ".") {
		utils.ShowMessage("Debe especificar un archivo con formato (ejemplo: nombre.txt).", true)
		return
	}
	fileName := lastPart
	if len(fileName) > 12 {
		utils.ShowMessage("El nombre del archivo (incluyendo su extensión) debe tener como máximo 12 caracteres: "+fileName, true)
		return
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_ .]+$`)
	if !re.MatchString(fileName) {
		utils.ShowMessage("Los nombres de los archivos deben tener solo caracteres alfanuméricos, guion bajo, espacio o punto: "+fileName, true)
		return
	}

	folders = folders[:len(folders)-1]

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
					return
				}
				currentInode = newInode
			} else {
				utils.ShowMessage("La carpeta ["+dir+"] no existe. Use -r para crearla.", true)
				return
			}
		}
	}

	// Buscar si existe el archivo y si lo vas a reemplazar
	exists, oldInodeNum := utils.SearchDirectoryEntry(file, &sb, currentInode, fileName)
	var inodeNum int32
	var inode Structs.Inode
	inodoSize := int32(binary.Size(Structs.Inode{}))

	if exists {
		utils.ShowMessage("El archivo ["+fileName+"] ya existe. Desea reemplazarlo?\nS: para reemplazar\nN: para cancelar", true)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">>> ")
		scanner.Scan()
		response := strings.TrimSpace(scanner.Text())
		if strings.ToUpper(response) != "S" {
			utils.ShowMessage("Operación cancelada por el usuario.", false)
			return
		}
		inodeNum = oldInodeNum
		utils.FreeFileBlocks(file, &sb, oldInodeNum)
		if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start+inodeNum*inodoSize)); err != nil {
			utils.ShowMessage("No se pudo leer el inodo viejo: "+err.Error(), true)
			return
		}
		for i := 0; i < 15; i++ {
			inode.I_block[i] = -1
		}
		inode.I_size = 0
	} else {
		tmp := utils.GetFreeInode(file, sb)
		if tmp == -1 {
			utils.ShowMessage("No hay inodos libres.", true)
			return
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
			return
		}
		content = data
	}
	if len(content) == 0 && sizeStr != "" {
		sz, err := strconv.Atoi(sizeStr)
		if err != nil || sz < 0 {
			utils.ShowMessage("El tamaño debe ser un número entero no negativo.", true)
			return
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
		return
	}

	offsetInode := int64(sb.S_inode_start + inodeNum*inodoSize)
	utils.WriteObject(file, inode, offsetInode)
	utils.WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_inode_start+inodeNum))

	if !exists {
		err = utils.AddDirectoryEntry(file, &sb, currentInode, fileName, inodeNum)
		if err != nil {
			utils.ShowMessage("Error al agregar entrada: "+err.Error(), true)
			return
		}
		sb.S_free_inodes_count--
	}
	utils.WriteObject(file, sb, int64(part.Start))
	utils.ShowMessage("Archivo ["+fileName+"] creado exitosamente.", false)
}
