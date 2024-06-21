package goreload

import "github.com/idreaminteractive/goreload/internal/commands"

func Reload(reloadServerUrl string) error {
	return commands.SignalReload(reloadServerUrl)
}
