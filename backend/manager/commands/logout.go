package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
)

func Fn_Logout(params string) {
	paramDefs := map[string]Structs.ParamDef{}
	_, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return
	}
	if globals.LoginSession.User == "" {
		utils.ShowMessage("No hay ninguna sesión iniciada.", true)
		return
	}

	// Limpiar sesión
	globals.LoginSession = Structs.LoginSession{
		User:        "",
		Password:    "",
		PartitionID: [4]byte{},
	}
	utils.ShowMessage("Sesión cerrada correctamente.", false)
}
