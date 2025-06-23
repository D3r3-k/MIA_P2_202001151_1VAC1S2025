package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"
)

func Fn_Rep(params string) (string, error) {
	paramDefs := map[string]Structs.ParamDef{
		"-name": {Required: true},
		"-path": {Required: true},
		"-id":   {Required: true},
		"-ruta": {Required: false},
	}
	parsed, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return "", err
	}
	name := strings.ToLower(parsed["-name"])
	path := parsed["-path"]
	id := parsed["-id"]
	ruta := parsed["-ruta"]

	validReport := false
	for _, report := range globals.Reports {
		if strings.ToLower(report) == name {
			validReport = true
			break
		}
	}
	if !validReport {
		utils.ShowMessage("Reporte inválido.", true)
		return "", fmt.Errorf("reporte inválido: %s", name)
	}

	return rep(name, path, id, ruta)
}

// rep -path=<ruta> -name=<nombre> -id=<id> [-ruta=<ruta>]
func rep(name, path, id, ruta string) (string, error) {
	switch name {
	case "mbr":
		err := reporteMBR(path, id)
		if err != nil {
			utils.ShowMessageCustom("Error REP", "No se pudo generar el reporte MBR: "+err.Error())
			return "", err
		} else {
			utils.ShowMessageCustom("Reporte generado", "El reporte [MBR] ha sido generado exitosamente en: "+path)
			return "El reporte [MBR] ha sido generado exitosamente en: " + path, nil
		}
	case "disk":
		err := reporteDISK(path, id)
		if err != nil {
			utils.ShowMessageCustom("Error REP", "No se pudo generar el reporte de disco: "+err.Error())
			return "", err
		} else {
			utils.ShowMessageCustom("Reporte generado", "El reporte [DISK] ha sido generado exitosamente en: "+path)
			return "El reporte [DISK] ha sido generado exitosamente en: " + path, nil
		}
	case "inode":
		err := reporteINODE(path, id)
		if err != nil {
			utils.ShowMessageCustom("Error REP", "No se pudo generar el reporte de inodo: "+err.Error())
			return "", err
		} else {
			utils.ShowMessageCustom("Reporte generado", "El reporte [INODE] ha sido generado exitosamente en: "+path)
			return "El reporte [INODE] ha sido generado exitosamente en: " + path, nil
		}
	case "block":
		err := reporteBLOCK(path, id)
		if err != nil {
			utils.ShowMessageCustom("Error REP", "No se pudo generar el reporte de bloque: "+err.Error())
			return "", err
		} else {
			utils.ShowMessageCustom("Reporte generado", "El reporte [BLOCK] ha sido generado exitosamente en: "+path)
			return "El reporte [BLOCK] ha sido generado exitosamente en: " + path, nil
		}
	case "bm_inode":
		err := reporteBM_INODE(path, id)
		if err != nil {
			utils.ShowMessageCustom("Error REP", "No se pudo generar el reporte bm_inode: "+err.Error())
			return "", err
		} else {
			utils.ShowMessageCustom("Reporte generado", "El reporte [BM_INODE] ha sido generado exitosamente en: "+path)
			return "El reporte [BM_INODE] ha sido generado exitosamente en: " + path, nil
		}
	case "bm_block":
		err := reporteBM_BLOCK(path, id)
		if err != nil {
			utils.ShowMessageCustom("Error REP", "No se pudo generar el reporte bm_block: "+err.Error())
			return "", err
		} else {
			utils.ShowMessageCustom("Reporte generado", "El reporte [BM_BLOCK] ha sido generado exitosamente en: "+path)
			return "El reporte [BM_BLOCK] ha sido generado exitosamente en: " + path, nil
		}
	case "tree":
		err := reporteTREE(path, id)
		if err != nil {
			utils.ShowMessageCustom("Error REP", "No se pudo generar el reporte tree: "+err.Error())
			return "", err
		} else {
			utils.ShowMessageCustom("Reporte generado", "El reporte [TREE] ha sido generado exitosamente en: "+path)
			return "El reporte [TREE] ha sido generado exitosamente en: " + path, nil
		}
	case "sb":
		err := reporteSB(path, id)
		if err != nil {
			utils.ShowMessageCustom("Error REP", "No se pudo generar el reporte sb: "+err.Error())
			return "", err
		} else {
			utils.ShowMessageCustom("Reporte generado", "El reporte [SB] ha sido generado exitosamente en: "+path)
			return "El reporte [SB] ha sido generado exitosamente en: " + path, nil
		}
	case "file":
		err := reporteFILE(path, id, ruta)
		if err != nil {
			utils.ShowMessageCustom("Error REP", "No se pudo generar el reporte file: "+err.Error())
			return "", err
		} else {
			utils.ShowMessageCustom("Reporte generado", "El reporte [FILE] ha sido generado exitosamente en: "+path)
			return "El reporte [FILE] ha sido generado exitosamente en: " + path, nil
		}

	case "ls":
		err := reporteLS(path, id, ruta)
		if err != nil {
			utils.ShowMessageCustom("Error REP", "No se pudo generar el reporte ls: "+err.Error())
			return "", err
		} else {
			utils.ShowMessageCustom("Reporte generado", "El reporte [LS] ha sido generado exitosamente en: "+path)
			return "El reporte [LS] ha sido generado exitosamente en: " + path, nil
		}
	case "journaling":
		err := reporteJOURNALING(path, id)
		if err != nil {
			utils.ShowMessageCustom("Error REP", "No se pudo generar el reporte journaling: "+err.Error())
			return "", err
		} else {
			utils.ShowMessageCustom("Reporte generado", "El reporte [JOURNALING] ha sido generado exitosamente en: "+path)
			return "El reporte [JOURNALING] ha sido generado exitosamente en: " + path, nil
		}
	default:
		utils.ShowMessageCustom("Error REP", "Reporte no implementado: "+name)
		return "", fmt.Errorf("reporte no implementado: %s", name)
	}
}

// TODO: rep -path=<ruta> -id=<id> -name=mbr
func reporteMBR(path string, id string) error {
	part := utils.GetPartitionById(id)
	if part == nil {
		return fmt.Errorf("partición no encontrada con id [%s]", id)
	}

	drive := strings.ToUpper(string(part.Id[0]))
	file, _, err := utils.OpenDisk(drive)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	var mbr Structs.MBR
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		return fmt.Errorf("no se pudo leer el MBR: %v", err)
	}

	dotPath, dotFormat, fdot, cleanup, err := utils.PrepareDotFile(path)
	if err != nil {
		return err
	}
	defer cleanup()

	fmt.Fprintln(fdot, "digraph MBR {")
	fmt.Fprintln(fdot, "rankdir=LR;")
	fmt.Fprintln(fdot, "node [shape=plaintext fontname=\"Consolas\"];")
	fmt.Fprintln(fdot, "MBRtable [label=<")
	fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='5' cellpadding='12' bgcolor='#E0FFFF'>")
	fmt.Fprintln(fdot, "<tr><td colspan='2' bgcolor='#26c6da' height='60' width='440'><b><font point-size='24'>Reporte MBR</font></b></td></tr>")

	fmt.Fprintf(fdot, "<tr><td bgcolor='#B2EBF2' width='180'><b>mbr_tamano</b></td><td width='260'>%d</td></tr>\n", mbr.MbrSize)
	fmt.Fprintf(fdot, "<tr><td bgcolor='#B2EBF2' width='180'><b>mbr_fecha_creacion</b></td><td width='260'>%s</td></tr>\n", utils.CleanDOTString(string(mbr.CreationDate[:])))
	fmt.Fprintf(fdot, "<tr><td bgcolor='#B2EBF2' width='180'><b>mbr_disk_signature</b></td><td width='260'>%d</td></tr>\n", mbr.Signature)

	for i, p := range mbr.Partitions {
		bgcolor := "#E0FFFF"
		if i%2 == 1 {
			bgcolor = "white"
		}
		fmt.Fprintf(fdot, "<tr><td colspan='2' bgcolor='#00BCD4' height='45' width='440'><b>Partición %d</b></td></tr>\n", i+1)
		fmt.Fprintf(fdot, "<tr><td bgcolor='%s' width='180'><b>part_status</b></td><td width='260'>%s</td></tr>\n", bgcolor, utils.CleanDOTString(string(p.Status[:])))
		fmt.Fprintf(fdot, "<tr><td bgcolor='%s' width='180'><b>part_type</b></td><td width='260'>%s</td></tr>\n", bgcolor, utils.CleanDOTString(string(p.Type[:])))
		fmt.Fprintf(fdot, "<tr><td bgcolor='%s' width='180'><b>part_fit</b></td><td width='260'>%s</td></tr>\n", bgcolor, utils.CleanDOTString(string(p.Fit[:])))
		fmt.Fprintf(fdot, "<tr><td bgcolor='%s' width='180'><b>part_start</b></td><td width='260'>%d</td></tr>\n", bgcolor, p.Start)
		fmt.Fprintf(fdot, "<tr><td bgcolor='%s' width='180'><b>part_size</b></td><td width='260'>%d</td></tr>\n", bgcolor, p.Size)
		fmt.Fprintf(fdot, "<tr><td bgcolor='%s' width='180'>part_name</td><td width='260'>%s</td></tr>\n", bgcolor, utils.CleanDOTString(string(p.Name[:])))
	}
	fmt.Fprintln(fdot, "</table>>];")
	fmt.Fprintln(fdot, "}")

	fdot.Close()

	return utils.GenerateGraphvizReport(dotPath, dotFormat, path)
}

// TODO: rep -path=<ruta> -id=<id> -name=disk
func reporteDISK(path string, id string) error {
	drive, err := utils.GetDriveletter(id)
	if err != nil {
		return err
	}
	file, _, err := utils.OpenDisk(drive)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	var mbr Structs.MBR
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		return fmt.Errorf("no se pudo leer el MBR: %v", err)
	}

	dotPath, dotFormat, fdot, _, err := utils.PrepareDotFile(path)
	if err != nil {
		return err
	}

	fmt.Fprintln(fdot, "digraph G {")
	fmt.Fprintln(fdot, "rankdir=LR;")
	fmt.Fprintln(fdot, "node [shape=plaintext fontname=\"Consolas\"];")
	fmt.Fprintln(fdot, "disk [label=<")
	fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='5' cellpadding='10'>")
	fmt.Fprintln(fdot, "<tr>")
	fmt.Fprintln(fdot, "<td bgcolor='#a5d6a7'><b>MBR</b></td>")

	mbrSize := int32(binary.Size(Structs.MBR{}))
	diskEnd := mbr.MbrSize
	diskUsable := diskEnd - mbrSize

	type partInfo struct {
		Part  Structs.Partition
		Index int
	}
	var parts []partInfo
	for i, p := range mbr.Partitions {
		if p.Size > 0 {
			parts = append(parts, partInfo{p, i})
		}
	}
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Part.Start < parts[j].Part.Start
	})

	pos := mbrSize
	writeFree := func(start, end int32) {
		size := end - start
		if size > 0 {
			percent := float64(size) / float64(diskUsable) * 100
			if percent > 0 {
				fmt.Fprintf(fdot, "<td bgcolor='#fffde7'><b>Libre<br/>%.2f%% del área útil</b></td>", percent)
			}
		}
	}

	for _, pinfo := range parts {
		p := pinfo.Part
		if p.Start > pos {
			writeFree(pos, p.Start)
		}
		percent := float64(p.Size) / float64(diskUsable) * 100
		name := utils.CleanDOTString(string(p.Name[:]))
		tipo := strings.ToUpper(utils.CleanDOTString(string(p.Type[:])))

		if tipo == "P" {
			fmt.Fprintf(fdot, "<td bgcolor='#81d4fa'><b>Primaria<br/>%s<br/>%.2f%%</b></td>", name, percent)
		} else if tipo == "E" {
			fmt.Fprintf(fdot, "<td bgcolor='#ffd54f'><b>Extendida<br/>%s<br/>%.2f%%</b></td>", name, percent)
		} else {
			fmt.Fprintf(fdot, "<td bgcolor='#e1bee7'><b>Desconocida<br/>%s<br/>%.2f%%</b></td>", name, percent)
		}
		pos = p.Start + p.Size
	}
	if pos < diskEnd {
		writeFree(pos, diskEnd)
	}

	fmt.Fprintln(fdot, "</tr>")
	fmt.Fprintln(fdot, "</table>>];")
	fmt.Fprintln(fdot, "}")

	fdot.Close()

	return utils.GenerateGraphvizReport(dotPath, dotFormat, path)
}

// TODO: rep -path=<ruta> -id=<id> -name=inode
func reporteINODE(path string, id string) error {
	drive, err := utils.GetDriveletter(id)
	if err != nil {
		return err
	}
	file, _, err := utils.OpenDisk(drive)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	part := utils.GetPartitionById(id)
	if part == nil {
		return fmt.Errorf("partición no encontrada con id [%s]", id)
	}

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return fmt.Errorf("no se pudo leer el superbloque: %v", err)
	}

	inodeCount := sb.S_inodes_count
	bitmapStart := sb.S_bm_inode_start
	inodeStart := sb.S_inode_start
	inodeSize := sb.S_inode_size

	bitmap := make([]byte, inodeCount)
	_, err = file.ReadAt(bitmap, int64(bitmapStart))
	if err != nil {
		return fmt.Errorf("no se pudo leer el bitmap de inodos: %v", err)
	}

	dotPath, dotFormat, fdot, _, err := utils.PrepareDotFile(path)
	if err != nil {
		return err
	}

	fmt.Fprintln(fdot, "digraph INODES {")
	fmt.Fprintln(fdot, "rankdir=TB;")
	fmt.Fprintln(fdot, "node [shape=plaintext fontname=\"Consolas\"];")

	for i := int32(0); i < inodeCount; i++ {
		if bitmap[i] == '1' {
			var inode Structs.Inode
			offset := int64(inodeStart) + int64(i)*int64(inodeSize)
			if err := utils.ReadObject(file, &inode, offset); err != nil {
				return fmt.Errorf("no se pudo leer el inodo %d: %v", i, err)
			}
			fmt.Fprintf(fdot, "inode%d [label=<\n", i)
			fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='5' cellpadding='6' bgcolor='#f5f5f5'>")
			fmt.Fprintf(fdot, "<tr><td colspan='2' bgcolor='#00bcd4'><b>Inodo %d</b></td></tr>\n", i)
			fmt.Fprintf(fdot, "<tr><td><b>UID</b></td><td>%d</td></tr>\n", inode.I_uid)
			fmt.Fprintf(fdot, "<tr><td><b>GID</b></td><td>%d</td></tr>\n", inode.I_gid)
			fmt.Fprintf(fdot, "<tr><td><b>Tamaño (bytes)</b></td><td>%d</td></tr>\n", inode.I_size)
			fmt.Fprintf(fdot, "<tr><td><b>atime</b></td><td>%s</td></tr>\n", string(inode.I_atime[:16]))
			fmt.Fprintf(fdot, "<tr><td><b>ctime</b></td><td>%s</td></tr>\n", string(inode.I_ctime[:16]))
			fmt.Fprintf(fdot, "<tr><td><b>mtime</b></td><td>%s</td></tr>\n", string(inode.I_mtime[:16]))
			tipo := "Desconocido"
			if inode.I_type[0] == '0' {
				tipo = "Carpeta"
			} else if inode.I_type[0] == '1' {
				tipo = "Archivo"
			}
			fmt.Fprintf(fdot, "<tr><td><b>Tipo</b></td><td>%s</td></tr>\n", tipo)
			fmt.Fprintf(fdot, "<tr><td><b>Permisos</b></td><td>%s</td></tr>\n", string(inode.I_perm[:]))
			for j := 0; j < 14; j++ {
				fmt.Fprintf(fdot, "<tr><td><b>AD</b></td><td>%d</td></tr>\n", inode.I_block[j])
			}
			fmt.Fprintf(fdot, "<tr><td><b>AI</b></td><td>%d</td></tr>\n", inode.I_block[14])
			fmt.Fprintln(fdot, "</table>>];")
		}
	}
	fmt.Fprintln(fdot, "}")

	fdot.Close()

	return utils.GenerateGraphvizReport(dotPath, dotFormat, path)
}

// TODO: rep -path=<ruta> -id=<id> -name=block
func reporteBLOCK(path string, id string) error {
	drive, err := utils.GetDriveletter(id)
	if err != nil {
		return err
	}
	file, _, err := utils.OpenDisk(drive)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	part := utils.GetPartitionById(id)
	if part == nil {
		return fmt.Errorf("partición no encontrada con id [%s]", id)
	}

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return fmt.Errorf("no se pudo leer el superbloque: %v", err)
	}

	blockCount := sb.S_blocks_count
	bitmapStart := sb.S_bm_block_start
	blockStart := sb.S_block_start
	blockSize := sb.S_block_size

	bitmap := make([]byte, blockCount)
	_, err = file.ReadAt(bitmap, int64(bitmapStart))
	if err != nil {
		return fmt.Errorf("no se pudo leer el bitmap de bloques: %v", err)
	}

	tipoBloque := make(map[int32]string)
	for i := int32(0); i < sb.S_inodes_count; i++ {
		var inode Structs.Inode
		offsetInodo := int64(sb.S_inode_start) + int64(i)*int64(sb.S_inode_size)
		if err := utils.ReadObject(file, &inode, offsetInodo); err == nil {
			if inode.I_type[0] == '0' || inode.I_type[0] == '1' {
				for j := 0; j < 14; j++ {
					idx := inode.I_block[j]
					if idx >= 0 {
						if inode.I_type[0] == '0' {
							tipoBloque[idx] = "carpeta"
						} else if inode.I_type[0] == '1' {
							tipoBloque[idx] = "archivo"
						}
					}
				}
				idx := inode.I_block[14]
				if idx >= 0 {
					tipoBloque[idx] = "apuntador"
				}
			}
		}
	}

	dotPath, dotFormat, fdot, _, err := utils.PrepareDotFile(path)
	if err != nil {
		return err
	}

	fmt.Fprintln(fdot, "digraph BLOCKS {")
	fmt.Fprintln(fdot, "rankdir=TB;")
	fmt.Fprintln(fdot, "node [shape=plaintext fontname=\"Consolas\"];")

	yaImpreso := make(map[int32]bool)

	for i := int32(0); i < blockCount; i++ {
		if bitmap[i] == '1' {
			if yaImpreso[i] {
				continue
			}
			offset := int64(blockStart) + int64(i)*int64(blockSize)
			tipo, ok := tipoBloque[i]

			if ok && tipo == "carpeta" {
				var folder Structs.Folderblock
				if err := utils.ReadObject(file, &folder, offset); err == nil {
					fmt.Fprintf(fdot, "block%d [label=<\n", i)
					fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='5' cellpadding='6' bgcolor='#fffde7'>")
					fmt.Fprintf(fdot, "<tr><td colspan='2' bgcolor='#ffb300'><b>Bloque Carpeta %d</b></td></tr>\n", i)
					for _, entry := range folder.B_content {
						name := strings.Trim(string(entry.B_name[:]), "\x00")
						if entry.B_inodo != -1 && name != "" {
							fmt.Fprintf(fdot, "<tr><td>%s</td><td>%d</td></tr>", name, entry.B_inodo)
						} else {
							fmt.Fprintf(fdot, "<tr><td>%s</td><td>%d</td></tr>", "--", -1)
						}
					}
					fmt.Fprintln(fdot, "</table>>];")
				}
			} else if ok && tipo == "archivo" {
				var fileblock Structs.Fileblock
				if err := utils.ReadObject(file, &fileblock, offset); err == nil {
					content := strings.Trim(string(fileblock.B_content[:]), "\x00")
					content = strings.ReplaceAll(content, "\n", "<br/>")
					fmt.Fprintf(fdot, "block%d [label=<\n", i)
					fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='5' cellpadding='6' bgcolor='#c5e1a5'>")
					fmt.Fprintf(fdot, "<tr><td bgcolor='#388e3c'><b>Bloque Archivo %d</b></td></tr>\n", i)
					fmt.Fprintf(fdot, "<tr><td align='left'><font face=\"Consolas\">%s</font></td></tr>\n", content)
					fmt.Fprintln(fdot, "</table>>];")
				}
			} else if ok && tipo == "apuntador" {
				var pointer Structs.Pointerblock
				if err := utils.ReadObject(file, &pointer, offset); err == nil {
					fmt.Fprintf(fdot, "block%d [label=<\n", i)
					fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='5' cellpadding='6' bgcolor='#ce93d8'>")
					fmt.Fprintf(fdot, "<tr><td colspan='2' bgcolor='#8e24aa'><b>Bloque Apuntador %d</b></td></tr>\n", i)
					for j, ptr := range pointer.B_pointers {
						if ptr != -1 {
							fmt.Fprintf(fdot, "<tr><td>Puntero %d</td><td>%d</td></tr>", j, ptr)
						} else {
							fmt.Fprintf(fdot, "<tr><td>Puntero %d</td><td>-1</td></tr>", j)
						}
					}
					fmt.Fprintln(fdot, "</table>>];")
					for _, ptr := range pointer.B_pointers {
						if ptr != -1 && !yaImpreso[ptr] {
							offTarget := int64(blockStart) + int64(ptr)*int64(blockSize)
							var fileblock Structs.Fileblock
							if err := utils.ReadObject(file, &fileblock, offTarget); err == nil {
								content := strings.Trim(string(fileblock.B_content[:]), "\x00")
								content = strings.ReplaceAll(content, "\n", "<br/>")
								fmt.Fprintf(fdot, "block%d [label=<\n", ptr)
								fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='5' cellpadding='6' bgcolor='#dcedc8'>")
								fmt.Fprintf(fdot, "<tr><td bgcolor='#388e3c'><b>Bloque Archivo %d</b></td></tr>\n", ptr)
								fmt.Fprintf(fdot, "<tr><td align='left'><font face=\"Consolas\">%s</font></td></tr>\n", content)
								fmt.Fprintln(fdot, "</table>>];")
								fmt.Fprintf(fdot, "block%d -> block%d;\n", i, ptr)
								yaImpreso[ptr] = true
								continue
							}
							var folder Structs.Folderblock
							if err := utils.ReadObject(file, &folder, offTarget); err == nil {
								fmt.Fprintf(fdot, "block%d [label=<\n", ptr)
								fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='5' cellpadding='6' bgcolor='#ffe0b2'>")
								fmt.Fprintf(fdot, "<tr><td colspan='2' bgcolor='#ffb300'><b>Bloque Carpeta %d</b></td></tr>\n", ptr)
								for _, entry := range folder.B_content {
									name := strings.Trim(string(entry.B_name[:]), "\x00")
									if entry.B_inodo != -1 && name != "" {
										fmt.Fprintf(fdot, "<tr><td>%s</td><td>%d</td></tr>", name, entry.B_inodo)
									} else {
										fmt.Fprintf(fdot, "<tr><td>%s</td><td>%d</td></tr>", "--", -1)
									}
								}
								fmt.Fprintln(fdot, "</table>>];")
								fmt.Fprintf(fdot, "block%d -> block%d;\n", i, ptr)
								yaImpreso[ptr] = true
							}
						}
					}
				}
			} else {
				fmt.Fprintf(fdot, "block%d [label=<\n", i)
				fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='5' cellpadding='6' bgcolor='#ececec'>")
				fmt.Fprintf(fdot, "<tr><td><b>Bloque %d</b></td></tr>\n", i)
				fmt.Fprintln(fdot, "<tr><td><i>Desconocido o no enlazado por inodos activos</i></td></tr>")
				fmt.Fprintln(fdot, "</table>>];")
			}
		}
	}
	fmt.Fprintln(fdot, "}")
	fdot.Close()

	return utils.GenerateGraphvizReport(dotPath, dotFormat, path)
}

// TODO: rep -path=<ruta> -id=<id> -name=bm_inode
func reporteBM_INODE(path string, id string) error {
	if !strings.HasSuffix(path, ".txt") {
		return fmt.Errorf("el archivo de reporte debe tener extensión .txt")
	}
	part := utils.GetPartitionById(id)
	if part == nil {
		return fmt.Errorf("partición no encontrada con id [%s]", id)
	}

	drive := strings.ToUpper(string(part.Id[0]))
	file, _, err := utils.OpenDisk(drive)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return fmt.Errorf("no se pudo leer el superbloque: %v", err)
	}

	inodeCount := sb.S_inodes_count
	bitmapStart := sb.S_bm_inode_start

	bitmap := make([]byte, inodeCount)
	_, err = file.ReadAt(bitmap, int64(bitmapStart))
	if err != nil {
		return fmt.Errorf("no se pudo leer el bitmap de inodos: %v", err)
	}

	if err := utils.CreateFile(path); err != nil {
		return fmt.Errorf("no se pudo crear el archivo de reporte: %v", err)
	}

	out, err := utils.OpenFile(path)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo de reporte: %v", err)
	}
	defer out.Close()

	for i := int32(0); i < inodeCount; i++ {
		fmt.Fprintf(out, "%c", bitmap[i])
		if (i+1)%20 == 0 || i == inodeCount-1 {
			fmt.Fprintln(out)
		}
	}

	return nil
}

// TODO: rep -path=<ruta> -id=<id> -name=bm_block
func reporteBM_BLOCK(path string, id string) error {
	if !strings.HasSuffix(path, ".txt") {
		return fmt.Errorf("el archivo de reporte debe tener extensión .txt")
	}

	part := utils.GetPartitionById(id)
	if part == nil {
		return fmt.Errorf("partición no encontrada con id [%s]", id)
	}

	drive := strings.ToUpper(string(part.Id[0]))
	file, _, err := utils.OpenDisk(drive)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return fmt.Errorf("no se pudo leer el superbloque: %v", err)
	}

	blockCount := sb.S_blocks_count
	bitmapStart := sb.S_bm_block_start

	bitmap := make([]byte, blockCount)
	_, err = file.ReadAt(bitmap, int64(bitmapStart))
	if err != nil {
		return fmt.Errorf("no se pudo leer el bitmap de bloques: %v", err)
	}

	if err := utils.CreateFile(path); err != nil {
		return fmt.Errorf("no se pudo crear el archivo de reporte: %v", err)
	}

	out, err := utils.OpenFile(path)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo de reporte: %v", err)
	}
	defer out.Close()

	for i := int32(0); i < blockCount; i++ {
		fmt.Fprintf(out, "%c", bitmap[i])
		if (i+1)%20 == 0 || i == blockCount-1 {
			fmt.Fprintln(out)
		}
	}

	return nil
}

// TODO: rep -path=<ruta> -id=<id> -name=tree
func reporteTREE(path string, id string) error {
	part := utils.GetPartitionById(id)
	if part == nil {
		return fmt.Errorf("partición no encontrada con id [%s]", id)
	}
	drive := strings.ToUpper(string(part.Id[0]))
	file, _, err := utils.OpenDisk(drive)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return fmt.Errorf("no se pudo leer el superbloque: %v", err)
	}

	dotPath, dotFormat, fdot, cleanup, err := utils.PrepareDotFile(path)
	if err != nil {
		return err
	}
	defer cleanup()

	fmt.Fprintln(fdot, "digraph TREE {")
	fmt.Fprintln(fdot, "rankdir=LR;")
	fmt.Fprintln(fdot, "node [fontname=\"Consolas\"];")
	fmt.Fprintln(fdot, "graph [nodesep=1, ranksep=1];")

	visitedInodes := make(map[int32]bool)
	visitedBlocks := make(map[int32]bool)

	portName := func(name string) string {
		name = strings.ReplaceAll(name, ".", "")
		name = strings.ReplaceAll(name, " ", "_")
		name = strings.ReplaceAll(name, "-", "_")
		if name == "" {
			name = "vacio"
		}
		return strings.ToLower(name)
	}

	printInodoNode := func(i int32, inode Structs.Inode) {
		tipo := "Desconocido"
		if inode.I_type[0] == '0' {
			tipo = "Carpeta"
		} else if inode.I_type[0] == '1' {
			tipo = "Archivo"
		}
		fmt.Fprintf(fdot, "inode%d [shape=\"record\", label=\"{", i)
		fmt.Fprintf(fdot, "{ <header> Inodo %d | { UID |<uid> %d} | { GID |<gid> %d} | { Tamaño |<tam> %d} | { atime |<atime> %s} | { ctime |<ctime> %s} | { mtime |<mtime> %s} | { Tipo |<tipo> %s} | { Permisos |<perm> %s}",
			i, inode.I_uid, inode.I_gid, inode.I_size,
			strings.Trim(string(inode.I_atime[:]), "\x00"),
			strings.Trim(string(inode.I_ctime[:]), "\x00"),
			strings.Trim(string(inode.I_mtime[:]), "\x00"),
			tipo, strings.Trim(string(inode.I_perm[:]), "\x00"),
		)
		for j := 0; j < 15; j++ {
			if j < 14 {
				fmt.Fprintf(fdot, " | {AD%-2d |<block%d> %d}", j, j, inode.I_block[j])
			} else {
				fmt.Fprintf(fdot, " | {AI%-2d |<block%d> %d}", j, j, inode.I_block[j])

			}
		}
		fmt.Fprintf(fdot, " }")
		fmt.Fprintf(fdot, "}\"];\n")
	}

	printFolderBlock := func(idx int32, folder Structs.Folderblock) {
		fmt.Fprintf(fdot, "block%d [shape=\"record\", label=\"{", idx)
		fmt.Fprintf(fdot, "{ <header> Bloque Carpeta %d", idx)
		for entryIdx, entry := range folder.B_content {
			name := strings.Trim(string(entry.B_name[:]), "\x00")
			if name == "" {
				name = "--"
			}
			fmt.Fprintf(fdot, " | { %s |<%s%d> %d}", name, portName(name), entryIdx, entry.B_inodo)
		}
		fmt.Fprintf(fdot, " }")
		fmt.Fprintf(fdot, "}\"];\n")
	}

	printFileBlock := func(idx int32, fileblock Structs.Fileblock) {
		content := strings.Trim(string(fileblock.B_content[:]), "\x00")
		content = strings.ReplaceAll(content, "\"", "\\\"")
		content = strings.ReplaceAll(content, "\n", "\\n")
		fmt.Fprintf(fdot, "block%d [shape=\"record\", label=\"{", idx)
		fmt.Fprintf(fdot, "{ <header> Bloque Archivo %d | <contenido> %s }", idx, content)
		fmt.Fprintf(fdot, "}\"];\n")
	}

	printPointerBlock := func(idx int32, pointer Structs.Pointerblock) {
		fmt.Fprintf(fdot, "block%d [shape=\"record\", label=\"{", idx)
		fmt.Fprintf(fdot, "{ <header> Bloque Apuntador %d", idx)
		for j, ptr := range pointer.B_pointers {
			fmt.Fprintf(fdot, " | { Puntero %d| <ptr%d> %d}", j, j, ptr)
		}
		fmt.Fprintf(fdot, " }")
		fmt.Fprintf(fdot, "}\"];\n")
	}

	var processPointerBlock func(idx int32, parentInode int32)
	var processInode func(inodeIdx int32)

	processInode = func(inodeIdx int32) {
		if visitedInodes[inodeIdx] {
			return
		}
		visitedInodes[inodeIdx] = true

		inodeOffset := int64(sb.S_inode_start) + int64(inodeIdx)*int64(sb.S_inode_size)
		var inode Structs.Inode
		if err := utils.ReadObject(file, &inode, inodeOffset); err != nil {
			return
		}
		printInodoNode(inodeIdx, inode)

		for j := 0; j < 14; j++ {
			blkIdx := inode.I_block[j]
			if blkIdx < 0 {
				continue
			}
			if !visitedBlocks[blkIdx] {
				visitedBlocks[blkIdx] = true
				offset := int64(sb.S_block_start) + int64(blkIdx)*int64(sb.S_block_size)
				if inode.I_type[0] == '0' {
					var folder Structs.Folderblock
					if err := utils.ReadObject(file, &folder, offset); err == nil {
						printFolderBlock(blkIdx, folder)
						fmt.Fprintf(fdot, "inode%d:block%d -> block%d:header;\n", inodeIdx, j, blkIdx)
						for entryIdx, entry := range folder.B_content {
							entryName := strings.Trim(string(entry.B_name[:]), "\x00")
							if entry.B_inodo >= 0 && entryName != "" && entryName != "." && entryName != ".." {
								fmt.Fprintf(fdot, "block%d:%s%d -> inode%d:header;\n", blkIdx, portName(entryName), entryIdx, entry.B_inodo)
								processInode(entry.B_inodo)
							}
						}
					}
				} else if inode.I_type[0] == '1' {
					var fileblock Structs.Fileblock
					if err := utils.ReadObject(file, &fileblock, offset); err == nil {
						printFileBlock(blkIdx, fileblock)
						fmt.Fprintf(fdot, "inode%d:block%d -> block%d:header;\n", inodeIdx, j, blkIdx)
					}
				}
			} else {
				fmt.Fprintf(fdot, "inode%d:block%d -> block%d:header;\n", inodeIdx, j, blkIdx)
			}
		}

		blkIdx := inode.I_block[14]
		if blkIdx >= 0 {
			processPointerBlock(blkIdx, inodeIdx)
		}
	}

	processPointerBlock = func(idx int32, parentInode int32) {
		if visitedBlocks[idx] {
			fmt.Fprintf(fdot, "inode%d:block14 -> block%d:header;\n", parentInode, idx)
			return
		}
		visitedBlocks[idx] = true

		offset := int64(sb.S_block_start) + int64(idx)*int64(sb.S_block_size)
		var pointer Structs.Pointerblock
		if err := utils.ReadObject(file, &pointer, offset); err != nil {
			return
		}
		printPointerBlock(idx, pointer)
		fmt.Fprintf(fdot, "inode%d:block14 -> block%d:header;\n", parentInode, idx)

		for j, ptr := range pointer.B_pointers {
			if ptr >= 0 {
				ptrOffset := int64(sb.S_block_start) + int64(ptr)*int64(sb.S_block_size)
				var fileblock Structs.Fileblock
				if err := utils.ReadObject(file, &fileblock, ptrOffset); err == nil {
					printFileBlock(ptr, fileblock)
					fmt.Fprintf(fdot, "block%d:ptr%d -> block%d:header;\n", idx, j, ptr)
					continue
				}
				var folder Structs.Folderblock
				if err := utils.ReadObject(file, &folder, ptrOffset); err == nil {
					printFolderBlock(ptr, folder)
					fmt.Fprintf(fdot, "block%d:ptr%d -> block%d:header;\n", idx, j, ptr)
					for entryIdx, entry := range folder.B_content {
						entryName := strings.Trim(string(entry.B_name[:]), "\x00")
						if entry.B_inodo >= 0 && entryName != "" && entryName != "." && entryName != ".." {
							fmt.Fprintf(fdot, "block%d:%s%d -> inode%d:header;\n", ptr, portName(entryName), entryIdx, entry.B_inodo)
							processInode(entry.B_inodo)
						}
					}
					continue
				}
				fmt.Fprintf(fdot, "block%d [label=\"Bloque %d desconocido\"];\n", ptr, ptr)
				fmt.Fprintf(fdot, "block%d:ptr%d -> block%d:header;\n", idx, j, ptr)
			}
		}
	}
	processInode(0)

	fmt.Fprintln(fdot, "}")
	fdot.Close()

	return utils.GenerateGraphvizReport(dotPath, dotFormat, path)
}

// TODO: rep -path=<ruta> -id=<id> -name=sb
func reporteSB(path string, id string) error {
	part := utils.GetPartitionById(id)
	if part == nil {
		return fmt.Errorf("partición no encontrada con id [%s]", id)
	}

	drive := strings.ToUpper(string(part.Id[0]))
	file, _, err := utils.OpenDisk(drive)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return fmt.Errorf("no se pudo leer el superbloque: %v", err)
	}

	dotPath, dotFormat, fdot, cleanup, err := utils.PrepareDotFile(path)
	if err != nil {
		return err
	}
	defer cleanup()

	fmt.Fprintln(fdot, "digraph SB {")
	fmt.Fprintln(fdot, "rankdir=LR;")
	fmt.Fprintln(fdot, "node [shape=plaintext fontname=\"Consolas\"];")
	fmt.Fprintln(fdot, "SBTable [label=<")
	fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='5' cellpadding='10' bgcolor='#F8BBD0'>")
	fmt.Fprintln(fdot, "<tr><td colspan='2' bgcolor='#D81B60'><b><font point-size='20' color='white'>Reporte de SuperBloque</font></b></td></tr>")
	fmt.Fprintf(fdot, "<tr><td><b>sb_nombre_hd</b></td><td>%s</td></tr>\n", drive+".dsk")

	fmt.Fprintf(fdot, "<tr><td><b>S_filesystem_type</b></td><td>%d</td></tr>\n", sb.S_filesystem_type)
	fmt.Fprintf(fdot, "<tr><td><b>S_inodes_count</b></td><td>%d</td></tr>\n", sb.S_inodes_count)
	fmt.Fprintf(fdot, "<tr><td><b>S_blocks_count</b></td><td>%d</td></tr>\n", sb.S_blocks_count)
	fmt.Fprintf(fdot, "<tr><td><b>S_free_blocks_count</b></td><td>%d</td></tr>\n", sb.S_free_blocks_count)
	fmt.Fprintf(fdot, "<tr><td><b>S_free_inodes_count</b></td><td>%d</td></tr>\n", sb.S_free_inodes_count)
	fmt.Fprintf(fdot, "<tr><td><b>S_mtime</b></td><td>%s</td></tr>\n", strings.Trim(string(sb.S_mtime[:]), " \x00"))
	fmt.Fprintf(fdot, "<tr><td><b>S_umtime</b></td><td>%s</td></tr>\n", strings.Trim(string(sb.S_umtime[:]), " \x00"))
	fmt.Fprintf(fdot, "<tr><td><b>S_mnt_count</b></td><td>%d</td></tr>\n", sb.S_mnt_count)
	fmt.Fprintf(fdot, "<tr><td><b>S_magic</b></td><td>%d</td></tr>\n", sb.S_magic)
	fmt.Fprintf(fdot, "<tr><td><b>S_inode_size</b></td><td>%d</td></tr>\n", sb.S_inode_size)
	fmt.Fprintf(fdot, "<tr><td><b>S_block_size</b></td><td>%d</td></tr>\n", sb.S_block_size)
	fmt.Fprintf(fdot, "<tr><td><b>S_fist_ino</b></td><td>%d</td></tr>\n", sb.S_fist_ino)
	fmt.Fprintf(fdot, "<tr><td><b>S_first_blo</b></td><td>%d</td></tr>\n", sb.S_first_blo)
	fmt.Fprintf(fdot, "<tr><td><b>S_bm_inode_start</b></td><td>%d</td></tr>\n", sb.S_bm_inode_start)
	fmt.Fprintf(fdot, "<tr><td><b>S_bm_block_start</b></td><td>%d</td></tr>\n", sb.S_bm_block_start)
	fmt.Fprintf(fdot, "<tr><td><b>S_inode_start</b></td><td>%d</td></tr>\n", sb.S_inode_start)
	fmt.Fprintf(fdot, "<tr><td><b>S_block_start</b></td><td>%d</td></tr>\n", sb.S_block_start)

	fmt.Fprintln(fdot, "</table>>];")
	fmt.Fprintln(fdot, "}")

	fdot.Close()

	return utils.GenerateGraphvizReport(dotPath, dotFormat, path)
}

// TODO: rep -path=<destino> -id=<id> -name=file -ruta=<pathlogico>
func reporteFILE(path, id, ruta string) error {
	if ruta == "" {
		return fmt.Errorf("debe especificar -ruta con la ruta del archivo dentro del sistema de archivos")
	}
	part := utils.GetPartitionById(id)
	if part == nil {
		return fmt.Errorf("partición no encontrada con id [%s]", id)
	}

	drive := strings.ToUpper(string(part.Id[0]))
	file, _, err := utils.OpenDisk(drive)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return fmt.Errorf("no se pudo leer el superbloque: %v", err)
	}

	var currentInode int32 = 0
	parts := strings.Split(ruta, "/")
	inodoSize := int32(binary.Size(Structs.Inode{}))

	for _, part := range parts {
		if part == "" {
			continue
		}
		encontrado, nextInode := utils.SearchDirectoryEntry(file, &sb, currentInode, part)
		if !encontrado {
			return fmt.Errorf("no se encontró la ruta: %s", ruta)
		}
		currentInode = nextInode
	}

	var inode Structs.Inode
	if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start+currentInode*inodoSize)); err != nil {
		return fmt.Errorf("no se pudo leer el inodo del archivo destino")
	}
	if inode.I_type[0] != '1' {
		return fmt.Errorf("la ruta especificada no apunta a un archivo regular")
	}

	var contenido string
	blockSizeFile := int32(binary.Size(Structs.Fileblock{}))
	for i := 0; i < 14; i++ {
		blockIdx := inode.I_block[i]
		if blockIdx == -1 {
			continue
		}
		var blk Structs.Fileblock
		offset := int64(sb.S_block_start + blockIdx*blockSizeFile)
		if err := utils.ReadObject(file, &blk, offset); err == nil {
			contenido += string(blk.B_content[:])
		}
	}

	if inode.I_block[14] != -1 {
		var pb Structs.Pointerblock
		offsetPB := int64(sb.S_block_start + inode.I_block[14]*blockSizeFile)
		if err := utils.ReadObject(file, &pb, offsetPB); err == nil {
			for _, ptr := range pb.B_pointers {
				if ptr != -1 {
					var blk Structs.Fileblock
					offset := int64(sb.S_block_start + ptr*blockSizeFile)
					if err := utils.ReadObject(file, &blk, offset); err == nil {
						contenido += string(blk.B_content[:])
					}
				}
			}
		}
	}

	contenido = strings.Trim(contenido, "\x00")
	if err := utils.CreateFile(path); err != nil {
		return fmt.Errorf("no se pudo crear el archivo de salida: %v", err)
	}
	out, err := utils.OpenFile(path)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo de salida: %v", err)
	}
	defer out.Close()

	fmt.Fprintln(out, contenido)
	return nil
}

// TODO: rep -path=<destino> -id=<id> -name=ls [-ruta=<carpeta>]
func reporteLS(path, id, ruta string) error {
	part := utils.GetPartitionById(id)
	if part == nil {
		return fmt.Errorf("partición no encontrada con id [%s]", id)
	}

	drive := strings.ToUpper(string(part.Id[0]))
	file, _, err := utils.OpenDisk(drive)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return fmt.Errorf("no se pudo leer el superbloque: %v", err)
	}

	carpetaRuta := "/"
	if ruta != "" {
		carpetaRuta = ruta
	}
	if carpetaRuta[0] != '/' {
		return fmt.Errorf("la ruta lógica debe ser absoluta (empezar con /)")
	}

	rutaLimpia := strings.Trim(carpetaRuta, "/")
	partes := []string{}
	if rutaLimpia != "" {
		partes = strings.Split(rutaLimpia, "/")
	}

	inodoActual := int32(0)
	encontrado := true

	for _, nombre := range partes {
		var inode Structs.Inode
		offset := int64(sb.S_inode_start) + int64(inodoActual)*int64(sb.S_inode_size)
		if err := utils.ReadObject(file, &inode, offset); err != nil {
			return fmt.Errorf("no se pudo leer inodo %d: %v", inodoActual, err)
		}
		if inode.I_type[0] != '0' {
			return fmt.Errorf("la ruta intermedia '%s' no es una carpeta", nombre)
		}
		encontrado = false
		for _, bIdx := range inode.I_block {
			if bIdx < 0 {
				continue
			}
			var folder Structs.Folderblock
			offBlk := int64(sb.S_block_start) + int64(bIdx)*int64(sb.S_block_size)
			if err := utils.ReadObject(file, &folder, offBlk); err == nil {
				for _, entry := range folder.B_content {
					n := strings.Trim(string(entry.B_name[:]), "\x00")
					if n == nombre {
						inodoActual = entry.B_inodo
						encontrado = true
						break
					}
				}
			}
			if encontrado {
				break
			}
		}
		if !encontrado {
			return fmt.Errorf("no se encontró la carpeta '%s' en la ruta", nombre)
		}
	}
	var inode Structs.Inode
	offset := int64(sb.S_inode_start) + int64(inodoActual)*int64(sb.S_inode_size)
	if err := utils.ReadObject(file, &inode, offset); err != nil {
		return fmt.Errorf("no se pudo leer el inodo destino: %v", err)
	}
	if inode.I_type[0] != '0' {
		return fmt.Errorf("la ruta dada no es una carpeta")
	}

	type Entry struct {
		Perm  string
		Owner string
		Grupo string
		Size  int
		Fecha string
		Hora  string
		Tipo  string
		Name  string
	}
	var entries []Entry

	for _, bIdx := range inode.I_block {
		if bIdx < 0 {
			continue
		}
		var folder Structs.Folderblock
		offBlk := int64(sb.S_block_start) + int64(bIdx)*int64(sb.S_block_size)
		if err := utils.ReadObject(file, &folder, offBlk); err == nil {
			for _, entry := range folder.B_content {
				n := strings.Trim(string(entry.B_name[:]), "\x00")
				if n == "" || n == "." || n == ".." || entry.B_inodo < 0 {
					continue
				}
				var inodoHijo Structs.Inode
				offInodoHijo := int64(sb.S_inode_start) + int64(entry.B_inodo)*int64(sb.S_inode_size)
				if err := utils.ReadObject(file, &inodoHijo, offInodoHijo); err != nil {
					continue
				}
				perms := "-"
				if inodoHijo.I_type[0] == '0' {
					perms = "d"
				}
				perms += utils.InodePermString(inodoHijo.I_perm[:])

				owner, grupo := utils.GetUserAndGroupNames(id, inodoHijo.I_uid, inodoHijo.I_gid)

				size := int(inodoHijo.I_size)

				fchRaw := strings.Trim(string(inodoHijo.I_mtime[:]), "\x00")

				fecha, hora := strings.Split(fchRaw, " ")[0], strings.Split(fchRaw, " ")[1]

				tipo := "Archivo"
				if inodoHijo.I_type[0] == '0' {
					tipo = "Carpeta"
				}
				entries = append(entries, Entry{
					Perm:  perms,
					Owner: owner,
					Grupo: grupo,
					Size:  size,
					Fecha: fecha,
					Hora:  hora,
					Tipo:  tipo,
					Name:  n,
				})
			}
		}
	}

	dotPath, dotFormat, fdot, cleanup, err := utils.PrepareDotFile(path)
	if err != nil {
		return err
	}
	defer cleanup()

	fmt.Fprintln(fdot, "digraph LS {")
	fmt.Fprintln(fdot, "node [shape=plaintext fontname=\"Consolas\"];")
	fmt.Fprintln(fdot, "LS [label=<")
	fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='2' cellpadding='8' bgcolor='#E3F2FD'>")
	fmt.Fprintln(fdot, "<tr>")
	fmt.Fprintln(fdot, "<td bgcolor='#1565C0'><font color='white'><b>Permisos</b></font></td>")
	fmt.Fprintln(fdot, "<td bgcolor='#1565C0'><font color='white'><b>Owner</b></font></td>")
	fmt.Fprintln(fdot, "<td bgcolor='#1565C0'><font color='white'><b>Grupo</b></font></td>")
	fmt.Fprintln(fdot, "<td bgcolor='#1565C0'><font color='white'><b>Size (en Bytes)</b></font></td>")
	fmt.Fprintln(fdot, "<td bgcolor='#1565C0'><font color='white'><b>Fecha</b></font></td>")
	fmt.Fprintln(fdot, "<td bgcolor='#1565C0'><font color='white'><b>Hora</b></font></td>")
	fmt.Fprintln(fdot, "<td bgcolor='#1565C0'><font color='white'><b>Tipo</b></font></td>")
	fmt.Fprintln(fdot, "<td bgcolor='#1565C0'><font color='white'><b>Name</b></font></td>")
	fmt.Fprintln(fdot, "</tr>")
	for _, e := range entries {
		fmt.Fprintf(fdot, "<tr><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
			e.Perm, e.Owner, e.Grupo, e.Size, e.Fecha, e.Hora, e.Tipo, e.Name)
	}
	fmt.Fprintln(fdot, "</table>>];")
	fmt.Fprintln(fdot, "}")
	fdot.Close()

	return utils.GenerateGraphvizReport(dotPath, dotFormat, path)
}

// TODO: rep -path=<destino> -id=<id> -name=journaling
func reporteJOURNALING(path string, id string) error {
	part := utils.GetPartitionById(id)
	if part == nil {
		return fmt.Errorf("partición no encontrada con id [%s]", id)
	}
	drive := strings.ToUpper(string(part.Id[0]))
	file, _, err := utils.OpenDisk(drive)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %v", err)
	}
	defer file.Close()

	var sb Structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		return fmt.Errorf("no se pudo leer el superbloque: %v", err)
	}
	if sb.S_filesystem_type != 3 {
		return fmt.Errorf("el sistema de archivos no es ext3, no se puede generar reporte de journaling")
	}
	journalingPos := int64(part.Start) + int64(binary.Size(Structs.Superblock{}))

	var journ Structs.Journaling
	if err := utils.ReadObject(file, &journ, journalingPos); err != nil {
		return fmt.Errorf("no se pudo leer journaling: %v", err)
	}

	dotPath, dotFormat, fdot, cleanup, err := utils.PrepareDotFile(path)
	if err != nil {
		return err
	}
	defer cleanup()

	fmt.Fprintln(fdot, "digraph JOURNALING {")
	fmt.Fprintln(fdot, "node [shape=plaintext fontname=\"Consolas\"];")
	fmt.Fprintln(fdot, "JOURNAL [label=<")
	fmt.Fprintln(fdot, "<table border='1' cellborder='1' cellspacing='2' cellpadding='8' bgcolor='#F0F4C3'>")
	fmt.Fprintln(fdot, "<tr>")
	fmt.Fprintln(fdot, "<td bgcolor='#827717'><font color='white'><b>Operacion</b></font></td>")
	fmt.Fprintln(fdot, "<td bgcolor='#827717'><font color='white'><b>Path</b></font></td>")
	fmt.Fprintln(fdot, "<td bgcolor='#827717'><font color='white'><b>Contenido</b></font></td>")
	fmt.Fprintln(fdot, "<td bgcolor='#827717'><font color='white'><b>Fecha</b></font></td>")
	fmt.Fprintln(fdot, "</tr>")
	for i := 0; i < int(journ.Ultimo) && i < len(journ.Contenido); i++ {
		entry := journ.Contenido[i]
		op := utils.CleanDOTString(strings.Trim(string(entry.Operation[:]), "\x00 "))
		if op == "" {
			break
		}
		pathJ := utils.CleanDOTString(strings.Trim(string(entry.Path[:]), "\x00 "))
		content := utils.CleanDOTString(strings.Trim(string(entry.Content[:]), "\x00 "))
		fecha := utils.CleanDOTString(strings.Trim(string(entry.Date[:]), "\x00 "))

		fmt.Fprintf(fdot, "<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n", op, pathJ, content, fecha)
	}

	fmt.Fprintln(fdot, "</table>>];")
	fmt.Fprintln(fdot, "}")
	fdot.Close()

	return utils.GenerateGraphvizReport(dotPath, dotFormat, path)
}
