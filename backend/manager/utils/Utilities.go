package utils

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// METODOS PARA MENU
func ParseParameters(raw string, defs map[string]Structs.ParamDef) (map[string]string, error) {
	params := map[string]string{}
	for k, v := range defs {
		if !v.Required {
			params[k] = v.Default
		}
	}
	re := regexp.MustCompile(`-(\w+)(=("[^"]+"|\S+))?`)
	matches := re.FindAllStringSubmatch(raw, -1)
	found := map[string]bool{}
	for _, m := range matches {
		key := "-" + strings.ToLower(m[1])
		val := ""
		if len(m) > 3 && m[3] != "" {
			val = m[3]
			val = strings.Trim(val, "\"")
		}
		def, ok := defs[key]
		if !ok {
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
		found[key] = true

		if def.NotValue {
			if val != "" {
				return nil, fmt.Errorf("el parámetro [%s] no debe tener valor", key)
			}
			params[key] = "true"
		} else {
			if val == "" {
				return nil, fmt.Errorf("el parámetro [%s] requiere un valor", key)
			}
			params[key] = val
		}
	}
	for k, v := range defs {
		if v.Required && !found[k] {
			return nil, fmt.Errorf("falta el parámetro obligatorio: %s", k)
		}
	}
	return params, nil
}

func ParseCatParameters(raw string) ([]string, error) {
	matches := globals.Re.FindAllStringSubmatch(raw, -1)
	filesMap := map[int]string{}
	for _, m := range matches {
		key := strings.ToLower(m[1])
		val := strings.Trim(m[2], "\"")
		if strings.HasPrefix(key, "file") {
			numPart := key[4:]
			if num, err := strconv.Atoi(numPart); err == nil {
				filesMap[num] = val
			} else {
				return nil, fmt.Errorf("nombre de parámetro incorrecto: -%s", key)
			}
		} else {
			return nil, fmt.Errorf("parámetro desconocido: -%s", key)
		}
	}
	// Ordenar los archivos por número
	var indices []int
	for idx := range filesMap {
		indices = append(indices, idx)
	}
	sort.Ints(indices)
	var files []string
	for _, idx := range indices {
		files = append(files, filesMap[idx])
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("debe especificar al menos un parámetro -fileN")
	}
	return files, nil
}

// METODOS PARA COMANDOS
func CreateFile(name string) error {
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		ShowMessage("Error al crear el directorio: "+dir, true)
		return err
	}
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			ShowMessage("No se pudo crear el archivo: "+name, true)
			return err
		}
		defer file.Close()
	}
	return nil
}

func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		ShowMessage("No se pudo abrir el archivo: "+name, true)
		defer file.Close()
		return nil, err
	}
	return file, nil
}

func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		ShowMessage("No se pudo escribir el objeto: "+err.Error(), true)
		defer file.Close()
		return err
	}
	return nil
}

func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		ShowMessage("No se pudo leer el objeto del archivo: "+err.Error(), true)
		defer file.Close()
		return err
	}
	return nil
}

func GetRealSize(size int, unit string) int {
	switch unit {
	case "k":
		return size * 1024
	case "m":
		return size * 1024 * 1024
	default:
		return size
	}
}

func GetFreeSpaces(mbr *Structs.MBR, diskSize int32) [][2]int32 {
	type partition struct {
		Start, Size int32
	}
	var parts []partition
	for _, p := range mbr.Partitions {
		if p.Size > 0 {
			parts = append(parts, partition{p.Start, p.Size})
		}
	}
	sort.Slice(parts, func(i, j int) bool { return parts[i].Start < parts[j].Start })

	var freeSpaces [][2]int32
	mbrEnd := int32(binary.Size(*mbr))
	if len(parts) == 0 {
		freeSpaces = append(freeSpaces, [2]int32{mbrEnd, diskSize - mbrEnd})
	} else {
		if parts[0].Start > mbrEnd {
			freeSpaces = append(freeSpaces, [2]int32{mbrEnd, parts[0].Start - mbrEnd})
		}
		for i := 0; i < len(parts)-1; i++ {
			end := parts[i].Start + parts[i].Size
			if parts[i+1].Start > end {
				freeSpaces = append(freeSpaces, [2]int32{end, parts[i+1].Start - end})
			}
		}
		lastEnd := parts[len(parts)-1].Start + parts[len(parts)-1].Size
		if diskSize > lastEnd {
			freeSpaces = append(freeSpaces, [2]int32{lastEnd, diskSize - lastEnd})
		}
	}
	return freeSpaces
}

func GetPartitionById(id string) *Structs.Partition {
	drive := strings.ToUpper(string(id[0]))
	path := globals.PathDisks + drive + ".dsk"

	file, err := OpenFile(path)
	if err != nil {
		ShowMessage("No se pudo abrir el disco: "+path, true)
		return nil
	}
	defer file.Close()

	var mbr Structs.MBR
	if err := ReadObject(file, &mbr, 0); err != nil {
		ShowMessage("No se pudo leer el MBR.", true)
		return nil
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
		ShowMessage("Partición no montada o no existe con ID: "+id, true)
		return nil
	}
	return part
}

func GetFreeBlock(file *os.File, sb Structs.Superblock) int32 {
	bitmap := make([]byte, sb.S_blocks_count)
	_, err := file.ReadAt(bitmap, int64(sb.S_bm_block_start))
	if err != nil {
		return -1
	}
	for i, b := range bitmap {
		if b == '0' {
			return int32(i)
		}
	}
	return -1
}

func GetFreeInode(file *os.File, sb Structs.Superblock) int32 {
	bitmap := make([]byte, sb.S_inodes_count)
	_, err := file.ReadAt(bitmap, int64(sb.S_bm_inode_start))
	if err != nil {
		return -1
	}
	for i, b := range bitmap {
		if b == '0' {
			return int32(i)
		}
	}
	return -1
}

func CreateDirectory(file *os.File, sb *Structs.Superblock, parentInodeNum int32, nombre string, part Structs.Partition) (int32, error) {
	inodoSize := int32(binary.Size(Structs.Inode{}))
	blockSize := int32(binary.Size(Structs.Folderblock{}))
	newInode := GetFreeInode(file, *sb)
	if newInode == -1 {
		return -1, errors.New("no hay inodos libres")
	}
	newBlock := GetFreeBlock(file, *sb)
	if newBlock == -1 {
		return -1, errors.New("no hay bloques libres")
	}
	var fb Structs.Folderblock
	copy(fb.B_content[0].B_name[:], ".")
	fb.B_content[0].B_inodo = newInode
	copy(fb.B_content[1].B_name[:], "..")
	fb.B_content[1].B_inodo = parentInodeNum
	fb.B_content[2].B_inodo = -1
	fb.B_content[3].B_inodo = -1

	offsetBlock := int64(sb.S_block_start + newBlock*blockSize)
	WriteObject(file, fb, offsetBlock)

	var inode Structs.Inode
	inode.I_uid = globals.LoginSession.UID
	inode.I_gid = globals.LoginSession.GID
	inode.I_size = int32(binary.Size(Structs.Folderblock{}))
	copy(inode.I_atime[:], GetCurrentTimeString(16))
	copy(inode.I_ctime[:], GetCurrentTimeString(16))
	copy(inode.I_mtime[:], GetCurrentTimeString(16))
	inode.I_block[0] = newBlock
	for i := 1; i < 15; i++ {
		inode.I_block[i] = -1
	}
	inode.I_type[0] = '0'
	copy(inode.I_perm[:], "755")

	offsetInode := int64(sb.S_inode_start + newInode*inodoSize)
	WriteObject(file, inode, offsetInode)
	WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_inode_start+newInode))
	WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_block_start+newBlock))
	err := AddDirectoryEntry(file, sb, parentInodeNum, nombre, newInode)
	if err != nil {
		return -1, err
	}
	sb.S_free_inodes_count--
	sb.S_free_blocks_count--
	WriteObject(file, *sb, int64(part.Start))
	return newInode, nil
}

func SearchDirectoryEntry(file *os.File, sb *Structs.Superblock, dirInodeNum int32, nombre string) (bool, int32) {
	var inode Structs.Inode
	inodoSize := int32(binary.Size(Structs.Inode{}))
	blockSize := int32(binary.Size(Structs.Folderblock{}))
	if err := ReadObject(file, &inode, int64(sb.S_inode_start+dirInodeNum*inodoSize)); err != nil {
		return false, -1
	}
	for i := 0; i < 14; i++ {
		blockNum := inode.I_block[i]
		if blockNum == -1 {
			continue
		}
		var folderBlock Structs.Folderblock
		offset := int64(sb.S_block_start + blockNum*blockSize)
		if err := ReadObject(file, &folderBlock, offset); err != nil {
			continue
		}
		for _, entry := range folderBlock.B_content {
			n := strings.Trim(string(entry.B_name[:]), "\x00")
			if n == nombre && entry.B_inodo != -1 {
				return true, entry.B_inodo
			}
		}
	}
	return false, -1
}

func AddDirectoryEntry(file *os.File, sb *Structs.Superblock, parentInodeNum int32, nombre string, newInodeNum int32) error {
	inodoSize := int32(binary.Size(Structs.Inode{}))
	blockSize := int32(binary.Size(Structs.Folderblock{}))
	var parentInode Structs.Inode
	if err := ReadObject(file, &parentInode, int64(sb.S_inode_start+parentInodeNum*inodoSize)); err != nil {
		return errors.New("no se pudo leer el inodo del directorio padre")
	}
	added := false
	for i := 0; i < 14; i++ {
		bn := parentInode.I_block[i]
		if bn == -1 {
			freeBlock := GetFreeBlock(file, *sb)
			if freeBlock == -1 {
				return errors.New("no hay bloques libres para agregar entrada en directorio")
			}
			var fb Structs.Folderblock
			copy(fb.B_content[0].B_name[:], nombre)
			fb.B_content[0].B_inodo = newInodeNum
			for j := 1; j < 4; j++ {
				fb.B_content[j].B_inodo = -1
			}
			offset := int64(sb.S_block_start + freeBlock*blockSize)
			WriteObject(file, fb, offset)
			parentInode.I_block[i] = freeBlock
			WriteObject(file, parentInode, int64(sb.S_inode_start+parentInodeNum*inodoSize))
			WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_block_start+freeBlock))
			added = true
			break
		} else {
			var fb Structs.Folderblock
			offset := int64(sb.S_block_start + bn*blockSize)
			if err := ReadObject(file, &fb, offset); err != nil {
				continue
			}
			for j := 0; j < 4; j++ {
				if fb.B_content[j].B_inodo == -1 {
					copy(fb.B_content[j].B_name[:], nombre)
					fb.B_content[j].B_inodo = newInodeNum
					WriteObject(file, fb, offset)
					added = true
					break
				}
			}
			if added {
				break
			}
		}
	}

	if !added {
		if parentInode.I_block[14] == -1 {
			indirectBlock := GetFreeBlock(file, *sb)
			if indirectBlock == -1 {
				return errors.New("no hay bloques libres para el bloque indirecto")
			}
			parentInode.I_block[14] = indirectBlock
			sb.S_free_blocks_count--
			WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_block_start+indirectBlock))

			var pb Structs.Pointerblock
			for i := 0; i < 16; i++ {
				pb.B_pointers[i] = -1
			}
			offsetPB := int64(sb.S_block_start + blockSize*indirectBlock)
			WriteObject(file, pb, offsetPB)
		}

		var pb Structs.Pointerblock
		offsetPB := int64(sb.S_block_start + blockSize*parentInode.I_block[14])
		if err := ReadObject(file, &pb, offsetPB); err != nil {
			return errors.New("no se pudo leer el bloque indirecto")
		}

		for i := 0; i < 16; i++ {
			if pb.B_pointers[i] == -1 {
				newBlock := GetFreeBlock(file, *sb)
				if newBlock == -1 {
					return errors.New("no hay bloques libres para agregar entrada indirecta")
				}
				var fb Structs.Folderblock
				copy(fb.B_content[0].B_name[:], nombre)
				fb.B_content[0].B_inodo = newInodeNum
				for j := 1; j < 4; j++ {
					fb.B_content[j].B_inodo = -1
				}
				offset := int64(sb.S_block_start + blockSize*newBlock)
				WriteObject(file, fb, offset)
				pb.B_pointers[i] = newBlock
				WriteObject(file, pb, offsetPB)
				WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_block_start+newBlock))
				added = true
				break
			} else {
				var fb Structs.Folderblock
				offset := int64(sb.S_block_start + blockSize*pb.B_pointers[i])
				if err := ReadObject(file, &fb, offset); err != nil {
					continue
				}
				for j := 0; j < 4; j++ {
					if fb.B_content[j].B_inodo == -1 {
						copy(fb.B_content[j].B_name[:], nombre)
						fb.B_content[j].B_inodo = newInodeNum
						WriteObject(file, fb, offset)
						added = true
						break
					}
				}
				if added {
					break
				}
			}
		}
		WriteObject(file, parentInode, int64(sb.S_inode_start+parentInodeNum*inodoSize))
	}

	if !added {
		return errors.New("no hay espacio ni en bloques directos ni en el indirecto para una nueva entrada")
	}
	return nil
}

func FreeFileBlocks(file *os.File, sb *Structs.Superblock, inodeNum int32) {
	inodoSize := int32(binary.Size(Structs.Inode{}))
	var inode Structs.Inode
	if err := ReadObject(file, &inode, int64(sb.S_inode_start+inodeNum*inodoSize)); err != nil {
		return
	}
	for i := 0; i < 15; i++ {
		blockNum := inode.I_block[i]
		if blockNum != -1 {
			WriteObject(file, [1]byte{'0'}, int64(sb.S_bm_block_start+blockNum))
		}
	}
}

func AppendToFileBlocks(file *os.File, sb *Structs.Superblock, inode *Structs.Inode, content []byte) error {
	if inode.I_type[0] != '1' {
		return fmt.Errorf("AppendToFileBlocks solo debe usarse con archivos (I_type=1)")
	}
	blockSize := int32(binary.Size(Structs.Fileblock{}))
	contentLen := int32(len(content))

	// 1. Leer bloques usados directos (0–13)
	usedBlocks := []int32{}
	for i := 0; i < 14; i++ {
		if inode.I_block[i] != -1 {
			usedBlocks = append(usedBlocks, inode.I_block[i])
		}
	}
	espacioUsado := inode.I_size % 64
	espacioLibre := int32(64) - espacioUsado

	contentAgregado := int32(0)
	// 2. Completar en el último bloque directo si hay espacio
	if espacioLibre > 0 && len(usedBlocks) > 0 {
		lastBlockIdx := usedBlocks[len(usedBlocks)-1]
		var lastBlock Structs.Fileblock
		offset := int64(sb.S_block_start + blockSize*lastBlockIdx)
		if err := ReadObject(file, &lastBlock, offset); err != nil {
			return fmt.Errorf("no se pudo leer el último bloque: %v", err)
		}
		maxCopy := contentLen
		if maxCopy > espacioLibre {
			maxCopy = espacioLibre
		}
		copy(lastBlock.B_content[espacioUsado:], content[:maxCopy])
		if err := WriteObject(file, lastBlock, offset); err != nil {
			return fmt.Errorf("no se pudo escribir el último bloque: %v", err)
		}
		contentAgregado = maxCopy
	}

	// 3. Añadir bloques directos nuevos hasta 14 (0–13)
	for contentAgregado < contentLen && len(usedBlocks) < 14 {
		newBlock := GetFreeBlock(file, *sb)
		if newBlock == -1 {
			return fmt.Errorf("no hay bloques libres para el archivo")
		}
		WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_block_start+int32(newBlock)))
		// Asigna el bloque al primer slot libre
		asignado := false
		for i := 0; i < 14; i++ {
			if inode.I_block[i] == -1 {
				inode.I_block[i] = int32(newBlock)
				asignado = true
				break
			}
		}
		if !asignado {
			return fmt.Errorf("no hay apuntadores directos libres (esto no debería ocurrir aquí)")
		}
		sb.S_free_blocks_count--
		// Escribir el bloque de datos
		var blk Structs.Fileblock
		cuantos := contentLen - contentAgregado
		if cuantos > 64 {
			cuantos = 64
		}
		copy(blk.B_content[:], content[contentAgregado:contentAgregado+cuantos])
		offset := int64(sb.S_block_start + blockSize*int32(newBlock))
		if err := WriteObject(file, blk, offset); err != nil {
			return fmt.Errorf("no se pudo escribir el bloque nuevo: %v", err)
		}
		contentAgregado += cuantos
		usedBlocks = append(usedBlocks, int32(newBlock))
	}

	// 4. Si hace falta más, usa apuntador indirecto simple (I_block[14])
	if contentAgregado < contentLen {
		var pointerBlock Structs.Pointerblock
		var pointerBlockIdx int32
		if inode.I_block[14] == -1 {
			pointerBlockIdx = GetFreeBlock(file, *sb)
			if pointerBlockIdx == -1 {
				return fmt.Errorf("no hay bloques libres para el apuntador indirecto")
			}
			for i := 0; i < 16; i++ {
				pointerBlock.B_pointers[i] = -1
			}
			inode.I_block[14] = pointerBlockIdx
			sb.S_free_blocks_count--
			WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_block_start+pointerBlockIdx))
		} else {
			pointerBlockIdx = inode.I_block[14]
			// **CORREGIDO: Siempre usar blockSize aquí**
			offsetPointer := int64(sb.S_block_start + blockSize*pointerBlockIdx)
			if err := ReadObject(file, &pointerBlock, offsetPointer); err != nil {
				return fmt.Errorf("no se pudo leer el bloque de apuntadores: %v", err)
			}
		}

		// 4.2. Buscar slots libres en el pointerBlock y agregar bloques de datos ahí
		for contentAgregado < contentLen {
			slot := -1
			for i := 0; i < 16; i++ {
				if pointerBlock.B_pointers[i] == -1 {
					slot = i
					break
				}
			}
			if slot == -1 {
				return fmt.Errorf("no hay más punteros libres en el bloque indirecto simple (soporta hasta 16)")
			}
			newBlock := GetFreeBlock(file, *sb)
			if newBlock == -1 {
				return fmt.Errorf("no hay bloques libres para el archivo (indirecto)")
			}
			pointerBlock.B_pointers[slot] = int32(newBlock)
			sb.S_free_blocks_count--
			WriteObject(file, [1]byte{'1'}, int64(sb.S_bm_block_start+int32(newBlock)))
			// Escribir el bloque de datos
			var blk Structs.Fileblock
			cuantos := contentLen - contentAgregado
			if cuantos > 64 {
				cuantos = 64
			}
			copy(blk.B_content[:], content[contentAgregado:contentAgregado+cuantos])
			offset := int64(sb.S_block_start + blockSize*int32(newBlock))
			if err := WriteObject(file, blk, offset); err != nil {
				return fmt.Errorf("no se pudo escribir el bloque nuevo (indirecto): %v", err)
			}
			contentAgregado += cuantos
		}
		// **CORREGIDO: Siempre usar blockSize aquí**
		offsetPointer := int64(sb.S_block_start + blockSize*pointerBlockIdx)
		if err := WriteObject(file, pointerBlock, offsetPointer); err != nil {
			return fmt.Errorf("no se pudo escribir el bloque de apuntadores: %v", err)
		}
	}

	// 5. Actualizar el tamaño lógico del archivo (inodo)
	inode.I_size += contentLen
	return nil
}

/*
EXTRAS
*/
func ShowMessage(message string, isError bool) {
	if len(message) > 57 && !strings.Contains(message, "\n") {
		re := regexp.MustCompile(`.{1,57}(\s|$)`)
		message = re.ReplaceAllString(message, "$0\n")
	}
	arrMsg := strings.Split(message, "\n")
	if isError {
		fmt.Println("╔══════════════════════════[ERROR]══════════════════════════╗")
	} else {
		fmt.Println("╔══════════════════════════[INFO]═══════════════════════════╗")
	}
	for i, line := range arrMsg {
		if i == len(arrMsg)-1 && line == "" {
			continue
		}
		fmt.Printf("║ %-57s ║\n", line)
	}
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
}

func ShowMessageCustom(title string, msg string) {
	if len(msg) > 57 && !strings.Contains(msg, "\n") {
		re := regexp.MustCompile(`.{1,57}(\s|$)`)
		msg = re.ReplaceAllString(msg, "$0\n")
	}
	arrMsg := strings.Split(msg, "\n")
	if title != "" {
		padding := 57 - len(title)
		left := padding / 2
		right := padding - left
		fmt.Printf("╔%s[%s]%s╗\n", strings.Repeat("═", left), title, strings.Repeat("═", right))
	}
	for i, line := range arrMsg {
		if i == len(arrMsg)-1 && line == "" {
			continue
		}
		fmt.Printf("║ %-57s ║\n", line)
	}
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
}

func GetCurrentTimeString(x int) string {
	return time.Now().Format("2006-01-02 15:04:05")[:x]
}

func GenerateRandomSignature() int32 {
	signature := int32(0)
	for i := 0; i < 4; i++ {
		signature = (signature << 8) | int32('A'+i)
	}
	return signature
}

func ValidateRegex(value, pattern string) bool {
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

func PrepareDotFile(path string) (dotPath string, dotFormat string, fdot *os.File, cleanup func(), err error) {
	pathLower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(pathLower, ".png"):
		dotFormat = "png"
	case strings.HasSuffix(pathLower, ".jpg"):
		dotFormat = "jpg"
	case strings.HasSuffix(pathLower, ".pdf"):
		dotFormat = "pdf"
	case strings.HasSuffix(pathLower, ".svg"):
		dotFormat = "svg"
	default:
		err = fmt.Errorf("formato de salida no soportado (use .png, .jpg, .svg o .pdf)")
		return
	}
	dir := filepath.Dir(path)
	if _, er := os.Stat(dir); os.IsNotExist(er) {
		if er := os.MkdirAll(dir, 0755); er != nil {
			err = fmt.Errorf("no se pudo crear la carpeta de salida: %v", er)
			return
		}
	}
	baseName := filepath.Base(path)
	ext := filepath.Ext(baseName)
	baseNameNoExt := baseName[:len(baseName)-len(ext)]
	dotPath = filepath.Join(dir, baseNameNoExt+".dot")

	fdot, err = os.Create(dotPath)
	cleanup = func() { fdot.Close() }
	return
}

func GenerateGraphvizReport(dotPath, dotFormat, path string) error {
	cmd := exec.Command("dot", dotPath, "-T"+dotFormat, "-o", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("no se pudo ejecutar graphviz (dot): %v\nSalida: %s", err, string(output))
	}
	//os.Remove(dotPath)
	return nil
}

func OpenDisk(id string) (*os.File, string, error) {
	drive := strings.ToUpper(string(id[0]))
	diskPath := globals.PathDisks + drive + ".dsk"
	file, err := os.Open(diskPath)
	return file, diskPath, err
}
func GetDriveletter(id string) (string, error) {
	part := GetPartitionById(id)
	if part == nil {
		return "", fmt.Errorf("partición no encontrada con id [%s]", id)
	}
	return strings.ToUpper(string(part.Id[0])), nil
}

func InodePermString(perm []byte) string {
	if len(perm) < 3 {
		return "---"
	}
	result := ""
	for i := 0; i < 3; i++ {
		switch perm[i] {
		case '0':
			result += "---"
		case '1':
			result += "--x"
		case '2':
			result += "-w-"
		case '3':
			result += "-wx"
		case '4':
			result += "r--"
		case '5':
			result += "r-x"
		case '6':
			result += "rw-"
		case '7':
			result += "rwx"
		default:
			result += "---"
		}
	}
	return result
}

func GetUserAndGroupNames(partitionId string, uid int32, gid int32) (string, string) {
	part := GetPartitionById(partitionId)
	if part == nil {
		return "", ""
	}

	drive := strings.ToUpper(string(part.Id[0]))
	diskPath := globals.PathDisks + drive + ".dsk"
	file, err := OpenFile(diskPath)
	if err != nil {
		return "", ""
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := ReadObject(file, &sb, int64(part.Start)); err != nil {
		return "", ""
	}
	inodeSize := int32(binary.Size(Structs.Inode{}))
	var inodeUser Structs.Inode
	if err := ReadObject(file, &inodeUser, int64(sb.S_inode_start+inodeSize*1)); err != nil {
		return "", ""
	}
	blockSize := int32(binary.Size(Structs.Fileblock{}))
	var fullContent string
	for i := 0; i < 15; i++ {
		blockNum := inodeUser.I_block[i]
		if blockNum == -1 {
			continue
		}
		var blk Structs.Fileblock
		offset := int64(sb.S_block_start + blockSize*int32(blockNum))
		if err := ReadObject(file, &blk, offset); err != nil {
			continue
		}
		fullContent += string(blk.B_content[:])
	}
	var usuario, grupo string
	for _, line := range strings.Split(fullContent, "\n") {
		parts := strings.Split(line, ",")
		if len(parts) >= 5 && parts[1] == "U" && parts[0] != "0" {
			lineUID, err := strconv.Atoi(strings.Trim(parts[0], "\x00"))
			if err == nil && int32(lineUID) == uid {
				usuario = strings.TrimSpace(parts[3])
			}
		}
		if len(parts) >= 3 && parts[1] == "G" && parts[0] != "0" {
			lineGID, err := strconv.Atoi(strings.Trim(parts[0], "\x00"))
			if err == nil && int32(lineGID) == gid {
				grupo = strings.TrimSpace(parts[2])
			}
		}
		if usuario != "" && grupo != "" {
			break
		}
	}
	return usuario, grupo
}

func StringToInt32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("error al convertir a int32: %v", err)
	}
	if i < -2147483648 || i > 2147483647 {
		return 0, fmt.Errorf("valor fuera de rango para int32: %d", i)
	}
	return int32(i), nil
}
func CleanDOTString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\x00")
	var out strings.Builder
	for _, r := range s {
		if r >= 32 && r <= 126 {
			switch r {
			case '<':
				out.WriteString("&lt;")
			case '>':
				out.WriteString("&gt;")
			case '&':
				out.WriteString("&amp;")
			case '"':
				out.WriteString("&quot;")
			default:
				out.WriteRune(r)
			}
		}
	}
	return out.String()
}

func WriteJournaling(sb Structs.Superblock, part Structs.Partition, file *os.File, operation, path, content []byte) {
	if sb.S_filesystem_type == 3 {
		var journal Structs.Journaling
		journalingPos := int64(part.Start) + int64(binary.Size(sb))
		ReadObject(file, &journal, journalingPos)

		var entrada Structs.Content_J
		copy(entrada.Operation[:], operation)
		copy(entrada.Path[:], path)
		copy(entrada.Content[:], content)
		copy(entrada.Date[:], GetCurrentTimeString(16))

		if int(journal.Ultimo) < len(journal.Contenido) {
			journal.Contenido[journal.Ultimo] = entrada
			journal.Ultimo++
		} else {
			journal.Ultimo = 0
			journal.Contenido[journal.Ultimo] = entrada
			journal.Ultimo++
		}
		WriteObject(file, journal, journalingPos)
	}
}

func CleanByteArray(arr []byte) {
	for i := range arr {
		arr[i] = 0
	}
}
