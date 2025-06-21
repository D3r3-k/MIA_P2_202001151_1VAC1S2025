package Structs

import (
	"fmt"
	"strings"
)

type MBR struct {
	MbrSize      int32
	CreationDate [10]byte
	Signature    int32
	Fit          [1]byte
	Partitions   [4]Partition
}

func PrintMBR(data MBR, drive string) {
	fmt.Println("╔═════════════════════════════════╗═════╗")
	fmt.Println("║      MBR - MASTER BOOT RECORD   ║  " + drive + "  ║")
	fmt.Println("╠═════════════════════════════════╣═════╝")
	fmt.Printf("║ Creacion     : %-16s ║\n", string(data.CreationDate[:]))
	fmt.Printf("║ Fit          : %-16s ║\n", string(data.Fit[:]))
	fmt.Printf("║ Tamaño       : %-16d ║\n", data.MbrSize)
	fmt.Printf("║ Signature    : %-16d ║\n", data.Signature)
	fmt.Println("╠═════════════════════════════════╣")
	fmt.Println("║ Particiones:                    ║")
	fmt.Println("╠════╦══════════════╦══════╦══════╬═══════════╦════════════╦═════════╦════════╗")
	fmt.Println("║ ID ║   Nombre     ║ Tipo ║ Fit  ║  Inicio   ║  Tamaño    ║ Estado  ║   ID   ║")
	fmt.Println("╠════╬══════════════╬══════╬══════╬═══════════╬════════════╬═════════╬════════╣")
	for i := 0; i < 4; i++ {
		PrintPartitionTable(i, data.Partitions[i])
	}
	fmt.Println("╚════╩══════════════╩══════╩══════╩═══════════╩════════════╩═════════╩════════╝")
}

type Partition struct {
	Status      [1]byte
	Type        [1]byte
	Fit         [1]byte
	Start       int32
	Size        int32
	Name        [16]byte
	Correlative int32
	Id          [4]byte
}

func PrintPartitions(data MBR, drive string) {
	fmt.Println("╔═════════════════════════════════╗═════╗")
	fmt.Println("║ Particiones:                    ║  " + drive + "  ║")
	fmt.Println("╠════╦══════════════╦══════╦══════╬═══════════╦════════════╦═════════╦════════╗")
	fmt.Println("║ ID ║   Nombre     ║ Tipo ║ Fit  ║  Inicio   ║  Tamaño    ║ Estado  ║   ID   ║")
	fmt.Println("╠════╬══════════════╬══════╬══════╬═══════════╬════════════╬═════════╬════════╣")
	for i := 0; i < 4; i++ {
		PrintPartitionTable(i, data.Partitions[i])
	}
	fmt.Println("╚════╩══════════════╩══════╩══════╩═══════════╩════════════╩═════════╩════════╝")
}

func PrintPartition(data Partition, drive string) {
	fmt.Println("╔═════════════════════════════════╗═════╗")
	fmt.Println("║ Partición:                      ║  " + drive + "  ║")
	fmt.Println("╠═════════════════════════════════╣═════╝")
	fmt.Printf("║ Nombre      : %-17s ║\n", strings.TrimRight(string(data.Name[:]), "\x00"))
	fmt.Printf("║ Tipo        : %-17s ║\n", strings.TrimRight(string(data.Type[:]), "\x00"))
	fmt.Printf("║ Fit         : %-17s ║\n", strings.TrimRight(string(data.Fit[:]), "\x00"))
	fmt.Printf("║ Inicio      : %-17d ║\n", data.Start)
	fmt.Printf("║ Tamaño      : %-17d ║\n", data.Size)
	estado := "Desmontado"
	if strings.TrimRight(string(data.Status[:]), "\x00") == "1" {
		estado = "Montado"
	}
	fmt.Printf("║ Estado      : %-17s ║\n", estado)
	fmt.Printf("║ ID          : %-17s ║\n", strings.TrimRight(string(data.Id[:]), "\x00"))
	fmt.Println("╠═════════════════════════════════╣")
	fmt.Printf("║ Correlativo : %-17d ║\n", data.Correlative)
	fmt.Println("╚═════════════════════════════════╝")
}

func PrintPartitionTable(idx int, data Partition) {
	name := strings.TrimRight(string(data.Name[:]), "\x00")
	name = strings.TrimSpace(name)
	typeChar := strings.TrimRight(string(data.Type[:]), "\x00")
	typeChar = strings.TrimSpace(typeChar)
	fit := strings.TrimRight(string(data.Fit[:]), "\x00")
	fit = strings.TrimSpace(fit)
	status := strings.TrimRight(string(data.Status[:]), "\x00")
	status = strings.TrimSpace(status)
	id := strings.TrimRight(string(data.Id[:]), "\x00")
	id = strings.TrimSpace(id)

	fmt.Printf(
		"║ %-2d ║ %-12s ║ %-4s ║ %-4s ║ %9d ║ %10d ║ %-7s ║ %-6s ║\n",
		idx+1,
		name,
		typeChar,
		fit,
		data.Start,
		data.Size,
		status,
		id,
	)
}

type PartitionFreeSpace struct {
	Start int32
	Size  int32
}

//  =============================================================

type Superblock struct {
	S_filesystem_type   int32 // ext2, ext3
	S_inodes_count      int32 // total number of inodes
	S_blocks_count      int32 // total number of blocks
	S_free_blocks_count int32 // free blocks
	S_free_inodes_count int32 // free inodes
	S_mtime             [17]byte
	S_umtime            [17]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_fist_ino          int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

//  =============================================================

type Inode struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [16]byte
	I_ctime [16]byte
	I_mtime [16]byte
	I_block [15]int32
	I_type  [1]byte
	I_perm  [3]byte
}

//  =============================================================

type Fileblock struct {
	B_content [64]byte
}

//  =============================================================

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

type Folderblock struct {
	B_content [4]Content
}

//  =============================================================

type Pointerblock struct {
	B_pointers [16]int32
}

//  =============================================================

type Content_J struct {
	Operation [10]byte
	Path      [100]byte
	Content   [100]byte
	Date      [17]byte
}

type Journaling struct {
	Size      int32
	Ultimo    int32
	Contenido [50]Content_J
}
