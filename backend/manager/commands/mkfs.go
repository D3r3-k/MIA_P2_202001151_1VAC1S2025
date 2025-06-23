package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Fn_Mkfs(params string) (string, error) {
	paramDefs := map[string]Structs.ParamDef{
		"-id":   {Required: true},
		"-type": {Required: false, Default: "full"},
		"-fs":   {Required: false, Default: "2fs"},
	}

	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return "", err
	}

	id := parsed["-id"]
	fsType := strings.ToLower(parsed["-type"])
	fs := strings.ToLower(parsed["-fs"])

	if fsType != "full" {
		utils.ShowMessage("Solo se permite formateo <full> actualmente.", true)
		return "", fmt.Errorf("tipo de formateo no soportado: %s", fsType)
	}
	if fs != "2fs" && fs != "3fs" {
		utils.ShowMessage("Solo se permite sistema de archivos <2fs> o <3fs>.", true)
		return "", fmt.Errorf("solo se permite sistema de archivos <2fs> o <3fs>")
	}

	return Mkfs(id, fs)
}

// Mkfs -id=<id> [-type=full] [-fs=2fs|3fs]
func Mkfs(id string, fs string) (string, error) {
	if globals.LoginSession.User != "" {
		utils.ShowMessage("Debe cerrar sesión antes de formatear un disco.", true)
		return "", fmt.Errorf("sesión activa, cierre sesión antes de formatear")
	}
	drive := strings.ToUpper(string(id[0]))
	path := globals.PathDisks + drive + ".dsk"

	file, err := utils.OpenFile(path)
	if err != nil {
		utils.ShowMessage("No se pudo abrir el disco: "+path, true)
		return "", fmt.Errorf("no se pudo abrir el disco: %s", path)
	}
	defer file.Close()

	var mbr Structs.MBR
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		utils.ShowMessage("No se pudo leer el MBR.", true)
		return "", fmt.Errorf("no se pudo leer el MBR: %v", err)
	}

	var part *Structs.Partition
	for i := 0; i < 4; i++ {
		p := &mbr.Partitions[i]
		if strings.Contains(string(p.Id[:]), id) && strings.TrimSpace(string(p.Status[:])) == "1" {
			part = p
			break
		}
	}
	if part == nil {
		utils.ShowMessage("mkfs Partición no montada o no existe con ID: "+id, true)
		return "", fmt.Errorf("partición no montada o no existe con ID: %s", id)
	}

	partSize := part.Size
	SuperBlockSize := int32(binary.Size(Structs.Superblock{}))
	InodoSize := int32(binary.Size(Structs.Inode{}))
	blockSize := int32(binary.Size(Structs.Folderblock{}))
	JournSize := int32(0)

	n := (partSize - SuperBlockSize) / (4 + InodoSize + (3 * blockSize))
	if fs == "3fs" {
		JournSize = int32(binary.Size(Structs.Journaling{}))
	}
	n = int32(n)
	if fs == "2fs" {
		return mkfs_ext2(n, *part, file)
	} else {
		n = (partSize - SuperBlockSize - JournSize) / (4 + InodoSize + (3 * blockSize))
		return mkfs_ext3(n, *part, file)
	}
}

// mkfs -type=full -fs=2fs -id=<id>
func mkfs_ext2(n int32, part Structs.Partition, file *os.File) (string, error) {
	// 0. calcular tamaños
	inodoSize := int32(binary.Size(Structs.Inode{}))
	blockSize := int32(binary.Size(Structs.Folderblock{}))
	// 1. Inicializar Superblock
	var sb Structs.Superblock
	sb.S_filesystem_type = 2
	sb.S_inodes_count = n
	sb.S_blocks_count = 3 * n
	sb.S_free_blocks_count = 3*n - 2
	sb.S_free_inodes_count = n - 2
	copy(sb.S_mtime[:], utils.GetCurrentTimeString(16))
	copy(sb.S_umtime[:], utils.GetCurrentTimeString(16))
	sb.S_mnt_count = 1
	sb.S_magic = 0xEF53
	sb.S_inode_size = inodoSize
	sb.S_block_size = blockSize
	sb.S_fist_ino = 2
	sb.S_first_blo = 2
	sb.S_bm_inode_start = part.Start + int32(binary.Size(Structs.Superblock{}))
	sb.S_bm_block_start = sb.S_bm_inode_start + n
	sb.S_inode_start = sb.S_bm_block_start + 3*n
	sb.S_block_start = sb.S_inode_start + n*inodoSize
	// 2. Escribir Superblock
	utils.WriteObject(file, sb, int64(part.Start))
	// 3. Inicializar bitmaps
	for i := int32(0); i < n; i++ {
		utils.WriteObject(file, [1]byte{'0'}, int64(sb.S_bm_inode_start+i))
	}
	for i := int32(0); i < 3*n; i++ {
		utils.WriteObject(file, [1]byte{'0'}, int64(sb.S_bm_block_start+i))
	}
	// 4. Inicializar inodos vacíos
	inodeEmpty := Structs.Inode{}
	for i := 0; i < 15; i++ {
		inodeEmpty.I_block[i] = -1
	}
	for i := int32(0); i < n; i++ {
		utils.WriteObject(file, inodeEmpty, int64(sb.S_inode_start+i*inodoSize))
	}
	// 5. Inicializar bloques vacíos
	blkEmpty := Structs.Folderblock{}
	for i := int32(0); i < 3*n; i++ {
		utils.WriteObject(file, blkEmpty, int64(sb.S_block_start+i*blockSize))
	}
	// 6. Crear inodo 0 (raíz)
	var inodo0 Structs.Inode
	inodo0.I_uid = 1
	inodo0.I_gid = 1
	inodo0.I_size = 0
	inodo0.I_type[0] = '0'
	copy(inodo0.I_perm[:], "664")
	copy(inodo0.I_ctime[:], utils.GetCurrentTimeString(16))
	copy(inodo0.I_mtime[:], utils.GetCurrentTimeString(16))
	copy(inodo0.I_atime[:], utils.GetCurrentTimeString(16))
	for i := 0; i < 15; i++ {
		inodo0.I_block[i] = -1
	}
	inodo0.I_block[0] = 0
	// 7. Crear inodo 1 (users.txt)
	data := "1,G,root\n1,U,root,root,123\n"
	var inodo1 Structs.Inode
	inodo1.I_uid = 1
	inodo1.I_gid = 1
	inodo1.I_size = int32(len(data))
	inodo1.I_type[0] = '1'
	copy(inodo1.I_perm[:], "664")
	copy(inodo1.I_ctime[:], utils.GetCurrentTimeString(16))
	copy(inodo1.I_mtime[:], utils.GetCurrentTimeString(16))
	copy(inodo1.I_atime[:], utils.GetCurrentTimeString(16))
	for i := 0; i < 15; i++ {
		inodo1.I_block[i] = -1
	}
	inodo1.I_block[0] = 1
	// 8. Bloque 0 - carpeta raíz
	var folder Structs.Folderblock
	copy(folder.B_content[0].B_name[:], ".")
	folder.B_content[0].B_inodo = 0
	copy(folder.B_content[1].B_name[:], "..")
	folder.B_content[1].B_inodo = 0
	copy(folder.B_content[2].B_name[:], "users.txt")
	folder.B_content[2].B_inodo = 1
	folder.B_content[3].B_inodo = -1
	// 9. Bloque 1 - users.txt
	var blkUsers Structs.Fileblock
	copy(blkUsers.B_content[:], data)
	// 10. Bitmaps
	utils.WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_inode_start+0)) // Inodo 0
	utils.WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_inode_start+1)) // Inodo 1
	utils.WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_block_start+0)) // Bloque 0
	utils.WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_block_start+1)) // Bloque 1

	// 11. Escribir estructuras
	utils.WriteObject(file, inodo0, int64(sb.S_inode_start+0*inodoSize))
	utils.WriteObject(file, inodo1, int64(sb.S_inode_start+1*inodoSize))
	utils.WriteObject(file, folder, int64(sb.S_block_start+0*blockSize))
	utils.WriteObject(file, blkUsers, int64(sb.S_block_start+1*blockSize))
	utils.WriteObject(file, sb, int64(part.Start))

	// Terminado
	utils.ShowMessage("Sistema de archivos EXT2 creado exitosamente en la partición: "+string(part.Id[:]), false)
	return "Partición " + string(part.Id[:]) + " formateada exitosamente con sistema de archivos EXT2", nil
}

// mkfs -type=full -fs=3fs -id=<id>
func mkfs_ext3(n int32, part Structs.Partition, file *os.File) (string, error) {
	// 0. calcular tamaños
	inodoSize := int32(binary.Size(Structs.Inode{}))
	blockSize := int32(binary.Size(Structs.Folderblock{}))
	journSize := int32(binary.Size(Structs.Journaling{}))

	// 1. Inicializar Superblock
	var sb Structs.Superblock
	sb.S_filesystem_type = 3
	sb.S_inodes_count = n
	sb.S_blocks_count = 3 * n
	sb.S_free_blocks_count = 3*n - 2
	sb.S_free_inodes_count = n - 2
	copy(sb.S_mtime[:], utils.GetCurrentTimeString(16))
	copy(sb.S_umtime[:], utils.GetCurrentTimeString(16))
	sb.S_mnt_count = 1
	sb.S_magic = 0xEF53
	sb.S_inode_size = inodoSize
	sb.S_block_size = blockSize
	sb.S_fist_ino = 2
	sb.S_first_blo = 2
	sb.S_bm_inode_start = part.Start + int32(binary.Size(Structs.Superblock{})) + journSize
	sb.S_bm_block_start = sb.S_bm_inode_start + n
	sb.S_inode_start = sb.S_bm_block_start + 3*n
	sb.S_block_start = sb.S_inode_start + n*inodoSize

	// 2. Escribir Superblock
	utils.WriteObject(file, sb, int64(binary.Size(part.Start)+int(journSize)))

	// 3. Inicializar Journaling (solo EXT3)
	var journal Structs.Journaling
	journal.Size = 50
	journal.Ultimo = 0
	utils.WriteObject(file, journal, int64(part.Start+int32(binary.Size(Structs.Superblock{}))))

	// 4. Inicializar bitmaps
	for i := int32(0); i < n; i++ {
		utils.WriteObject(file, [1]byte{'0'}, int64(sb.S_bm_inode_start+i))
	}
	for i := int32(0); i < 3*n; i++ {
		utils.WriteObject(file, [1]byte{'0'}, int64(sb.S_bm_block_start+i))
	}

	// 5. Inicializar inodos vacíos
	inodeEmpty := Structs.Inode{}
	for i := 0; i < 15; i++ {
		inodeEmpty.I_block[i] = -1
	}
	for i := int32(0); i < n; i++ {
		utils.WriteObject(file, inodeEmpty, int64(sb.S_inode_start+i*inodoSize))
	}

	// 6. Inicializar bloques vacíos
	blkEmpty := Structs.Folderblock{}
	for i := int32(0); i < 3*n; i++ {
		utils.WriteObject(file, blkEmpty, int64(sb.S_block_start+i*blockSize))
	}

	// 7. Crear inodo 0 (raíz)
	var inodo0 Structs.Inode
	inodo0.I_uid = 1
	inodo0.I_gid = 1
	inodo0.I_size = 0
	inodo0.I_type[0] = '0'
	copy(inodo0.I_perm[:], "664")
	copy(inodo0.I_ctime[:], utils.GetCurrentTimeString(16))
	copy(inodo0.I_mtime[:], utils.GetCurrentTimeString(16))
	copy(inodo0.I_atime[:], utils.GetCurrentTimeString(16))
	for i := 0; i < 15; i++ {
		inodo0.I_block[i] = -1
	}
	inodo0.I_block[0] = 0

	// 8. Crear inodo 1 (users.txt)
	data := "1,G,root\n1,U,root,root,123\n"
	var inodo1 Structs.Inode
	inodo1.I_uid = 1
	inodo1.I_gid = 1
	inodo1.I_size = int32(len(data))
	inodo1.I_type[0] = '1'
	copy(inodo1.I_perm[:], "664")
	copy(inodo1.I_ctime[:], utils.GetCurrentTimeString(16))
	copy(inodo1.I_mtime[:], utils.GetCurrentTimeString(16))
	copy(inodo1.I_atime[:], utils.GetCurrentTimeString(16))
	for i := 0; i < 15; i++ {
		inodo1.I_block[i] = -1
	}
	inodo1.I_block[0] = 1

	// 9. Bloque 0 - carpeta raíz
	var folder Structs.Folderblock
	copy(folder.B_content[0].B_name[:], ".")
	folder.B_content[0].B_inodo = 0
	copy(folder.B_content[1].B_name[:], "..")
	folder.B_content[1].B_inodo = 0
	copy(folder.B_content[2].B_name[:], "users.txt")
	folder.B_content[2].B_inodo = 1
	folder.B_content[3].B_inodo = -1

	// 10. Bloque 1 - users.txt
	var blkUsers Structs.Fileblock
	copy(blkUsers.B_content[:], data)

	// 11. Bitmaps (asignados)
	utils.WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_inode_start+0)) // Inodo 0
	utils.WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_inode_start+1)) // Inodo 1
	utils.WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_block_start+0)) // Bloque 0
	utils.WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_block_start+1)) // Bloque 1

	// 12. Escribir estructuras
	utils.WriteObject(file, inodo0, int64(sb.S_inode_start+0*inodoSize))
	utils.WriteObject(file, inodo1, int64(sb.S_inode_start+1*inodoSize))
	utils.WriteObject(file, folder, int64(sb.S_block_start+0*blockSize))
	utils.WriteObject(file, blkUsers, int64(sb.S_block_start+1*blockSize))
	utils.WriteObject(file, sb, int64(part.Start))

	utils.ShowMessage("Sistema de archivos EXT3 creado exitosamente en la partición: "+string(part.Id[:]), false)
	return "Partición " + string(part.Id[:]) + " formateada exitosamente con sistema de archivos EXT3", nil
}
