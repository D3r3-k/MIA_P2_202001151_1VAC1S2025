package commands

import (
	"MIA_PI_202001151_1VAC1S2025/manager/utils"
	"fmt"
)

func Pause() {
	utils.ShowMessage("Ejecucion pausada. Precione cualquier tecla para continuar...", false)
	var input string
	fmt.Scanln(&input)
}
