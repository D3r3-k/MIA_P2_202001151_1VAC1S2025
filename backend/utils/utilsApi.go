package utilsApi

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"fmt"
	"os"
	"strings"
)

// [drives] Funciones para obtener información sobre los discos
func CountDisks() int {
	totalDisks := 0
	files, err := os.ReadDir(globals.PathDisks)
	if err != nil {
		return 0
	}
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if len(name) >= 4 && name[len(name)-4:] == ".dsk" {
				totalDisks++
			}
		}
	}
	return totalDisks
}

func CalculateTotalPartitions() (int, string) {
	totalPartitions := 0
	totalSize := int32(0)
	files, err := os.ReadDir(globals.PathDisks)
	if err != nil {
		return 0, "0 B"
	}
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if len(name) >= 4 && name[len(name)-4:] == ".dsk" {
				file, _, err := utils.OpenDisk(string(name[0]))
				if err != nil {
					continue
				}
				defer file.Close()

				var mbr Structs.MBR
				if err := utils.ReadObject(file, &mbr, 0); err != nil {
					continue
				}
				for _, partition := range mbr.Partitions {
					if partition.Status[0] == '1' {
						totalPartitions++
						totalSize += partition.Size
					}
				}
			}
		}
	}
	totalSizeStr := ConvertSizeToString(totalSize)
	return totalPartitions, totalSizeStr
}

func ConvertSizeToString(size int32) string {
	// convertir el tamaño a una cadena legible
	if size < 1024 {
		// mostrar 2 decimales
		return fmt.Sprintf("%.2f B", float64(size))
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
	}
}

// [drives] Funciones para obtener información de los discos
type DiskInfo struct {
	Name       string
	Path       string
	Size       int64
	Fit        string
	Partitions int
}

func GetDiskInfo() ([]DiskInfo, error) {
	var diskData []DiskInfo
	files, err := os.ReadDir(globals.PathDisks)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if len(name) >= 4 && name[len(name)-4:] == ".dsk" {
				f, _, err := utils.OpenDisk(string(name[0]))
				if err != nil {
					continue
				}
				defer f.Close()

				var mbr Structs.MBR
				if err := utils.ReadObject(f, &mbr, 0); err != nil {
					continue
				}
				diskInfo := DiskInfo{
					Name:       string(name[0]),
					Path:       strings.Split(f.Name(), "/")[len(strings.Split(f.Name(), "/"))-1],
					Size:       int64(mbr.MbrSize),
					Fit:        string(mbr.Fit[:]),
					Partitions: 0,
				}
				for _, partition := range mbr.Partitions {
					if partition.Status[0] == '1' {
						diskInfo.Partitions++
					}
				}
				diskData = append(diskData, diskInfo)
			}
		}
	}
	return diskData, nil
}
