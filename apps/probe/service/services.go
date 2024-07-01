package service

import (
	"apps/probe/service/filesystem"
	"apps/probe/service/shortcut"
)

var FileSystem *filesystem.FileSystemService
var ShortcutManage *shortcut.ShortcutManageService

func InitServices() {
	FileSystem = filesystem.NewFileSystemService()
	ShortcutManage = shortcut.NewShortcutManageService()
}
