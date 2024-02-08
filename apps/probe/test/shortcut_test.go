package test

import (
	"apps/probe/service/shortcut"
	"fmt"
	"testing"
)

func TestShortcut(t *testing.T) {

}

func TestCreateScript(t *testing.T) {

	path, err := shortcut.CreateScript("create file", "scripts", "echo 1231231231321\necho wadaw")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(*path)
}
