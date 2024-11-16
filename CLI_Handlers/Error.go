package CLI_Handlers

import (
	"NullOps/Interface"
	"fmt"
)

func LogError(err error) {
	if err != nil {
		Interface.Info(fmt.Sprintf("Error: %s", err))
	}
}
