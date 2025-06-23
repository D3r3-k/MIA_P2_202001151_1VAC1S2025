package lib

import (
	"MIA_PI_202001151_1VAC1S2025/manager/commands"
	"fmt"
)

func AnalyzeCommand(comando string, parametros string) (string, error) {
	switch comando {
	case "mkdisk":
		output, err := commands.Fn_Mkdisk(parametros)
		return output, err
	case "rmdisk":
		output, err := commands.Fn_rmdisk(parametros)
		return output, err
	case "fdisk":
		output, err := commands.Fn_Fdisk(parametros)
		return output, err
	case "mount":
		output, err := commands.Fn_Mount(parametros)
		return output, err
	case "unmount":
		output, err := commands.Fn_Unmount(parametros)
		return output, err
	case "mkfs":
		output, err := commands.Fn_Mkfs(parametros)
		return output, err
	case "login":
		output, err := commands.Fn_Login(parametros)
		msg := "Inicio de sesión exitoso"
		if !output {
			msg = "Error al iniciar sesión"
		}
		return msg, err
	case "logout":
		output, err := commands.Fn_Logout(parametros)
		return output, err
	case "mkgrp":
		output, err := commands.Fn_Mkgrp(parametros)
		return output, err
	case "rmgrp":
		output, err := commands.Fn_Rmgrp(parametros)
		return output, err
	case "mkusr":
		output, err := commands.Fn_Mkusr(parametros)
		return output, err
	case "rmusr":
		output, err := commands.Fn_Rmusr(parametros)
		return output, err
	case "mkfile":
		output, err := commands.Fn_Mkfile(parametros)
		return output, err
	case "cat":
		output, err := commands.Fn_Cat(parametros)
		msg := ""
		if output != nil {
			msg = "[Cat " + output.Name + "]\n" + "\"" + output.Content + "\""
		}
		return msg, err
	case "mkdir":
		output, err := commands.Fn_Mkdir(parametros)
		return output, err
	case "find":
		output := commands.Fn_Find(parametros)
		if output.Error != nil {
			return "", output.Error
		}
		msg := "[Find " + output.Object.Name + "]\n"
		if output.Tree != "" {
			msg += output.Tree + "\n"
		}
		return msg, nil
	case "pause":
		return "Pausa omitida", nil
	case "rep":
		output, err := commands.Fn_Rep(parametros)
		return output, err
	default:
		return "", fmt.Errorf("comando '%s' no reconocido", comando)
	}
}

func AnalyzeExistCommand(comando string) bool {
	switch comando {
	case "mkdisk", "rmdisk", "fdisk", "mount", "unmount", "mkfs", "login", "logout",
		"mkgrp", "rmgrp", "mkusr", "rmusr", "mkfile", "cat", "mkdir", "find", "pause", "rep":
		return true
	default:
		return false
	}
}
