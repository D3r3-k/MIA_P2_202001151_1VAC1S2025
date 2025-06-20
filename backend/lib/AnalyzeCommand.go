package lib

import (
	"MIA_PI_202001151_1VAC1S2025/commands"
	"MIA_PI_202001151_1VAC1S2025/utils"
	"fmt"
)

func AnalyzeCommand(comando string, parametros string) {
	switch comando {
	case "mkdisk":
		commands.Fn_Mkdisk(parametros)
	case "rmdisk":
		commands.Fn_rmdisk(parametros)
	case "fdisk":
		commands.Fn_Fdisk(parametros)
	case "mount":
		commands.Fn_Mount(parametros)
	case "unmount":
		commands.Fn_Unmount(parametros)
	case "mkfs":
		commands.Fn_Mkfs(parametros)
	case "login":
		commands.Fn_Login(parametros)
	case "logout":
		commands.Fn_Logout(parametros)
	case "mkgrp":
		commands.Fn_Mkgrp(parametros)
	case "rmgrp":
		commands.Fn_Rmgrp(parametros)
	case "mkusr":
		commands.Fn_Mkusr(parametros)
	case "rmusr":
		commands.Fn_Rmusr(parametros)
	case "mkfile":
		commands.Fn_Mkfile(parametros)
	case "cat":
		commands.Fn_Cat(parametros)
	case "mkdir":
		commands.Fn_Mkdir(parametros)
	case "find":
		commands.Fn_Find(parametros)
	case "pause":
		commands.Pause()
	case "rep":
		commands.Fn_Rep(parametros)
	default:
		utils.ShowMessage(fmt.Sprintf("Comando '%s' no reconocido.", comando), true)
	}
}
