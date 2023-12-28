package test

import (
	"config-manager/common/model/shortcut"
	shortcut2 "config-manager/probe/service/shortcut"
	"fmt"
	"testing"
	"time"
)

func TestExecuteOneLineShortcut(t *testing.T) {

	oneLineShortcut := shortcut.NewOneLineShortcut("test", "echo hello", 10*time.Second, false, "no desc")

	lineShortcut := shortcut2.ExecuteOneLineShortcut(oneLineShortcut)

	fmt.Println(lineShortcut)

}

func TestScriptShortcut(t *testing.T) {

	scriptShortcut, err := shortcut.NewScriptShortcut("load", "/Users/meichuankutou/vscode/config-manage/server-probe/probe/test/scripts/tsest.sh", 10*time.Second, true, "no desc")

	if err != nil {
		t.Fatal(err)
	}

	result := shortcut2.ExecuteScriptShortcut(scriptShortcut)

	fmt.Println(result)

}
