package utilsApi

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"fmt"
	"os"
	"strings"
)

// [Drives] Funciones para obtener información sobre los discos
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

// [Drives] Funciones para obtener información de los discos
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

// [Drive] Funciones para obtener información de un disco específico
type DriveInfo struct {
	Name       string
	Path       string
	Size       int32
	Fit        string
	Partitions int
}

func GetDriveInfo(driveLetter string) (DriveInfo, error) {
	if len(driveLetter) != 1 {
		return DriveInfo{}, fmt.Errorf("invalid drive letter: %s", driveLetter)
	}

	f, _, err := utils.OpenDisk(driveLetter)
	if err != nil {
		return DriveInfo{}, err
	}
	defer f.Close()

	var mbr Structs.MBR
	if err := utils.ReadObject(f, &mbr, 0); err != nil {
		return DriveInfo{}, err
	}
	_fit := string(mbr.Fit[:])
	_fit = strings.Trim(_fit, "\x00")
	switch strings.ToUpper(_fit) {
	case "F":
		_fit = "Primer Ajuste"
	case "B":
		_fit = "Mejor Ajuste"
	case "W":
		_fit = "Peor Ajuste"
	default:
		_fit = "Sin Ajuste"
	}
	diskInfo := DriveInfo{
		Name:       string(driveLetter[0]),
		Path:       strings.Split(f.Name(), "/")[len(strings.Split(f.Name(), "/"))-1],
		Size:       mbr.MbrSize,
		Fit:        _fit,
		Partitions: 0,
	}
	for _, partition := range mbr.Partitions {
		if partition.Status[0] == '1' {
			diskInfo.Partitions++
		}
	}
	return diskInfo, nil
}

type PartitionInfo struct {
	Status     string
	Type       string
	Fit        string
	Size       string
	Start      int32
	Name       string
	ID         string
	Path       string
	Date       string
	Filesystem string
	Signature  string
}

func GetDrivePartitions(driveLetter string) ([]PartitionInfo, error) {
	Partitions := []PartitionInfo{}
	if len(driveLetter) != 1 {
		return nil, fmt.Errorf("invalid drive letter: %s", driveLetter)
	}

	f, _, err := utils.OpenDisk(driveLetter)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var mbr Structs.MBR
	if err := utils.ReadObject(f, &mbr, 0); err != nil {
		return nil, err
	}
	for _, partition := range mbr.Partitions {
		var sb Structs.Superblock
		if err := utils.ReadObject(f, &sb, int64(partition.Start)); err != nil {
			return nil, err
		}
		var newP PartitionInfo
		_type := string(partition.Type[:])
		_fit := string(partition.Fit[:])
		_name := string(partition.Name[:])
		_id := string(partition.Id[:])
		_date := string(mbr.CreationDate[:])
		_type = strings.Trim(_type, "\x00")
		_fit = strings.Trim(_fit, "\x00")
		_name = strings.Trim(_name, "\x00")
		_id = strings.Trim(_id, "\x00")
		_date = strings.Trim(_date, "\x00")
		_signature := strings.Trim(fmt.Sprintf("%d", mbr.Signature), "\x00")
		if partition.Status[0] == '1' {
			newP.Status = "Montada"
		} else {
			newP.Status = "Desmontada"
		}
		newP.Type = string(partition.Type[:])
		switch strings.ToUpper(_type) {
		case "P":
			newP.Type = "Primaria"
		case "E":
			newP.Type = "Extendida"
		default:
			newP.Type = "Sin Formato"
		}
		newP.Fit = strings.ToUpper(_fit)
		newP.Size = ConvertSizeToString(partition.Size)
		newP.Start = partition.Start
		newP.Name = _name
		newP.ID = _id
		newP.Path = globals.PathDisks + driveLetter
		newP.Date = _date
		switch sb.S_filesystem_type {
		case 2:
			newP.Filesystem = "Ext2"
		case 3:
			newP.Filesystem = "Ext3"
		default:
			newP.Filesystem = "Sin Formato"
		}
		newP.Signature = _signature
		Partitions = append(Partitions, newP)
	}
	return Partitions, nil
}

type StandardResponse struct {
	Error    string      `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
	Status   string      `json:"status"`
}
