package commands

import (
	globals "MIA_PI_202001151_1VAC1S2025/manager/global"
	Structs "MIA_PI_202001151_1VAC1S2025/manager/structs"
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"fmt"
)

func Fn_Logout(params string) (string, error) {
	paramDefs := map[string]Structs.ParamDef{}
	_, err := utils.ParseParameters(params, paramDefs)
	if err != nil {
		utils.ShowMessage(err.Error(), true)
		return "", err
	}
	if globals.LoginSession.User == "" {
		utils.ShowMessage("No hay ninguna sesión iniciada.", true)
		return "", fmt.Errorf("no hay ninguna sesión iniciada")
	}

	// Limpiar sesión
	globals.LoginSession = Structs.LoginSession{
		User:        "",
		Password:    "",
		PartitionID: [4]byte{},
	}
	utils.ShowMessage("Sesión cerrada correctamente.", false)
	return "Sesión cerrada correctamente.", nil
}
