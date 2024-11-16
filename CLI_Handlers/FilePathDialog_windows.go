//go:build windows
// +build windows

package CLI_Handlers

import (
	"NullOps/Interface"
	"fmt"
	"github.com/sqweek/dialog"
)

func GetFilePath() string {
	Interface.Clear()
	FilePath, err := dialog.File().Title("Select Input File").Load()
	if err != nil {
		fmt.Println(err)
	}

	return FilePath
}
