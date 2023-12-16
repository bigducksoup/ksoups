package test

import (
	fileservice "config-manager/probe/service/fileService"
	"fmt"
	"testing"
)

func TestFile(t *testing.T) {

	path := "/Users/meichuankutou/Public/testrewrite.txt"

	file, err := fileservice.GetFile(path)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(file.GetContent())

	file.InsertLines(2, []string{"linex", "liney"})

	fmt.Println(file.GetContent())

	file.Flush()

}
