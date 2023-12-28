package test

import (
	"config-manager/common/model/shortcut"
	"config-manager/probe/db"
	"fmt"
	"testing"
	"time"
)

func TestShortcut(t *testing.T) {

	err := db.InitDataBase()

	if err != nil {
		t.Fatal(err)
	}

	oneLineShortcut := shortcut.NewOneLineShortcut(
		"test",
		"echo 123",
		10*time.Second,
		false,
		"test",
	)

	r := db.DB.Create(&oneLineShortcut)
	fmt.Println(r.Error)

}
