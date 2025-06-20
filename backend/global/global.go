package globals

import (
	Structs "MIA_PI_202001151_1VAC1S2025/structs"
	"regexp"
)

var Commands = map[string]string{
	"mkdisk":  "Crear un disco virtual",
	"rmdisk":  "Eliminar un disco virtual",
	"fdisk":   "Crear o modificar particiones en un disco virtual",
	"mount":   "Montar un sistema de archivos",
	"unmount": "Desmontar un sistema de archivos",
	"mkfs":    "Formatear un sistema de archivos",
	"login":   "Iniciar sesión en el sistema",
	"logout":  "Cerrar sesión en el sistema",
	"mkgrp":   "Crear un grupo de usuarios",
	"rmgrp":   "Eliminar un grupo de usuarios",
	"mkusr":   "Crear un usuario",
	"rmusr":   "Eliminar un usuario",
	"mkfile":  "Crear un archivo en el sistema de archivos",
	"cat":     "Leer el contenido de un archivo",
	"mkdir":   "Crear un directorio",
	"find":    "Buscar archivos o directorios",
	"pause":   "Pausar la ejecución del programa",
	"execute": "Ejecutar un script o comando",
	"rep":     "Generar un reporte del sistema de archivos",
}

var Reports = []string{
	"mbr",
	"disk",
	"inode",
	"block",
	"bm_inode",
	"bm_block",
	"tree",
	"sb",
	"file",
	"ls",
	"Journaling",
}

var Re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

var PathDisks = "./disks/MIA/P1/"
var LoginSession Structs.LoginSession = Structs.LoginSession{
	User:        "",
	Password:    "",
	PartitionID: [4]byte{},
}
